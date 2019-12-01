package api

import (
	"fmt"
)

// get the property for app as a string, if property does not exist return err
func (api *API) GetAppProperty(app string, property string) (string, error) {
	if api.AppConfig.Get(app+"."+property) != nil {
		return api.AppConfig.Get(app + "." + property).(string), nil
	} else {
		return "", fmt.Errorf("Configuration missing '%s' section under [%s] heading.\n", property, app)
	}
}
