package handlers

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// HandlerFunc : a convenience type for our usual net/http Handler function signature
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// ContextHandlerFunc : a type that has the function signature of what we want our http Handlers
// to look like when we have an AppEngine context available
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// ToHandlerHOF : a type which gives us the the function signature that we need to
// take in a ContextHandlerFunc and convert it into a HandlerFunc,
// so that it can be used with regular http routing APIs.
type ToHandlerHOF func(f ContextHandlerFunc) HandlerFunc

// ToHTTPHandler : returns a HandlerFunc which internally creates a new appengine.Context
// and passes it through to our ContextHandlerFunc
func ToHTTPHandler(f ContextHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		f(ctx, w, r)
	}
}

// ToHTTPHandlerConverter : returns a higher order function that converts
// an aetest.Context handler to a standard HTTP handler.
func ToHTTPHandlerConverter(ctx context.Context) ToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(ctx, w, r)
		}
	}
}
