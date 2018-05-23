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
func StashRegistrationForm(ctx context.Context, regform *RegistrationForm) (*datastore.Key, error) {
	regform.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "RegistrationForm", regform.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, regform)
	return k, err
}

// GetRegistrationForm gets the signup code matching the given email address.
// This should only be called during testing.
func GetRegistrationForm(ctx context.Context, email string) (RegistrationForm, error) {
	key := datastore.NewKey(ctx, "RegistrationForm", email, 0, nil)
	var regform RegistrationForm
	err := datastore.Get(ctx, key, &regform)
	return regform, err
}

// AddUser does a thing
func AddUser(ctx context.Context, user *User) (*datastore.Key, error) {
	user.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "User", user.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, user)
	return k, err
}

func CreateConvention(ctx context.Context, convention *Convention) (*datastore.Key, error) {
	convention.Creation_Date = time.Now()
	key := datastore.NewKey(ctx, "Convention", "", 0, nil) // TODO: get it to use ID as the unique ID
	k, err := datastore.Put(ctx, key, convention)
	return k, err
}

func GetLatestConvention(ctx context.Context) (Convention, error) {
	var conventions []Convention
	q := datastore.NewQuery("Convention").Order("-Creation_Date")
	_, err := q.GetAll(ctx, &conventions)
	CheckErr(err)
	if len(conventions) < 1 {
		return Convention{}, errors.New("No conventions in DB")
	}
	return conventions[0], nil
}
