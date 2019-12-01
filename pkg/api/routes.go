// routes handles setting up routes for our API
package api

// SetupRoutes defines and associates routes to handlers.
func (api *API) SetupRoutes() {
	api.WebService.Router.HandleFunc("/zefram/api/v1/contact/", api.ContactHandler).Methods("POST")
}
