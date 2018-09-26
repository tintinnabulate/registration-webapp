package mockverifier

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var testEmailAddress string

// Email holds our JSON response for GET and POST /signup/{email}
type Email struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Note    string `json:"note"`
}

func isSignupVerified(w http.ResponseWriter, r *http.Request) {
	var email Email
	email.Address = testEmailAddress
	email.Success = mux.Vars(r)["email"] == testEmailAddress
	email.Note = ""
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(email)
}

// Start : starts the mock verifier
func Start(email string) {
	testEmailAddress = email
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup/{site_code}/{email}", isSignupVerified).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", appRouter))
}
