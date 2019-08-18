package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tintinnabulate/gonfig"

	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	routerInit()
	templatesInit()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// signupHandler : show the signup form (SignupURL)
func signupHandler(w http.ResponseWriter, r *http.Request) {
	c := convention{
		Name:              "Name",
		Creation_Date:     time.Now(),
		Year:              2019,
		Country:           Albania_,
		City:              "City",
		Cost:              2000,
		Currency_Code:     "EUR",
		Start_Date:        time.Now(),
		End_Date:          time.Now(),
		Hotel:             "Hotel",
		Hotel_Is_Venue:    false,
		Venue:             "Venue",
		Stripe_Product_ID: "Stripe_Product_ID",
	}
	tmpl := templates.Lookup("signup_form.tmpl")
	page := &pageInfo{
		convention: c,
		localizer:  getLocalizer(r),
		r:          r}
	tmpl.Execute(w, getVars(page))
}

func getLocalizer(r *http.Request) *i18n.Localizer {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	return i18n.NewLocalizer(translator, lang, accept)
}

func createHTTPRouter() *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup", signupHandler).Methods("GET")
	return appRouter
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

func routerInit() {
	router := createHTTPRouter()
	csrfProtector := csrf.Protect(
		[]byte(config.CSRFKey),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

// templatesInit : parse the HTML templates, including any predefined functions (FuncMap)
func templatesInit() {
	templates = template.Must(template.New("").
		Funcs(funcMap).
		ParseGlob("templates/*.tmpl"))
}

// Config is our configuration file format
type Config struct {
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
	CSVUser              string `id:"CSVUser"              default:"CSVUser"`
}

var (
	templates  *template.Template
	translator *i18n.Bundle
	store      *sessions.CookieStore
	config     Config
)
