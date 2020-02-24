package model

import (
	"github.com/jinzhu/gorm"
)

// Publisher Defines the structure of a Publisher
type Publisher struct {
	gorm.Model
	Name string `sql:"unique;not null" json:"name"`
}
