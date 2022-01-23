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
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// version current version of the application
const version = "1.0.0"
const staticImages = "static/images/"

//config model for application configuration
type serverConfig struct {
	port    int
	env     string
	migrate bool
	seed    bool
	db      struct {
		dsn string
	}
}

type envSet struct {
	fs *flag.FlagSet
}

func (e *envSet) setFlag() {
	e.fs = flag.NewFlagSet("envSet", flag.ContinueOnError)
	e.fs.StringVar(
		&cfg.env, "env",
		"development",
		"Environment variable. By default set to development. Can select from development|test|production",
	)
	e.fs.BoolVar(&cfg.migrate, "migrate", false, "Should auto migrate data?")
	e.fs.BoolVar(
		&cfg.seed,
		"seed",
		false,
		"Should seed data on launch? Seeding data is only available when migrate is set to true",
	)
}

func (e *envSet) parseEnvFlag() {
	args := os.Args[1:]
	for _, farg := range args {
		if e.fs.Name() == farg {
			err := e.fs.Parse(os.Args[2:])
			if err != nil {
				log.Println("envSet flag err")
			}
		}
	}
}

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger
var cfg serverConfig
var e envSet

// main launches the application
func main() {
	setEnvCfg()
	setAppCfg()
	db, err := setupDatabase()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	createHandlers(db)
	err = createServer().ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
		return
	}
}

// setEnvironment set the environment file and sets the env name in config
func (e *envSet) setEnvironment(cfg *serverConfig) {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	switch cfg.env {
	case "production":
		err := loadEnvFile(root+"/.env", true)
		if err != nil {
			errorLog.Fatal("Error loading .env file", err)
		}
	case "development":
		err := loadEnvFile(root+"/.env."+cfg.env+".local", false)
		if err != nil {
			errorLog.Fatal("Error loading .env file", err)
		}
	case "test":
		err := loadEnvFile(root+"/.env."+cfg.env, false)
		if err != nil {
			errorLog.Fatal("Error loading .env file", err)
		}
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
func setServerPort(cfg *serverConfig) error {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return err
	}
	cfg.port = port
	return nil
}

func loadEnvFile(path string, appInProd bool) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	app.InProduction = appInProd
	return nil
}

func setEnvCfg() {
	e.setFlag()
	e.parseEnvFlag()
	e.setEnvironment(&cfg)
	setDSN(&cfg)
	err := setServerPort(&cfg)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func setAppCfg() {
	app.AppVersion = version
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	app.StaticImages = staticImages
}

func setupDatabase() (*driver.DB, error) {
	db, err := driver.ConnectSQL(cfg.db.dsn, app)
	if err != nil {
		return nil, err
	}
	if cfg.migrate {
		err = database.MigrateData(db)
		if err != nil {
			return nil, err
		}
	}
	if cfg.migrate && cfg.seed {
		err = database.SeedData(db)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func createHandlers(db *driver.DB) {
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
}

func createServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
