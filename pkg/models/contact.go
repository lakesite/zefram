package model

import (
	"github.com/jinzhu/gorm"
)

// Contact holds the minimal data for a web contact.
type Contact struct {
	gorm.Model
	FirstName  string
	LastName   string
	Email      string `gorm:"type:varchar(100);"valid:"email"`
	Message    string
	AddMailing bool `gorm:"default:false;"`
}
