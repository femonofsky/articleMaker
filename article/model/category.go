package model

import "github.com/jinzhu/gorm"

//  Category Defines the structure for an Category
type Category struct {
	gorm.Model
	Name string `sql:"unique;not null" json:"name"`
}
