package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"

	"github.com/tintinnabulate/register/mockverifier"
)

func createContextHandlerToHTTPHandler(ctx context.Context) contextHandlerToHandlerHOF {
	return func(f contextHandlerFunc) handlerFunc {
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

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	os.Exit(retCode)
}

func testSetup() {
	viper.Set("IsLiveSite", false)
	viper.Set("SignupServiceURL", "http://localhost:10000/signup/eury2019")
	stripeInit()
	go mockverifier.Start(viper.GetString("TestEmailAddress"))
}

// TestGetSignupPage does just that
func TestGetSignupPage(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	createConvention(ctx, cnv)

	c.Convey("When visit the signup page", t, func() {
		r := createHTTPRouter(createContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/signup", nil) // URL-encoded payload

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"Please enter your email address\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `Please enter your email address`)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `EURYPAA 2018 - Foo, Albania_`)
		})
	})
}

// TestGetRegisterPage does just that
func TestGetRegisterPage(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	createConvention(ctx, cnv)

	c.Convey("When you visit the register page", t, func() {
		r := createHTTPRouter(createContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/register", nil) // URL-encoded payload

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"Continue to checkout\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `Continue to checkout`)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `EURYPAA 2018 - Foo, Albania_`)
		})
	})
}

// TestSubmitEmptyEmailAddress does just that
func TestSubmitEmptyEmailAddress(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	createConvention(ctx, cnv)

	c.Convey("When you submit a blank email address", t, func() {
		r := createHTTPRouter(createContextHandlerToHTTPHandler(ctx))
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
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `Please check your email inbox, and click the link we've sent you`)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `EURYPAA 2018 - Foo, Albania_`)
		})
	})
}

// TestRegisterWithValidEmail does just that
func TestRegisterWithValidEmail(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	createConvention(ctx, cnv)

	c.Convey("When you register with a valid email address", t, func() {
		r := createHTTPRouter(createContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("Email_Address", viper.GetString("TestEmailAddress"))
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
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `stripe-button`)
			c.Convey("There should be a registration entry in the Registration table", func() {
				reg, err := getRegistrationForm(ctx, viper.GetString("TestEmailAddress"))
				checkErr(err)
				c.So(reg.City, c.ShouldEqual, "Foo")
			})
		})
	})
}

// TestRegisterWithInvalidEmail does just that
func TestRegisterWithInvalidEmail(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	createConvention(ctx, cnv)

	c.Convey("When you register with an invalid email address", t, func() {
		r := createHTTPRouter(createContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("Email_Address", "thewrongemailaddress@notsignedup.glom")
		formData.Set("Country", "1")
		formData.Set("City", "Foo")
		formData.Set("First_Name", "Bar")

		req, err := http.NewRequest("POST", "/register", strings.NewReader(formData.Encode())) // URL-encoded payload
		//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

		c.So(err, c.ShouldBeNil)

		c.Convey("It should return http.StatusFound", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, http.StatusFound)
		})
	})
}

func TestPayOverStripeCreatesUser(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	createConvention(ctx, cnv)

	c.Convey("When you register with a valid email address", t, func() {
		r := createHTTPRouter(createContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()
		record1 := httptest.NewRecorder()

		formData1 := url.Values{}
		formData1.Set("Email_Address", viper.GetString("TestEmailAddress"))
		formData1.Set("Country", "1")
		formData1.Set("City", "Foo")
		formData1.Set("First_Name", "Bar")

		req1, err := http.NewRequest("POST", "/register", strings.NewReader(formData1.Encode())) // URL-encoded payload
		//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
		req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req1.Header.Add("Content-Length", strconv.Itoa(len(formData1.Encode())))

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"stripe-button\"", func() {
			r.ServeHTTP(record1, req1)
			c.So(record1.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record1.Body), c.ShouldContainSubstring, `stripe-button`)
			formData := url.Values{}
			formData.Set("stripeEmail", viper.GetString("TestEmailAddress"))
			formData.Set("stripeToken", "tok_visa")

			req, err := http.NewRequest("POST", "/charge", strings.NewReader(formData.Encode())) // URL-encoded payload
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

			c.So(err, c.ShouldBeNil)

			c.Convey("The next page body should contain \"You are now registered!\"", func() {
				r.ServeHTTP(record, req)
				c.So(record.Code, c.ShouldEqual, http.StatusOK)
				c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `You are now registered!`)
				c.Convey("There should be a user entry in the User table", func() {
					c.So(1, c.ShouldEqual, 1)
				})
			})
		})
	})
}
