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
	stripeInit()
}

// TestSubmitEmptyEmailAddress does just that
func TestSubmitEmptyEmailAddress(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	CreateConvention(ctx, cnv)

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
	CreateConvention(ctx, cnv)

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
				reg, err := GetRegistrationForm(ctx, viper.GetString("TestEmailAddress"))
				checkErr(err)
				c.So(reg.City, c.ShouldEqual, "Foo")
			})
		})
	})
}

/* func TestPayOverStripeCreatesUser(t *testing.T) {
	ctx, inst := getContext()
	defer inst.Close()
	cnv := &convention{Country: 1, Year: 2018, City: "Foo", Cost: 2000, Currency_Code: "EUR", Name: "EURYPAA"}
	CreateConvention(ctx, cnv)

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
				c.So(stripe.Key, c.ShouldEqual, viper.GetString("StripeTestSK"))
			})
		})
	})
}
*/
