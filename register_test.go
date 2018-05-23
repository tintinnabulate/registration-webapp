package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	c "github.com/smartystreets/goconvey/convey"

	"golang.org/x/net/context"
	//"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func CreateContextHandlerToHttpHandler(ctx context.Context) ContextHandlerToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(ctx, w, r)
		}
	}
}

func TestMonkeys(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/monkeys/dong", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, "banana: dong")
		})
	})
}

func TestCreateSignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/signup/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)
		})
	})
}

func TestCreateAndCheckSignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()
		record2 := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/signup/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {

			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)

			req2, err2 := http.NewRequest("GET", "/signup/lolz", nil)
			c.So(err2, c.ShouldBeNil)

			c.Convey("It should return a 200 response", func() {
				r.ServeHTTP(record2, req2)
				c.So(record2.Code, c.ShouldEqual, 200)
				c.So(fmt.Sprint(record2.Body), c.ShouldEqual, `{"address":"lolz","success":false,"note":""}
`)
			})
		})

	})
}

func TestVerifySignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/verify/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"code":"lolz","Success":false,"Note":"no such verification code"}
`)
		})
	})
}

// TODO test adding a valid signup and look for that code
