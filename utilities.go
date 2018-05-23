package main

import (
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: this will need adapting to whatever format we request for Sobriety_Date and Birth_Date
func TimeConverter(value string) reflect.Value {
	tstamp, _ := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(time.Unix(tstamp, 0))
}

// Standard http handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Our context.Context http handler
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// Higher order function for changing a HandlerFunc to a ContextHandlerFunc,
// usually creating the context.Context along the way.
type ContextHandlerToHandlerHOF func(f ContextHandlerFunc) HandlerFunc

// Creates a new Context and uses it when calling f
func ContextHandlerToHttpHandler(f ContextHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		f(ctx, w, r)
	}
}
