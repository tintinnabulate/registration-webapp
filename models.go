package main

import "time"

type registrationForm struct {
	First_Name      string
	Email_Address   string
	Creation_Date   time.Time
	Member_Of       Fellowship
	Donation_Amount float32
}

type convention struct {
	ID                int64 // pk
	Name              string
	Creation_Date     time.Time
	Year              int
	Country           EURYPAA_Country
	City              string
	Cost              int
	Currency_Code     string
	Start_Date        time.Time
	End_Date          time.Time
	Hotel             string
	Hotel_Is_Venue    bool
	Venue             string
	Stripe_Product_ID string
}

// signup is used to parse JSON response from the signup microservice
type signup struct {
	Email_Address string `json:"address"`
	Success       bool   `json:"success"`
	Note          string `json:"note"`
}

type user struct {
	Email_Address      string
	Creation_Date      time.Time
	First_Name         string
	Member_Of          Fellowship
	Donation_Amount    int64
	Stripe_Customer_ID string
	Stripe_Charge_ID   string
}
