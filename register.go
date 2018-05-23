package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/stripe/stripe-go"
	stripeClient "github.com/stripe/stripe-go/client"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup", f(GetSignupHandler)).Methods("GET")
	appRouter.HandleFunc("/signup", f(PostSignupHandler)).Methods("POST")
	appRouter.HandleFunc("/register", f(GetRegistrationFormHandler)).Methods("GET")
	appRouter.HandleFunc("/register", f(PostRegistrationFormHandler)).Methods("POST")
	appRouter.HandleFunc("/charge", f(PostRegistrationFormPaymentHandler)).Methods("POST")
	appRouter.HandleFunc("/new_convention", f(GetNewConventionHandlerForm)).Methods("GET")
	appRouter.HandleFunc("/new_convention", f(PostNewConventionHandlerForm)).Methods("POST")

	return appRouter
}

func GetSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w,
		"signup_form.tmpl",
		map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func PostSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	CheckErr(err)
	var signup Signup
	err = schemaDecoder.Decode(&signup, req.PostForm)
	client := urlfetch.Client(ctx)
	resp, err := client.Post(fmt.Sprintf("%s/%s", config.SignupURL, signup.Email_Address), "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)
}

func GetRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	t, err := template.New("registration_form.tmpl").Funcs(funcMap).ParseFiles("registration_form.tmpl")
	CheckErr(err)
	t.ExecuteTemplate(w,
		"registration_form.tmpl",
		map[string]interface{}{
			"Countries":      Countries,
			"Fellowships":    Fellowships,
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func PostRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var regform RegistrationForm
	var signup Signup
	err := req.ParseForm()
	CheckErr(err)
	err = schemaDecoder.Decode(&regform, req.PostForm)
	CheckErr(err)
	client := urlfetch.Client(ctx)
	resp, err := client.Get(fmt.Sprintf("%s/%s", config.SignupURL, regform.Email_Address))
	CheckErr(err)
	json.NewDecoder(resp.Body).Decode(&signup)
	if signup.Success {
		regform.Creation_Date = time.Now()
		_, err := StashRegistrationForm(ctx, &regform)
		CheckErr(err)
		showPaymentForm(ctx, w, req, &regform)
	} else {
		fmt.Fprint(w, "I'm sorry, you need to sign up first. Go to /signup")
	}
}

func showPaymentForm(ctx context.Context, w http.ResponseWriter, req *http.Request, regform *RegistrationForm) {
	tmpl := templates.Lookup("stripe.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Key":            publishableKey,
			csrf.TemplateTag: csrf.TemplateField(req),
			"Email":          regform.Email_Address,
		})
}

func PostRegistrationFormPaymentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	emailAddress := r.Form.Get("stripeEmail")

	customerParams := &stripe.CustomerParams{Email: emailAddress}
	customerParams.SetSource(r.Form.Get("stripeToken"))

	httpClient := urlfetch.Client(ctx)
	sc := stripeClient.New(stripe.Key, stripe.NewBackends(httpClient))

	newCustomer, err := sc.Customers.New(customerParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   500,
		Currency: "usd",
		Desc:     "Sample Charge",
		Customer: newCustomer.ID,
	}
	charge, err := sc.Charges.New(chargeParams)
	if err != nil {
		fmt.Fprintf(w, "Could not process payment: %v", err)
		return
	}
	regform, err := GetRegistrationForm(ctx, emailAddress)
	CheckErr(err)
	user := &User{
		Creation_Date:      time.Now(),
		First_Name:         regform.First_Name,
		Last_Name:          regform.Last_Name,
		Email_Address:      regform.Email_Address,
		Password:           regform.Password,
		Country:            regform.Country,
		City:               regform.City,
		Sobriety_Date:      regform.Sobriety_Date,
		Member_Of:          regform.Member_Of,
		Stripe_Customer_ID: charge.Customer.ID}
	_, err = AddUser(ctx, user)
	CheckErr(err)
	fmt.Fprintf(w, "Completed payment! Well not really... this was a test :-P")
}

func GetNewConventionHandlerForm(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	t, err := template.New("new_convention.tmpl").Funcs(funcMap).ParseFiles("new_convention.tmpl")
	CheckErr(err)
	t.ExecuteTemplate(w,
		"new_convention.tmpl",
		map[string]interface{}{
			"Countries":      EURYPAA_Countries,
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func PostNewConventionHandlerForm(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var convention Convention
	err := req.ParseForm()
	CheckErr(err)
	err = schemaDecoder.Decode(&convention, req.PostForm)
	CheckErr(err)
	convention.Creation_Date = time.Now()
	_, err = CreateConvention(ctx, &convention)
	CheckErr(err)
	fmt.Fprint(w, "Convention created")
}

type Signup struct {
	Email_Address string `json:"address"`
	Success       bool   `json:"success"`
	Note          string `json:"note"`
}

type configuration struct {
	SiteName             string
	SiteDomain           string
	SMTPServer           string
	SMTPUsername         string
	SMTPPassword         string
	ProjectID            string
	CSRF_Key             string
	IsLiveSite           bool
	SignupURL            string
	StripePublishableKey string
	StripeSecretKey      string
}

var (
	config         configuration
	schemaDecoder  = schema.NewDecoder()
	funcMap        = template.FuncMap{"inc": func(i int) int { return i + 1 }}
	publishableKey string
	templates      = template.Must(template.ParseGlob("views/*.tmpl"))
)

func ConfigInit() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	CheckErr(err)
}

func SchemaDecoderInit() {
	schemaDecoder.RegisterConverter(time.Time{}, TimeConverter)
	schemaDecoder.IgnoreUnknownKeys(true)
}

func RouterInit() {
	// TODO: https://youtu.be/xyDkyFjzFVc?t=1308
	router := CreateHandler(ContextHandlerToHttpHandler)
	csrfProtector := csrf.Protect(
		[]byte(config.CSRF_Key),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

func StripeInit() {
	publishableKey = config.StripePublishableKey
	stripe.Key = config.StripeSecretKey
}

func init() {
	ConfigInit()
	SchemaDecoderInit()
	RouterInit()
	StripeInit()
}
