package main

import (
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: this will need adapting to whatever format we request for Sobriety_Date and Birth_Date
func timeConverter(value string) reflect.Value {
	tstamp, _ := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(time.Unix(tstamp, 0))
}

type templateVars map[string]interface{}

func getVars(convention convention, email string, r *http.Request) templateVars {
	return map[string]interface{}{
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
		csrf.TemplateTag: csrf.TemplateField(r),
		"Email":          email,
	}
}
