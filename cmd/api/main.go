package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"bitbucket.org/julius_liaudanskis/go-blog/models/database"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// version current version of the application
const version = "1.0.0"

//config model for application configuration
type serverConfig struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

// main launches the application
func main() {
	var cfg serverConfig
	setEnvironment(&cfg)
	setDSN(&cfg)
	setServerPort(&cfg)
	app.AppVersion = version

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	db, err := driver.ConnectSQL(cfg.db.dsn, app)
	if err != nil {
		errorLog.Fatal(err)
		return
	}
	err = database.SeedData(db)
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
		return
	}
}

// setEnvironment set the environment file and sets the env name in config
func setEnvironment(cfg *serverConfig) {
	flag.StringVar(
		&cfg.env,
		"env",
		"development",
		"Application env (development|production), default is set to development",
	)
	flag.Parse()
	switch cfg.env {
	case "production":
		err := godotenv.Load(".env")
		if err != nil {
			errorLog.Fatal("Error loading .env file", err)
		}
		app.InProduction = true
	case "development":
		err := godotenv.Load(".env." + cfg.env + ".local")
		if err != nil {
			errorLog.Fatal("Error loading .env file", err)
		}
		app.InProduction = false
	default:
		err := godotenv.Load(".env." + cfg.env + ".local")
		if err != nil {
			errorLog.Fatal("Error loading .env file", err)
		}
		app.InProduction = false
	}
	app.Environment = cfg.env
}

// setDSN generates and sets the dsn for the database connection in config
func setDSN(cfg *serverConfig) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	cfg.db.dsn = dsn
}

// setServerPort sets the server port in app config from the env file
func setServerPort(cfg *serverConfig) {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		errorLog.Fatal(err)
		return
	}
	cfg.port = port
}
