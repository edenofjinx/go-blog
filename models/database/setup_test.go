package database

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type seedTestSuite struct {
	suite.Suite
	db *driver.DB
}

func (suite *seedTestSuite) SetupSuite() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	a := config.AppConfig{
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		InProduction: false,
		AppVersion:   "test",
		Environment:  "test",
		StaticImages: "static/test/images/",
	}
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
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
	d, err := driver.ConnectSQL(dsn, a)
	if err != nil {
		suite.Fail("could not connect to the test sql")
	}
	suite.db = d
}

func (suite *seedTestSuite) TearDownTest() {
	err := suite.db.SQL.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
}

func (suite *seedTestSuite) TearDownSuite() {
	err := suite.db.SQL.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
}
