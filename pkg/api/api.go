// API contains the handlers to manage API endpoints
package api

import (
	"fmt"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/lakesite/ls-fibre/pkg/service"
	"github.com/pelletier/go-toml"

	"github.com/lakesite/zefram/pkg/mail"
	"github.com/lakesite/zefram/pkg/models"
)

type API struct {
	WebService *service.WebService
	AppConfig  *toml.Tree
	DBConfig   map[string]*model.DBConfig
}

// ContactHandler handles POST data for a Contact.
func (api *API) ContactHandler(w http.ResponseWriter, r *http.Request) {
	// parse the form
	err := r.ParseForm()
	if err != nil {
		api.WebService.JsonStatusResponse(w, "Error parsing form data.", http.StatusBadRequest)
		return
	}

	// create a new contact
	c := new(model.Contact)

	// using a new decoder, decode the form and bind it to the contact
	decoder := schema.NewDecoder()
	decoder.Decode(c, r.Form)

	// validate the structure:
	_, err = valid.ValidateStruct(c)
	if err != nil {
		api.WebService.JsonStatusResponse(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// insert the contact structure
	if dbc := api.DBConfig["zefram"].Connection.Create(c); dbc.Error != nil {
		api.WebService.JsonStatusResponse(w, fmt.Sprintf("Error: %s", dbc.Error.Error()), http.StatusInternalServerError)
		return
	}

	// send an email
	mailfrom, _ := api.GetAppProperty("zefram", "mailfrom")
	mailto, _ := api.GetAppProperty("zefram", "mailto")
	subject := "Contact from: " + c.Email
	body := "First Name: " + c.FirstName + "\nLast Name: " + c.LastName + "\nMessage: \n\n" + c.Message
	err = mail.LocalhostSendMail(mailfrom, mailto, subject, body)

	if err != nil {
		api.WebService.JsonStatusResponse(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return StatusOK with Contact made:
	api.WebService.JsonStatusResponse(w, "Contact made", http.StatusOK)
}

// Create a new API for the provided web service.
// ws is a ls-fibre WebService.
func NewAPI(ws *service.WebService, ac *toml.Tree, dbc map[string]*model.DBConfig) *API {
	api := &API{
		WebService: ws,
		AppConfig:  ac,
		DBConfig:   dbc,
	}

	return api
}
