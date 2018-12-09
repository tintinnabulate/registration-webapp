/*
	Implementation Note:
		None.

	Filename:
		db_operations.go
*/

package main

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// addUser : adds user to User table
func addUser(ctx context.Context, u *user) (*datastore.Key, error) {
	u.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "User", u.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, u)
	if err != nil {
		return nil, fmt.Errorf("could not add user to user table: %v", err)
	}
	return k, nil
}

// getUser : gets the user matching the given email address.
// This should only be called during testing.
func getUser(ctx context.Context, email string) (user, error) {
	key := datastore.NewKey(ctx, "User", email, 0, nil)
	var u user
	err := datastore.Get(ctx, key, &u)
	if err != nil {
		return user{},
			fmt.Errorf("could not get user: %v", err)
	}
	return u, nil
}

// createConvention : creates a convention in the Convention table
func createConvention(ctx context.Context, c *convention) (*datastore.Key, error) {
	c.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "Convention", "", 0, nil) // TODO: get it to use ID as the unique ID
	k, err := datastore.Put(ctx, key, c)
	if err != nil {
		return nil, fmt.Errorf("could not store convention in convention table: %v", err)
	}
	return k, nil
}

// getLatestConvention : gets the latest convention
func getLatestConvention(ctx context.Context) (convention, error) {
	var conventions []convention
	q := datastore.NewQuery("Convention").Order("-Creation_Date")
	_, err := q.GetAll(ctx, &conventions)
	if err != nil {
		return convention{}, fmt.Errorf("could not get latest convention: %v", err)
	}
	if len(conventions) < 1 {
		return convention{}, errors.New("No conventions in DB")
	}
	return conventions[0], nil
}
