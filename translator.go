package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type templateVars map[string]interface{}

type pageInfo struct {
	convention      convention
	email           string
	localizer       *i18n.Localizer
	r               *http.Request
	stripeSessionID string
	donationAmount  int
}

func getVars(i *pageInfo) templateVars {

	btnCompletePayment := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "btnCompletePayment",
			Other: "Complete payment to donate",
		},
	})
	btnSendVerifEmail := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "btnSendVerifEmail",
			Other: "Send verification email",
		},
	})
	btnContinueToCheckout := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "btnContinueToCheckout",
			Other: "Continue to checkout",
		},
	})
	errProcessPayment := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "errProcessPayment",
			Other: "Could not process payment",
		},
	})
	frmAmount := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmAmount",
			Other: "Amount {{ .CostPrint }} {{ .Currency }}",
		},
		TemplateData: map[string]string{
			"CostPrint": fmt.Sprintf("%d", i.convention.Cost/100),
			"Currency":  i.convention.Currency_Code,
		},
	})
	frmEnterEmail := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmEnterEmail",
			Other: "Please enter your email address",
		},
	})
	frmFirstName := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmFirstName",
			Other: "First name",
		},
	})
	frmDonationAmount := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmDonationAmount",
			Other: "Amount to donate ({{ .Currency }})",
		},
		TemplateData: map[string]string{
			"Currency": i.convention.Currency_Code,
		},
	})
	frmPaymentDetails := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmPaymentDetails",
			Other: "Payment Details",
		},
	})
	frmSameEmail := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmSameEmail",
			Other: "Email - use the same one you verified with",
		},
	})
	frmYourDetails := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmYourDetails",
			Other: "Your details",
		},
	})
	frmWhoAreYou := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "frmWhoAreYou",
			Other: "Who are you? Member of AA, Al-Anon or Guest?",
		},
	})
	pgCheckEmail := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgCheckEmail",
			Other: "Please check your email inbox, and click the link we've sent you",
		},
	})
	pgNowRegistered := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgNowRegistered",
			Other: "You have now donated! Thank you for your donation :)",
		},
	})
	pgRegisterFor := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgRegisterFor",
			Other: "Donate to {{ .Name }}",
		},
		TemplateData: map[string]string{
			"Name": i.convention.Name,
		},
	})
	pgRegisteredFor := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "pgRegisteredFor",
			Other: "Donated to {{ .Name }}",
		},
		TemplateData: map[string]string{
			"Name": i.convention.Name,
		},
	})
	valEnterEmail := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "valEnterEmail",
			Other: "Please enter a valid email address so we can send you convention details.",
		},
	})
	valFirstName := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "valFirstName",
			Other: "Valid first name is required.",
		},
	})
	valSameEmail := i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "valSameEmail",
			Other: "Please enter a valid email address so we can send you convention details.",
		},
	})

	return map[string]interface{}{
		"btnCompletePayment":    btnCompletePayment,
		"btnContinueToCheckout": btnContinueToCheckout,
		"btnSendVerifEmail":     btnSendVerifEmail,
		"errProcessPayment":     errProcessPayment,
		"frmAmount":             frmAmount,
		"frmDonationAmount":     frmDonationAmount,
		"frmEnterEmail":         frmEnterEmail,
		"frmFirstName":          frmFirstName,
		"frmPaymentDetails":     frmPaymentDetails,
		"frmSameEmail":          frmSameEmail,
		"frmYourDetails":        frmYourDetails,
		"frmWhoAreYou":          frmWhoAreYou,
		"pgCheckEmail":          pgCheckEmail,
		"pgNowRegistered":       pgNowRegistered,
		"pgRegisterFor":         pgRegisterFor,
		"pgRegisteredFor":       pgRegisteredFor,
		"valEnterEmail":         valEnterEmail,
		"valFirstName":          valFirstName,
		"valSameEmail":          valSameEmail,
		"Name":                  i.convention.Name,
		"Cost":                  i.convention.Cost,
		"CostPrint":             i.convention.Cost / 100,
		"DonationAmount":        i.donationAmount / 100,
		"Currency":              i.convention.Currency_Code,
		"Year":                  i.convention.Year,
		"City":                  i.convention.City,
		"Country":               i.convention.Country,
		"Fellowships":           Fellowships,
		"Key":                   stripePublishableKey,
		csrf.TemplateTag:        csrf.TemplateField(i.r),
		"Email":                 i.email,
		"StripeSessionID":       i.stripeSessionID,
	}
}
