package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

func getLatestConvention(context context.Context) (convention, error) {
	return convention{
		Name:              "EURYPAA",
		Creation_Date:     time.Now(),
		Year:              2020,
		Country:           Poland_,
		City:              "Warsaw",
		Cost:              1500,
		Currency_Code:     "EUR",
		Start_Date:        time.Now(),
		End_Date:          time.Now(),
		Hotel:             "Hotel",
		Hotel_Is_Venue:    false,
		Venue:             "Venue",
		Stripe_Product_ID: "Stripe_Product_ID",
	}, nil
}

// addUser : adds user to User table
func addUser(ctx context.Context, u *user) (*datastore.Key, error) {
	u.Creation_Date = time.Now()
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("could not create datastore client: %v", err)
	}
	key := datastore.IncompleteKey("User", nil)
	if _, err := client.Put(ctx, key, u); err != nil {
		return nil, fmt.Errorf("could not add user to user table: %v", err)
	}
	return key, nil
}
