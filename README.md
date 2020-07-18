# registration-webapp

[![Build Status](https://travis-ci.org/tintinnabulate/registration-webapp.svg?branch=master)](https://travis-ci.org/tintinnabulate/registration-webapp)

This is a very simple registration website for an event, supporting payment using Stripe.

It is currently very tailored to its use in production, so it can't be used for any other event without significant development work.

It runs on Google App Engine, using Google Datastore to store registrations.

It depends on Go 1.12, and the Google App Engine 'stable' platform 1.12

## Features

* uses an email verification link to verify that registrants are using a real email address.
* only stores a registration if the user has paid successfully
* comprehensive test suite that creates test registrations and payments, hitting the Stripe test API.
* builds on Travis continuous integration servers using the `travis.yml` file.

## Email verification

The email verification is implemented in a microservice called vmail: <https://github.com/tintinnabulate/vmail>.

You have to run build and run this microservice first, before you can make use of it.

The microservice is used in two ways for email verification:

1. Submitting a `POST` request to `vmail_microservice_url/signup/{email_address}` to email the verification link to `{email_address}`
2. Submitting a `GET` request to `vmail_microservice_url/verify/{code}` to validate the verification code. This happens when you click the verification link sent by (1).

## User Registration & Payment

User registration also uses the vmail microservice.

A user can take the following steps through the registration process:

1. When user clicks "Register", if `GET vmail_microservice_url/signup/{email_address}` response JSON = `{"Address": "email_address", "Success": true, "Note": ""}` GOTO (2), else GOTO (3)
2. Take payment. If successful payment, GOTO (4), else GOTO (1).
3. Redirect to signup page
4. Store user object in Registrations database.

## Contributing

There is a Dockerfile which sets up a docker image to run the test suite. 

If you are familiar with Docker, you can also use this docker image to set up a
development environment.
