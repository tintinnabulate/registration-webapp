package main

import (
	"github.com/gorilla/csrf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type templateVars map[string]interface{}

type pageInfo struct {
	convention convention
	email      string
	localizer  *i18n.Localizer
	r          *http.Request
}

func getVars(i *pageInfo) templateVars {
	return map[string]interface{}{
		"Name":           i.convention.Name,
		"Cost":           i.convention.Cost,
		"CostPrint":      i.convention.Cost / 100,
		"Currency":       i.convention.Currency_Code,
		"Year":           i.convention.Year,
		"City":           i.convention.City,
		"Country":        i.convention.Country,
		"Countries":      Countries,
		"Fellowships":    Fellowships,
		"Willings":       Willings,
		"HelpOutreaches": HelpOutreaches,
		"Tshirts":        Tshirts,
		csrf.TemplateTag: csrf.TemplateField(i.r),
		"Email":          i.email,
	}
}
