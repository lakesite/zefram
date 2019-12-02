// routes handles setting up routes for our API
package api

import (
	"net/http"

	"github.com/lakesite/ls-governor"
)

// SetupRoutes defines and associates routes to handlers.
// Use a wrapper convention to pass a governor API to each handler.
func SetupRoutes(gapi *governor.API) {
	gapi.WebService.Router.HandleFunc(
		"/zefram/api/v1/contact/", 
		func(w http.ResponseWriter, r *http.Request) {
			ContactHandler(w, r, gapi)
		},
	).Methods("POST")
}
