package manager

import (
	"log"
	"os"

	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
	"github.com/pelletier/go-toml"

	"github.com/lakesite/zefram/pkg/api"
	"github.com/lakesite/zefram/pkg/models"
)

// ManagerService has a toml Config property which contains zefram specific directives,
// a DBConfig array of app database configurations.
type ManagerService struct {
	Config   *toml.Tree
	DBConfig map[string]*model.DBConfig
}

// InitApp initializes an app configuration, return true if successful false
// otherwise
func (ms *ManagerService) InitApp(app string) bool {
	if ms.DBConfig[app] == nil {
		ms.DBConfig[app] = &model.DBConfig{}
	}

	success := true

	// pull in the database config to DBConfig struct
	ms.DBConfig[app].Server, _ = ms.GetAppProperty(app, "dbserver")
	ms.DBConfig[app].Port, _ = ms.GetAppProperty(app, "dbport")
	ms.DBConfig[app].Database, _ = ms.GetAppProperty(app, "database")
	ms.DBConfig[app].User, _ = ms.GetAppProperty(app, "dbuser")
	ms.DBConfig[app].Password, _ = ms.GetAppProperty(app, "dbpassword")
	ms.DBConfig[app].Driver, _ = ms.GetAppProperty(app, "dbdriver")
	ms.DBConfig[app].Path, _ = ms.GetAppProperty(app, "dbpath")

	// Init the DB, which pulls in our gorm DB struct;
	ms.DBConfig[app].Init()

	// Finally, migrate if we haven't already;
	ms.DBConfig[app].Migrate()

	return success
}

// Init is required to initialize the manager service via a config file.
func (ms *ManagerService) Init(cfgfile string) {
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist.\n", cfgfile)
	} else {
		ms.Config, _ = toml.LoadFile(cfgfile)
		ms.DBConfig = make(map[string]*model.DBConfig)
		ms.InitApp("zefram")
	}
}

// Daemonize sets up the web service and defines routes for the API.
func (ms *ManagerService) Daemonize() {
	address := config.Getenv("ZEFRAM_HOST", "127.0.0.1") + ":" + config.Getenv("ZEFRAM_PORT", "7990")
	ws := service.NewWebService("zefram", address)
	api := api.NewAPI(
		ws,          // web service
		ms.Config,   // app configuration
		ms.DBConfig, // database connection and configuration
	)
	api.SetupRoutes()
	api.WebService.RunWebServer()
}
