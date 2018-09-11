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
	//appRouter.HandleFunc("/new_convention", f(GetNewConventionHandlerForm)).Methods("GET")
	//appRouter.HandleFunc("/new_convention", f(PostNewConventionHandlerForm)).Methods("POST")

	return appRouter
}

func GetSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	convention, err := GetLatestConvention(ctx)
	CheckErr(err)
	tmpl := templates.Lookup("signup_form.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Name":           convention.Name,
			"Year":           convention.Year,
			"City":           convention.City,
			"Country":        convention.Country,
			"Countries":      Countries,
			"Fellowships":    Fellowships,
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func PostSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	convention, err := GetLatestConvention(ctx)
	CheckErr(err)
	err = req.ParseForm()
	CheckErr(err)
	var signup Signup
	err = schemaDecoder.Decode(&signup, req.PostForm)
	client := urlfetch.Client(ctx)
	_, err = client.Post(fmt.Sprintf("%s/%s", config.SignupServiceURL, signup.Email_Address), "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl := templates.Lookup("check_email.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Name":           convention.Name,
			"Year":           convention.Year,
			"City":           convention.City,
			"Country":        convention.Country,
			"Countries":      Countries,
			"Fellowships":    Fellowships,
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func GetRegistrationFormHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	convention, err := GetLatestConvention(ctx)
	CheckErr(err)
	tmpl := templates.Lookup("registration_form.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Name":           convention.Name,
			"Year":           convention.Year,
			"City":           convention.City,
			"Country":        convention.Country,
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
	resp, err := client.Get(fmt.Sprintf("%s/%s", config.SignupServiceURL, regform.Email_Address))
	CheckErr(err)
	json.NewDecoder(resp.Body).Decode(&signup)
	if signup.Success {
		_, err := StashRegistrationForm(ctx, &regform)
		CheckErr(err)
		showPaymentForm(ctx, w, req, &regform)
	} else {
		http.Redirect(w, req, config.SignupURL, 301)
	}
}

func showPaymentForm(ctx context.Context, w http.ResponseWriter, req *http.Request, regform *RegistrationForm) {
	convention, err := GetLatestConvention(ctx)
	CheckErr(err)
	tmpl := templates.Lookup("stripe.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Name":           convention.Name,
			"Cost":           convention.Cost,
			"CostPrint":      convention.Cost / 100,
			"Currency":       convention.Currency_Code,
			"Year":           convention.Year,
			"City":           convention.City,
			"Country":        convention.Country,
			"Countries":      Countries,
			"Fellowships":    Fellowships,
			"Key":            publishableKey,
			csrf.TemplateTag: csrf.TemplateField(req),
			"Email":          regform.Email_Address,
		})
}

func PostRegistrationFormPaymentHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	convention, err := GetLatestConvention(ctx)
	CheckErr(err)
	req.ParseForm()

	emailAddress := req.Form.Get("stripeEmail")

	customerParams := &stripe.CustomerParams{Email: stripe.String(emailAddress)}
	customerParams.SetSource(req.Form.Get("stripeToken"))

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
	regform, err := GetRegistrationForm(ctx, emailAddress)
	CheckErr(err)
	user := &User{
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
	tmpl := templates.Lookup("registration_successful.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Name":           convention.Name,
			"Year":           convention.Year,
			"City":           convention.City,
			"Country":        convention.Country,
			"Countries":      Countries,
			"Fellowships":    Fellowships,
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func GetNewConventionHandlerForm(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	tmpl := templates.Lookup("new_convention.tmpl")
	tmpl.Execute(w,
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
	_, err = CreateConvention(ctx, &convention)
	CheckErr(err)
	fmt.Fprint(w, "Convention created")
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
	SignupServiceURL     string
	StripePublishableKey string
	StripeSecretKey      string
	StripeTestPK         string
	StripeTestSK         string
	TestEmailAddress     string
}

var (
	config         configuration
	schemaDecoder  *schema.Decoder
	publishableKey string
	templates      *template.Template
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
	schemaDecoder = schema.NewDecoder()
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

func TemplatesInit() {
	templates = template.Must(template.New("").Funcs(FuncMap).ParseGlob("templates/*.tmpl"))
}

func init() {
	ConfigInit()
	TemplatesInit()
	SchemaDecoderInit()
	RouterInit()
	StripeInit()
}
