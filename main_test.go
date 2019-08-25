package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/tintinnabulate/vmail/mockverifier"
)

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	os.Exit(retCode)
}

func testSetup() {
	configInit("config.example.json")
	environmentInit()
	templatesInit()
	schemaDecoderInit()
	translatorInit()
	stripeInit()
	go mockverifier.Start(config.TestEmailAddress)
}

func TestSignupHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/signup", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getSignupHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "you@example.com"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)

	}
}

func TestPostSignupHandler(t *testing.T) {

	formData := url.Values{}
	formData.Set("Email_Address", config.TestEmailAddress)

	req, err := http.NewRequest("POST", "/signup", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postSignupHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "Please check your email"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)

	}
}

func TestPostSignupHandlerEmptyEmail(t *testing.T) {

	formData := url.Values{}
	formData.Set("Email_Address", "")

	req, err := http.NewRequest("POST", "/signup", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postSignupHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusNotFound,
		)
	}

	expected := "could not send verification email"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)

	}
}

func TestPostRegistrationHandler(t *testing.T) {

	formData := url.Values{}
	formData.Set("Email_Address", config.TestEmailAddress)
	formData.Set("Country", "1")
	formData.Set("City", "Foo")
	formData.Set("First_Name", "Bar")

	req, err := http.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postRegistrationFormHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "checkout-button"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)

	}
}

func TestGetRegistrationHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getRegistrationFormHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "Do you want to help outreach the convention?"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)

	}
}

func TestGetRegistrationHandlerSpanish(t *testing.T) {

	req, err := http.NewRequest("GET", "/register?lang=es", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getRegistrationFormHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "¿Quieres ayudar a promocionar el evento?"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)

	}
}
