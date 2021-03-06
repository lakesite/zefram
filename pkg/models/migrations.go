package model

import (
	"errors"
	"fmt"

	"github.com/lakesite/ls-governor"
)

// Migrate takes a governor API and app name and migrates models, returns error
func Migrate(gapi *governor.API, app string) error {
	if gapi == nil {
		return errors.New("Migrate: Governor API is not initialized.")
	}

	if app == "" {
		return errors.New("Migrate: App name cannot be empty.")
	}

	dbc := gapi.ManagerService.DBConfig[app]
	
	if dbc == nil {
		return fmt.Errorf("Migrate: Database configuration for '%s' does not exist.", app)
	}

	if dbc.Connection == nil {
		return fmt.Errorf("Migrate: Database connection for '%s' does not exist.", app)
	}

	dbc.Connection.AutoMigrate(&Contact{})
	return nil
}
