// handlers contains the handlers to manage API endpoints
package api

import (
	"fmt"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/lakesite/ls-mail"
	"github.com/lakesite/ls-governor"

	"github.com/lakesite/zefram/pkg/models"
)

// ContactHandler handles POST data for a Contact.
// The handler is wrapped to provide convenience access to a governor.API
func ContactHandler(w http.ResponseWriter, r *http.Request, gapi *governor.API) {
	// parse the form
	err := r.ParseForm()
	if err != nil {
		gapi.WebService.JsonStatusResponse(w, "Error parsing form data.", http.StatusBadRequest)
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
		gapi.WebService.JsonStatusResponse(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// insert the contact structure
	if dbc := gapi.ManagerService.DBConfig["zefram"].Connection.Create(c); dbc.Error != nil {
		gapi.WebService.JsonStatusResponse(w, fmt.Sprintf("Error: %s", dbc.Error.Error()), http.StatusInternalServerError)
		return
	}

	// send an email
	mailfrom, _ := gapi.ManagerService.GetAppProperty("zefram", "mailfrom")
	mailto, _ := gapi.ManagerService.GetAppProperty("zefram", "mailto")
	subject := "Contact from: " + c.Email
	body := "First Name: " + c.FirstName + "\nLast Name: " + c.LastName + "\nMessage: \n\n" + c.Message
	err = mail.LocalhostSendMail(mailfrom, mailto, subject, body)

	if err != nil {
		gapi.WebService.JsonStatusResponse(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return StatusOK with Contact made:
	gapi.WebService.JsonStatusResponse(w, "Contact made", http.StatusOK)
}
