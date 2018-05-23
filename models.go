package main

import (
	"time"
)

// Email_Address is primary key
type User struct {
	Email_Address      string
	Creation_Date      time.Time
	First_Name         string
	Last_Name          string
	Password           string
	Conf_Password      string
	Country            CountryType
	City               string
	Sobriety_Date      time.Time
	Member_Of          []Fellowship
	Stripe_Customer_ID string
}

type Registration struct {
	ID                 int64  // pk
	User_Email_Address string // fk
	Convention_ID      int64  // fk
	Creation_Date      time.Time
	Stripe_Charge_ID   string
}

type Convention struct {
	ID                int64 // pk
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

// Email_Address is primary key
type RegistrationForm struct {
	Email_Address string
	Creation_Date time.Time
	First_Name    string
	Last_Name     string
	Password      string
	Conf_Password string
	Country       CountryType
	City          string
	Sobriety_Date time.Time
	Member_Of     []Fellowship
}
