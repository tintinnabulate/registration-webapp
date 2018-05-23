/*
	Implementation Note:
		None.

	Filename:
		db_operations.go
*/

package main

import (
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// StashRegistrationForm adds a signup with the given verification code to the datastore,
// returning the key of the newly created entity.
func StashRegistrationForm(ctx context.Context, regform *RegistrationForm) (*datastore.Key, error) {
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
	key := datastore.NewKey(ctx, "User", user.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, user)
	return k, err
}

func CreateConvention(ctx context.Context, convention *Convention) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "Convention", "", 0, nil)
	k, err := datastore.Put(ctx, key, convention)
	return k, err
}
