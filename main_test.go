package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	os.Exit(retCode)
}

func testSetup() {
	configInit("config.example.json")
	templatesInit()
	schemaDecoderInit()
	translatorInit()
	stripeInit()
	//go mockverifier.Start(config.TestEmailAddress)
}

func TestSignupHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/signup", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(signupHandler)
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
