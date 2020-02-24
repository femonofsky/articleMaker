package model

import (
	"fmt"
	"github.com/femonofsky/ArticleMaker/article/config"
	"github.com/jinzhu/gorm"
	// This loads the mysql database driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// This loads the postgres database driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// This loads the sqlite database driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Global Db
var Db *gorm.DB

// New return a GORM Database connected to either postgres or MYSQL Database
func New(config *config.Config) (*gorm.DB, error) {
	var connect string
	switch config.DB.Driver {
	case "mysql":
		// create string connection for msql
		connect = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.DB.Name, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	case "postgres":
		// create string connection for postgres
		connect = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name, config.DB.Password)
	case "sqlite3":
		connect = config.DB.Name
	default:
		return nil, fmt.Errorf("DB_DRIVER (%s) not support ", config.DB.Driver)
	}
	// Open Connection to Database
	DB, err := gorm.Open(config.DB.Driver, connect)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to %s database", config.DB.Driver)
	}
	if config.DB.Driver == "sqlite3" {
		DB.Exec("PRAGMA foreign_keys = ON")
	}
	Db = DB
	return DB, nil
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func Migrate(db *gorm.DB) *gorm.DB {
	// Database Migration the schema
	db.Debug().AutoMigrate(&Article{}, &Category{}, &Publisher{})
	db.Model(&Article{}).AddForeignKey("category_name", "categories(name)", "CASCADE", "CASCADE")
	db.Model(&Article{}).AddForeignKey("publisher_name", "publishers(name)", "CASCADE", "CASCADE")
	return db
}

const DateTimeLayout = "2006-01-02 15:04:05"


