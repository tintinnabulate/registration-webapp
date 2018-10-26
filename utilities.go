package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func getVars(convention convention, email string, localizer *i18n.Localizer, r *http.Request) templateVars {

	btnCompletePayment := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "btnCompletePayment",
			Other: "Complete payment to Register",
		},
	})
	btnSendVerifEmail := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "btnSendVerifEmail",
			Other: "Send verification email",
		},
	})
	btnContinueToCheckout := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "btnContinueToCheckout",
			Other: "Continue to checkout",
		},
	})
	errProcessPayment := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "errProcessPayment",
			Other: "Could not process payment",
		},
	})
	frmAmount := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmAmount",
			Other: "Amount {{ .CostPrint }} {{ .Currency }}",
		},
		TemplateData: map[string]string{
			"CostPrint": fmt.Sprintf("%d", convention.Cost/100),
			"Currency":  convention.Currency_Code,
		},
	})
	frmCity := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmCity",
			Other: "City",
		},
	})
	frmCountry := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmCountry",
			Other: "Country",
		},
	})
	frmEnterEmail := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmEnterEmail",
			Other: "Please enter your email address",
		},
	})
	frmFirstName := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmFirstName",
			Other: "First name",
		},
	})
	frmILiveIn := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmILiveIn",
			Other: "I live in...",
		},
	})
	frmPaymentDetails := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmPaymentDetails",
			Other: "Payment Details",
		},
	})
	frmSameEmail := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmSameEmail",
			Other: "Email - use the same one you verified with",
		},
	})
	frmYourDetails := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmYourDetails",
			Other: "Your details",
		},
	})
	pgCheckEmail := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgCheckEmail",
			Other: "Please check your email inbox, and click the link we've sent you",
		},
	})
	pgNowRegistered := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgNowRegistered",
			Other: "You are now registered!",
		},
	})
	pgRegisterFor := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgRegisterFor",
			Other: "Register for {{ .Name }}",
		},
		TemplateData: map[string]string{
			"Name": convention.Name,
		},
	})
	pgRegisteredFor := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgRegisteredFor",
			Other: "Registered for {{ .Name }}",
		},
		TemplateData: map[string]string{
			"Name": convention.Name,
		},
	})
	valEnterEmail := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "valEnterEmail",
			Other: "Please enter a valid email address so we can send you convention details.",
		},
	})
	valFirstName := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "valFirstName",
			Other: "Valid first name is required.",
		},
	})
	valSameEmail := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "valSameEmail",
			Other: "Please enter a valid email address so we can send you convention details.",
		},
	})

	helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HelloPerson",
			Other: "Hello {{.Name}}",
		},
		TemplateData: map[string]string{
			"Name": "Bob",
		},
	})

	return map[string]interface{}{
		"btnCompletePayment":    btnCompletePayment,
		"btnContinueToCheckout": btnContinueToCheckout,
		"btnSendVerifEmail":     btnSendVerifEmail,
		"errProcessPayment":     errProcessPayment,
		"frmAmount":             frmAmount,
		"frmCity":               frmCity,
		"frmCountry":            frmCountry,
		"frmEnterEmail":         frmEnterEmail,
		"frmFirstName":          frmFirstName,
		"frmILiveIn":            frmILiveIn,
		"frmPaymentDetails":     frmPaymentDetails,
		"frmSameEmail":          frmSameEmail,
		"frmYourDetails":        frmYourDetails,
		"pgCheckEmail":          pgCheckEmail,
		"pgNowRegistered":       pgNowRegistered,
		"pgRegisterFor":         pgRegisterFor,
		"pgRegisteredFor":       pgRegisteredFor,
		"valEnterEmail":         valEnterEmail,
		"valFirstName":          valFirstName,
		"valSameEmail":          valSameEmail,
		"Name":                  convention.Name,
		"Cost":                  convention.Cost,
		"CostPrint":             convention.Cost / 100,
		"Currency":              convention.Currency_Code,
		"Year":                  convention.Year,
		"City":                  convention.City,
		"Country":               convention.Country,
		"Countries":             Countries,
		"Fellowships":           Fellowships,
		"Key":                   publishableKey,
		csrf.TemplateTag:        csrf.TemplateField(r),
		"Email":                 email,
	}
}
