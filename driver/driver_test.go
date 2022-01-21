package driver

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type driverTestSuite struct {
	suite.Suite
	a config.AppConfig
}

func (suite *driverTestSuite) SetupSuite() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := config.AppConfig{
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		InProduction: false,
		AppVersion:   "test",
		Environment:  "test",
		StaticImages: "static/test/images/",
	}
	suite.a = app
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(driverTestSuite))
}

func (suite *driverTestSuite) TestConnectSQL() {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../")
	err := godotenv.Load(root + "/.env.test")
	suite.Nil(err, "should not throw an error when loading .env file")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	d, err := ConnectSQL(dsn, suite.a)
	suite.Nil(err, "should not throw an error with correct dsn settings in .env file")
	assert.ObjectsAreEqualValues(gorm.DB{}, d)
	incorrectDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"testingIncorrectDsn",
		"withDummyData",
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	d, err = ConnectSQL(incorrectDsn, suite.a)
	suite.Error(err, "should contain an error as dummy dsn data is provided")
}

func (suite *driverTestSuite) TestNewDatabase() {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../")
	err := godotenv.Load(root + "/.env.test")
	suite.Nil(err, "should not throw an error when loading .env file")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	d, err := NewDatabase(dsn, suite.a)
	suite.Nil(err, "should not throw an error with correct dsn settings in .env file")
	assert.ObjectsAreEqualValues(gorm.DB{}, d)
	incorrectDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"testingIncorrectDsn",
		"withDummyData",
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	d, err = NewDatabase(incorrectDsn, suite.a)
	suite.Error(err, "should contain an error as dummy dsn data is provided")
}
