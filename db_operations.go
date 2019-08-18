package main

import (
	"context"
	"time"
)

func getLatestConvention(context context.Context) (convention, error) {
	return convention{
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
	}, nil
}
