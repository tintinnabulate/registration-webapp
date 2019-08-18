package main

import "time"

type registrationForm struct {
	Email_Address string
	Creation_Date time.Time
	First_Name    string
	Last_Name     string
	Password      string
	Conf_Password string
	Country       CountryType
	City          string
	Sobriety_Date time.Time
	Member_Of     Fellowship
	IsServant     Willing
	IsOutreacher  HelpOutreach
	IsTshirtBuyer Tshirt
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
