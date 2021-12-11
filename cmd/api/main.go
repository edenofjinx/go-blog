package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// version current version of the application
const version = "1.0.0"

//config model for application configuration
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

//Application holds the application data
type Application struct {
	Config   config
	Database models.MysqlModel
}

// main launches the application
func main() {
	var cfg config
	setEnvironment(&cfg)
	setDSN(&cfg)
	setServerPort(&cfg)

	db, err := openDatabase(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	err = models.SeedData(db)
	if err != nil {
		log.Println(err)
		return
	}

	app := &Application{
		Config:   cfg,
		Database: models.NewDatabase(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// openDatabase initialize the database session
func openDatabase(cfg config) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(cfg.db.dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

// loadEnvironment loads the environment file and sets the env name in config
func setEnvironment(cfg *config) {
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
			log.Fatal("Error loading .env file", err)
		}
	case "development":
		err := godotenv.Load(".env." + cfg.env + ".local")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	default:
		err := godotenv.Load(".env." + cfg.env + ".local")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}
}

// setDSN generates and sets the dsn for the database connection in config
func setDSN(cfg *config) {
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
func setServerPort(cfg *config) {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Println(err)
		return
	}
	cfg.port = port
}
