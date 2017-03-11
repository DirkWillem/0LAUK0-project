package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/gcfg.v1"
	"log"

	_ "github.com/lib/pq"
	"main/utils"
	"os"
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
			UseEnvDBString bool
		}

		// JSON Web Token settings
		JWT struct {
			Secret string
		}

		// Web host settings
		Host struct {
			Host       string
			Port       string
			UseEnvPort bool
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

	configFile := "app.cfg"
	if envConfigFile := os.Getenv("CONFIG_FILE"); len(envConfigFile) > 0 {
		configFile = envConfigFile
	}

	err := gcfg.ReadFileInto(&config, fmt.Sprintf("./config/%s", configFile))

	if err != nil {
		utils.LogErrorMessageFatal(fmt.Sprintf("Error reading app config: %s", err.Error()))
	}

	// Open database connection
	log.Println("Opening database connection")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.Database.User, config.Database.Password, config.Database.Host, config.Database.DBName)

	if config.Database.UseEnvDBString {
		connectionString = os.Getenv("DATABASE_URL")
	}

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		utils.LogErrorMessageFatal(fmt.Sprintf("Error opening database connection: %s", err.Error()))
	}

	if config.Host.UseEnvPort {
		config.Host.Port = os.Getenv("PORT")
	}
}
