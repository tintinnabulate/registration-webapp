package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stripe/stripe-go"
	stripeSession "github.com/stripe/stripe-go/checkout/session"
	"github.com/tintinnabulate/gonfig"
	"golang.org/x/text/language"
)

// globals
var (
	schemaDecoder        *schema.Decoder
	stripePublishableKey string
	templates            *template.Template
	config               Config
	cookieStore          *sessions.CookieStore
	translator           *i18n.Bundle
)

// init : initialize applicaiton
func init() {
	// Load configuration from config file
	configInit("config.json")
	// Set up environment variables
	environmentInit()
	// Load HTML templates from templates directory
	templatesInit()
	// Create a form decoder for parsing user form data
	schemaDecoderInit()
	// Create a translator for internationalisation
	translatorInit()
	// Set up all the routes e.g. /signup /register ...
	routerInit()
	// Set up Stripe payment platform
	stripeInit()
}

// main : main entry point to application
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Start here: http://localhost:%s/signup", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// createHTTPRouter : all of the routes that the web application handles
func createHTTPRouter() *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/", getSignupHandler).Methods("GET")
	appRouter.HandleFunc("/signup", getSignupHandler).Methods("GET")
	appRouter.HandleFunc("/signup", postSignupHandler).Methods("POST")
	appRouter.HandleFunc("/register", getRegistrationFormHandler).Methods("GET")
	appRouter.HandleFunc("/register", postRegistrationFormHandler).Methods("POST")
	appRouter.HandleFunc("/success", getSuccessHandler).Methods("GET")
	return appRouter
}

// getSignupHandler : (route) show the signup form (this is config.SignupURL)
func getSignupHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getLatestConvention(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("signup_form.tmpl")
	page := &pageInfo{
		convention: c,
		localizer:  getLocalizer(r),
		r:          r}
	tmpl.Execute(w, getVars(page))
}

// getRegistrationFormHandler : show the registration form
func getRegistrationFormHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getLatestConvention(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("registration_form.tmpl")
	page := &pageInfo{
		convention: c,
		localizer:  getLocalizer(r),
		r:          r,
	}
	tmpl.Execute(w, getVars(page))
}

// postSignupHandler : (route) use the signup service, vmail, to send the person a verification URL
func postSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse email form: %v", err), http.StatusInternalServerError)
		return
	}
	var s signup
	err = schemaDecoder.Decode(&s, r.PostForm)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode email address: %v", err), http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(fmt.Sprintf("%s/%s", config.SignupServiceURL, s.Email_Address), "", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to email verifier: %v", err), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "could not send verification email", resp.StatusCode)
		return
	}
	c, err := getLatestConvention(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("check_email.tmpl")
	page := &pageInfo{
		convention: c,
		localizer:  getLocalizer(r),
		r:          r,
	}
	tmpl.Execute(w, getVars(page))
}

// postRegistrationFormHandler : if they've signed up, show the payment form, otherwise redirect to SignupURL
func postRegistrationFormHandler(w http.ResponseWriter, r *http.Request) {
	var regform registrationForm
	var s signup
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse registration form: %v", err), http.StatusInternalServerError)
		return
	}
	err = schemaDecoder.Decode(&regform, r.PostForm)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode registration form: %v", err), http.StatusInternalServerError)
		return
	}
	resp, err := http.Get(fmt.Sprintf("%s/%s", config.SignupServiceURL, regform.Email_Address))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to email verifier: %v", err), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "could not verify email address", resp.StatusCode)
		return
	}
	json.NewDecoder(resp.Body).Decode(&s)
	session, err := cookieStore.Get(r, "regform")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create cookie session: %v", err), http.StatusInternalServerError)
		return
	}
	if s.Success {
		session.Values["regform"] = regform
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not save cookie session: %v", err), http.StatusInternalServerError)
			return
		}
		showPaymentForm(w, r, &regform)
	} else {
		http.Redirect(w, r, "/signup", http.StatusNotFound)
		return
	}
}

func showPaymentForm(w http.ResponseWriter, r *http.Request, regform *registrationForm) {
	convention, err := getLatestConvention(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}

	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(regform.Email_Address),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Name:        stripe.String(fmt.Sprintf("%s Donation", convention.Name)),
				Description: stripe.String(fmt.Sprintf("%s Donation", convention.Name)),
				Amount:      stripe.Int64(int64(regform.Donation_Amount * 100)),
				Currency:    stripe.String(string(convention.Currency_Code)),
				Quantity:    stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(config.SuccessURL),
		CancelURL:  stripe.String(config.SignupURL),
	}

	ss, err := stripeSession.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create stripe session: %v", err), http.StatusInternalServerError)
		return
	}

	user := &user{
		First_Name:         regform.First_Name,
		Email_Address:      regform.Email_Address,
		Member_Of:          regform.Member_Of,
		Donation_Amount:    regform.Donation_Amount * 100,
		Stripe_Customer_ID: "",
		Stripe_Charge_ID:   ss.PaymentIntent.ID,
	}
	_, err = addUser(r.Context(), user)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not add new user to user table: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl := templates.Lookup("stripe.tmpl")
	page := &pageInfo{
		convention: convention,
		email:      regform.Email_Address,
		localizer:  getLocalizer(r), r: r,
		stripeSessionID: ss.ID,
		donationAmount:  regform.Donation_Amount * 100,
	}
	tmpl.Execute(w, getVars(page))
}

