package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	c "github.com/smartystreets/goconvey/convey"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func CreateContextHandlerToHTTPHandler(ctx context.Context) ContextHandlerToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(ctx, w, r)
		}
	}
}

func getContext() (context.Context, aetest.Instance) {
	inst, _ := aetest.NewInstance(
		&aetest.Options{
			StronglyConsistentDatastore: true,
			// SuppressDevAppServerLog:     true,
		})
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
	}
	ctx := appengine.NewContext(req)
	return ctx, inst
}

// TestSubmitEmptyEmailAddress does just that
func TestSubmitEmptyEmailAddress(t *testing.T) {
	ConfigInit()

	ctx, inst := getContext()
	defer inst.Close()

	cnv := &Convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}

	CreateConvention(ctx, cnv)

	c.Convey("When you submit a blank email address", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("Email_Address", "")

		req, err := http.NewRequest("POST", "/signup", strings.NewReader(formData.Encode())) // URL-encoded payload
		//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"Please check your email...\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `Please check your email inbox, and click the link we've sent you`)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `EURYPAA 2018 - Foo, Albania_`)
		})
	})
}

// TestRegisterWithValidEmail does just that
func TestRegisterWithValidEmail(t *testing.T) {
	ConfigInit()

	ctx, inst := getContext()
	defer inst.Close()

	cnv := &Convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}

	CreateConvention(ctx, cnv)

	c.Convey("When you register with a valid email address", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("Email_Address", config.TestEmailAddress)
		formData.Set("Country", "1")
		formData.Set("City", "Foo")
		formData.Set("First_Name", "Bar")

		req, err := http.NewRequest("POST", "/register", strings.NewReader(formData.Encode())) // URL-encoded payload
		//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"stripe-button\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `stripe-button`)
			c.Convey("There should be a registration entry in the Registration table", func() {
				reg, err := GetRegistrationForm(ctx, config.TestEmailAddress)
				CheckErr(err)
				c.So(reg.City, c.ShouldEqual, "Foo")
			})
		})
	})
}

//// TestPayOverStripeCreatesUser does just that
//func TestPayOverStripeCreatesUser(t *testing.T) {
//	ConfigInit()
//
//	ctx, inst := getContext()
//	defer inst.Close()
//
//	cnv := &Convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
//
//	CreateConvention(ctx, cnv)
//	c.Convey("When you make a payment over Stripe", t, func() {
//		c.Convey("There should be an user in the User table", func() {
//		})
//	})
//}
