package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/stripe/stripe-go"
	stripeClient "github.com/stripe/stripe-go/client"
	"github.com/tintinnabulate/aecontext-handlers/handlers"
	"github.com/tintinnabulate/gonfig"

	"golang.org/x/net/context"

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
	checkErr(err)
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create cookie session: %v", err), http.StatusInternalServerError)
		return
	}
	session.Values["foo"] = "LAZERS"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "could not save cookie session", http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("signup_form.tmpl")
	tmpl.Execute(w, getVars(convention, "", r))
}

// postSignupHandler : use the signup service to send the person a verification URL
func postSignupHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	checkErr(err)
	err = r.ParseForm()
	checkErr(err)
	var s signup
	err = schemaDecoder.Decode(&s, r.PostForm)
	httpClient := urlfetch.Client(ctx)
	resp, err := httpClient.Post(fmt.Sprintf("%s/%s", config.SignupServiceURL, s.Email_Address), "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "could not send verification email", resp.StatusCode)
		return
	}
	tmpl := templates.Lookup("check_email.tmpl")
	tmpl.Execute(w, getVars(convention, "", r))
}

// getRegistrationFormHandler : show the registration form
func getRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	checkErr(err)
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "could not create cookie session", http.StatusInternalServerError)
		return
	}
	var email string
	if v, ok := session.Values["foo"].(string); !ok {
		http.Error(w, "could not type assert value from cookie", http.StatusInternalServerError)
		return
	} else {
		email = v
	}
	tmpl := templates.Lookup("registration_form.tmpl")
	tmpl.Execute(w, getVars(convention, email, r))
}

// postRegistrationFormHandler : if they've signed up, show the payment form, otherwise redirect to SignupURL
func postRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var regform registrationForm
	var s signup
	err := r.ParseForm()
	checkErr(err)
	err = schemaDecoder.Decode(&regform, r.PostForm)
	checkErr(err)
	httpClient := urlfetch.Client(ctx)
	resp, err := httpClient.Get(fmt.Sprintf("%s/%s", config.SignupServiceURL, regform.Email_Address))
	checkErr(err)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "could not verify email address", resp.StatusCode)
		return
	}
	json.NewDecoder(resp.Body).Decode(&s)
	if s.Success {
		_, err := stashRegistrationForm(ctx, &regform)
		checkErr(err)
		showPaymentForm(ctx, w, r, &regform)
	} else {
		http.Redirect(w, r, "/signup", http.StatusNotFound)
		return
	}
}

func showPaymentForm(ctx context.Context, w http.ResponseWriter, r *http.Request, regform *registrationForm) {
	convention, err := getLatestConvention(ctx)
	checkErr(err)
	tmpl := templates.Lookup("stripe.tmpl")
	tmpl.Execute(w, getVars(convention, regform.Email_Address, r))
}

// postRegistrationFormPaymentHandler : charge the customer, and create a User in the User table
func postRegistrationFormPaymentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	convention, err := getLatestConvention(ctx)
	checkErr(err)
	r.ParseForm()

	emailAddress := r.Form.Get("stripeEmail")

	customerParams := &stripe.CustomerParams{Email: stripe.String(emailAddress)}
	customerParams.SetSource(r.Form.Get("stripeToken"))

	httpClient := urlfetch.Client(ctx)
	sc := stripeClient.New(stripe.Key, stripe.NewBackends(httpClient))

	newCustomer, err := sc.Customers.New(customerParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	regform, err := getRegistrationForm(ctx, emailAddress)
	checkErr(err)
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
	checkErr(err)
	tmpl := templates.Lookup("registration_successful.tmpl")
	tmpl.Execute(w, getVars(convention, "", r))
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
)

func configInit(configName string) {
	err := gonfig.Load(&config, gonfig.Conf{
		FileDefaultFilename: configName,
		FileDecoder:         gonfig.DecoderJSON,
		FlagDisable:         true,
	})
	checkErr(err)
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
	routerInit()
	stripeInit()
}