// getSuccessHandler: (route) show the success message (this is config.SuccessURL)
func getSuccessHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getLatestConvention(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("registration_successful.tmpl")
	page := &pageInfo{
		convention: c,
		localizer:  getLocalizer(r),
		r:          r}
	tmpl.Execute(w, getVars(page))
}

// configInit : initialize the config using the config JSON file
func configInit(configName string) {
	err := gonfig.Load(&config, gonfig.Conf{
		FileDefaultFilename: configName,
		FileDecoder:         gonfig.DecoderJSON,
		FlagDisable:         true,
	})
	if err != nil {
		log.Fatalf("could not load configuration file: %v", err)
		return
	}
	gob.Register(&registrationForm{})
	cookieStore = sessions.NewCookieStore(
		[]byte(config.CookieStoreAuth),
		[]byte(config.CookieStoreEnc))
}

// templatesInit : parse the HTML templates, including any predefined functions (FuncMap)
func templatesInit() {
	templates = template.Must(template.New("").
		Funcs(funcMap).
		ParseGlob("templates/*.tmpl"))
}

// schemaDecoderInit : create the schema decoder for decoding req.PostForm
func schemaDecoderInit() {
	schemaDecoder = schema.NewDecoder()
	schemaDecoder.RegisterConverter(time.Time{}, timeConverter)
	schemaDecoder.IgnoreUnknownKeys(true)
}

// translatorInit : initialize the internationalisation handler
func translatorInit() {
	translator = i18n.NewBundle(language.English)
	translator.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	translator.MustLoadMessageFile("locales/active.es.toml")
}

// routerInit : initialize the application's CSRF-protected router
func routerInit() {
	router := createHTTPRouter()
	csrfProtector := csrf.Protect(
		[]byte(config.CSRFKey),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

// stripeInit : set up Stripe public and secret keys
func stripeInit() {
	if config.IsLiveSite {
		stripePublishableKey = config.StripePublishableKey
		stripe.Key = config.StripeSecretKey
	} else {
		stripePublishableKey = config.StripeTestPK
		stripe.Key = config.StripeTestSK
	}
}

// environmentInit : set up environment variables
func environmentInit() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GoogleCredentials)
}

// getLocalizer : used to localise content on a page given a *http.Request
func getLocalizer(r *http.Request) *i18n.Localizer {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	return i18n.NewLocalizer(translator, lang, accept)
}

// Config : the application's JSON configuration file format
type Config struct {
	// IsLiveSite : set this false during testing, true when deployed
	// SignupServiceURL : this is URL of the github.com/tintinnabulate/vmail deployment
	// TestEmailAddress : the email address that is used during testing
	// CSVUser : the special user that can download a CSV of all the registered Users
	// GoogleCredentials : the JSON Google Cloud Platform service account key
	SiteName             string `id:"SiteName"             default:"MyDomain"`
	ProjectID            string `id:"ProjectID"            default:"my-appspot-project-id"`
	CSRFKey              string `id:"CSRF_Key"             default:"my-random-32-bytes"`
	IsLiveSite           bool   `id:"IsLiveSite"           default:"false"`
	SignupURL            string `id:"SignupURL"            default:"this-apps-signup-endpoint.com/signup"`
	SignupServiceURL     string `id:"SignupServiceURL"     default:"http://localhost:10000/signup/eury2019"`
	StripePublishableKey string `id:"StripePublishableKey" default:"pk_live_foo"`
	StripeSecretKey      string `id:"StripeSecretKey"      default:"sk_live_foo"`
	StripeTestPK         string `id:"StripeTestPK"         default:"pk_test_UdWbULsYzTqKOob0SHEsTNN2"`
	StripeTestSK         string `id:"StripeTestSK"         default:"rk_test_xR1MFQcmds6aXvoDRKDD3HdR"`
	TestEmailAddress     string `id:"TestEmailAddress"     default:"foo@example.com"`
	CookieStoreAuth      string `id:"CookieStoreAuth"      default:"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
	CookieStoreEnc       string `id:"CookieStoreEnc"       default:"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
	CSVUser              string `id:"CSVUser"              default:"special-user@example.com"`
	GoogleCredentials    string `id:"GoogleCredentials"    default:"GoogleCredentialsJSONFile"`
	SuccessURL           string `id:"SuccessURL"           default:"this-apps-payment-success-url.com/success"`
}
