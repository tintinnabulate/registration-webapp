package testers

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

// GetTestingContext : gets an aetest.Context aetest.Instance to pass into our tests.
// Always call inst.Close() when you are done with it, and at least at the end of each test.
// Use ToHTTPHandlerConverter to convert the aetest.Context to our normal AppEngine Context.
func GetTestingContext() (context.Context, aetest.Instance) {
	inst, _ := aetest.NewInstance(
		&aetest.Options{
			SuppressDevAppServerLog:     true,
			StronglyConsistentDatastore: true,
		})
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
	}
	ctx := appengine.NewContext(req)
	return ctx, inst
}
