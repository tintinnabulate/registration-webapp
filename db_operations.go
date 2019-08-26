package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func getLatestConvention(ctx context.Context) (convention, error) {
	var theConvention convention
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return convention{}, fmt.Errorf("could not create datastore client: %v", err)
	}

	query := datastore.NewQuery("Convention").
		Order("-Creation_Date")

	it := client.Run(ctx, query)
	for {
		_, err := it.Next(&theConvention)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return convention{}, fmt.Errorf("convention not in DB: %v", err)
		}
		return theConvention, nil
	}
	return convention{}, fmt.Errorf("no conventions not in DB: %v", err)
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
