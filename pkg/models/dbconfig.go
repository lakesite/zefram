package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// DBConfig holds a database configuration.
type DBConfig struct {
	Connection *gorm.DB
	Server     string
	Port       string
	Database   string
	User       string
	Password   string
	Driver     string
	Path       string
}

func (db *DBConfig) Init() {
	if db.Driver == "sqlite3" {
		db.Connection, _ = gorm.Open("sqlite3", db.Path)
	} else {
		// handle connections with other drivers
	}
}

func (db *DBConfig) Migrate() {
	if db.Connection != nil {
		db.Connection.AutoMigrate(&Contact{})
	}
}
