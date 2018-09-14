/*
	Implementation Note:
		None.

	Filename:
		db_operations.go
*/

package main

import (
	"errors"
	"time"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// StashRegistrationForm adds a signup with the given verification code to the datastore,
// returning the key of the newly created entity.
func StashRegistrationForm(ctx context.Context, regform *registrationForm) (*datastore.Key, error) {
	regform.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "RegistrationForm", regform.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, regform)
	return k, err
}

// GetRegistrationForm gets the signup code matching the given email address.
// This should only be called during testing.
func GetRegistrationForm(ctx context.Context, email string) (registrationForm, error) {
	key := datastore.NewKey(ctx, "RegistrationForm", email, 0, nil)
	var regform registrationForm
	err := datastore.Get(ctx, key, &regform)
	return regform, err
}

// AddUser does a thing
func AddUser(ctx context.Context, u *user) (*datastore.Key, error) {
	u.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "User", u.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, u)
	return k, err
}

// CreateConvention : creates a convention in the Convention table
func CreateConvention(ctx context.Context, c *convention) (*datastore.Key, error) {
	c.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "Convention", "", 0, nil) // TODO: get it to use ID as the unique ID
	k, err := datastore.Put(ctx, key, c)
	return k, err
}

// getLatestConvention : gets the latest convention
func getLatestConvention(ctx context.Context) (convention, error) {
	var conventions []convention
	q := datastore.NewQuery("Convention").Order("-Creation_Date")
	_, err := q.GetAll(ctx, &conventions)
	checkErr(err)
	if len(conventions) < 1 {
		return convention{}, errors.New("No conventions in DB")
	}
	return conventions[0], nil
}
