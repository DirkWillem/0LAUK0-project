package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/gcfg.v1"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type (
	// AppConfig contains all application configuration
	AppConfig struct {
		// Database connection settings
		Database struct {
			Host     string
			User     string
			Password string
			DBName   string
		}

		// JSON Web Token settings
		JWT struct {
			Secret string
		}

		// Web host settings
		Host struct {
			Host string
			Port string
		}
	}
)

var (
	db     *sql.DB
	config AppConfig
)

func init() {
	// Load the app config
	log.Println("Loading app config")
	err := gcfg.ReadFileInto(&config, "./config/app.cfg")

	if err != nil {
		LogErrorMessageFatal(fmt.Sprintf("Error reading app config: %s", err.Error()))
	}

	// Open database connection
	log.Println("Opening database connection")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", config.Database.User, config.Database.Password, config.Database.Host, config.Database.DBName)

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		LogErrorMessageFatal(fmt.Sprintf("Error opening database connection: %s", err.Error()))
	}
}
