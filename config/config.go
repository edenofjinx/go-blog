package config

import (
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	InProduction bool
	Environment  string
	AppVersion   string
	StaticImages string
}
