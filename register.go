package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stripe/stripe-go"
	stripeClient "github.com/stripe/stripe-go/client"
	"github.com/tintinnabulate/aecontext-handlers/handlers"
	"github.com/tintinnabulate/gonfig"

	"golang.org/x/net/context"
	"golang.org/x/text/language"

	"google.golang.org/appengine/urlfetch"
)

// createHTTPRouter : create a HTTP router where each handler is wrapped by a given context
func createHTTPRouter(f handlers.ToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup", f(getSignupHandler)).Methods("GET")
	appRouter.HandleFunc("/signup", f(postSignupHandler)).Methods("POST")
	appRouter.HandleFunc("/register", f(getRegistrationFormHandler)).Methods("GET")
	appRouter.HandleFunc("/register", f(postRegistrationFormHandler)).Methods("POST")
	appRouter.HandleFunc("/charge", f(postRegistrationFormPaymentHandler)).Methods("POST")
	return appRouter
}

// getSignupHandler : show the signup form (SignupURL)
func getSignupHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("signup_form.tmpl")
	l := getLocalizer(r)
	tmpl.Execute(w, getVars(convention, "", l, r))
}

// postSignupHandler : use the signup service to send the person a verification URL
func postSignupHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	err = r.ParseForm()
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
	httpClient := urlfetch.Client(ctx)
	resp, err := httpClient.Post(fmt.Sprintf("%s/%s", config.SignupServiceURL, s.Email_Address), "", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to email verifier: %v", err), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "could not send verification email", resp.StatusCode)
		return
	}
	tmpl := templates.Lookup("check_email.tmpl")
	l := getLocalizer(r)
	tmpl.Execute(w, getVars(convention, "", l, r))
}

// getRegistrationFormHandler : show the registration form
func getRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("registration_form.tmpl")
	l := getLocalizer(r)
	tmpl.Execute(w, getVars(convention, "", l, r))
}

// postRegistrationFormHandler : if they've signed up, show the payment form, otherwise redirect to SignupURL
func postRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	httpClient := urlfetch.Client(ctx)
	resp, err := httpClient.Get(fmt.Sprintf("%s/%s", config.SignupServiceURL, regform.Email_Address))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to email verifier: %v", err), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "could not verify email address", resp.StatusCode)
		return
	}
	json.NewDecoder(resp.Body).Decode(&s)
	session, err := store.Get(r, "regform")
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
		showPaymentForm(ctx, w, r, &regform)
	} else {
		http.Redirect(w, r, "/signup", http.StatusNotFound)
		return
	}
}

func showPaymentForm(ctx context.Context, w http.ResponseWriter, r *http.Request, regform *registrationForm) {
	convention, err := getLatestConvention(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("stripe.tmpl")
	l := getLocalizer(r)
	tmpl.Execute(w, getVars(convention, regform.Email_Address, l, r))
}

// postRegistrationFormPaymentHandler : charge the customer, and create a User in the User table
func postRegistrationFormPaymentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get latest convention: %v", err), http.StatusInternalServerError)
		return
	}
	r.ParseForm()

	emailAddress := r.Form.Get("stripeEmail")

	customerParams := &stripe.CustomerParams{Email: stripe.String(emailAddress)}
	customerParams.SetSource(r.Form.Get("stripeToken"))

	httpClient := urlfetch.Client(ctx)
	sc := stripeClient.New(stripe.Key, stripe.NewBackends(httpClient))

	newCustomer, err := sc.Customers.New(customerParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create customer: %v", err), http.StatusInternalServerError)
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(convention.Cost)),
		Currency:    stripe.String(convention.Currency_Code),
		Description: stripe.String(fmt.Sprintf("%s Registration", convention.Name)),
		Customer:    stripe.String(newCustomer.ID),
	}
	charge, err := sc.Charges.New(chargeParams)
	if err != nil {
		fmt.Fprintf(w, "Could not process payment: %v", err)
		return
	}
	session, err := store.Get(r, "regform")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create cookie session: %v", err), http.StatusInternalServerError)
		return
	}
	var regform *registrationForm
	if v, ok := session.Values["regform"].(*registrationForm); !ok {
		http.Error(w, "could not type assert value from cookie", http.StatusInternalServerError)
		return
	} else {
		regform = v
	}
	user := &user{
		First_Name:         regform.First_Name,
		Last_Name:          regform.Last_Name,
		Email_Address:      regform.Email_Address,
		Password:           regform.Password,
		Country:            regform.Country,
		City:               regform.City,
		Sobriety_Date:      regform.Sobriety_Date,
		Member_Of:          regform.Member_Of,
		Stripe_Customer_ID: charge.Customer.ID}
	_, err = addUser(ctx, user)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not add new user to user table: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("registration_successful.tmpl")
	l := getLocalizer(r)
	tmpl.Execute(w, getVars(convention, "", l, r))
}

// Config is our configuration file format
type Config struct {
	SMTPUsername         string `id:"SMTPUsername"         default:"sender@mydomain.com"`
	SMTPPassword         string `id:"SMTPPassword"         default:"mypassword"`
	SMTPServer           string `id:"SMTPServer"           default:"smtp.mydomain.com"`
	SiteDomain           string `id:"SiteDomain"           default:"mydomain.com"`
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
}

var (
	schemaDecoder  *schema.Decoder
	publishableKey string
	templates      *template.Template
	config         Config
	store          *sessions.CookieStore
	translator     *i18n.Bundle
)

func translatorInit() {
	translator = &i18n.Bundle{DefaultLanguage: language.English}
	translator.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	translator.MustLoadMessageFile("active.es.toml")
}

func getLocalizer(r *http.Request) *i18n.Localizer {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	return i18n.NewLocalizer(translator, lang, accept)
}

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
	store = sessions.NewCookieStore(
		[]byte(config.CookieStoreAuth),
		[]byte(config.CookieStoreEnc))
}

// schemaDecoderInit : create the schema decoder for decoding req.PostForm
func schemaDecoderInit() {
	schemaDecoder = schema.NewDecoder()
	schemaDecoder.RegisterConverter(time.Time{}, timeConverter)
	schemaDecoder.IgnoreUnknownKeys(true)
}

// routerInit : initialise our CSRF protected HTTPRouter
func routerInit() {
	// TODO: https://youtu.be/xyDkyFjzFVc?t=1308
	router := createHTTPRouter(handlers.ToHTTPHandler)
	csrfProtector := csrf.Protect(
		[]byte(config.CSRFKey),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

// stripeInit : set up important Stripe variables
func stripeInit() {
	if config.IsLiveSite {
		publishableKey = config.StripePublishableKey
		stripe.Key = config.StripeSecretKey
	} else {
		publishableKey = config.StripeTestPK
		stripe.Key = config.StripeTestSK
	}
}

// templatesInit : parse the HTML templates, including any predefined functions (FuncMap)
func templatesInit() {
	templates = template.Must(template.New("").
		Funcs(funcMap).
		ParseGlob("templates/*.tmpl"))
}

func init() {
	configInit("config.json")
	templatesInit()
	schemaDecoderInit()
	translatorInit()
	routerInit()
	stripeInit()
}
