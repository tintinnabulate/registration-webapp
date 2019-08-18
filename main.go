package main

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
	"log"
	"net/http"
	"os"
	"rsc.io/quote"
	"time"
)

func main() {

	templatesInit()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, quote.Hello())
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

// templatesInit : parse the HTML templates, including any predefined functions (FuncMap)
func templatesInit() {
	templates = template.Must(template.New("").
		Funcs(funcMap).
		ParseGlob("templates/*.tmpl"))
}

var (
	templates  *template.Template
	translator *i18n.Bundle
)
