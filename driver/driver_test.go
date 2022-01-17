package driver

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

type driverTestSuite struct {
	suite.Suite
	a config.AppConfig
}

var testsForDatabaseConnection = []struct {
	name          string
	dsn           string
	expectedError bool
}{
	{
		name:          "valid connection",
		dsn:           "root:root@tcp(localhost:3306)/go_blog_test?charset=utf8&parseTime=True&loc=Local",
		expectedError: false,
	},
	{
		name:          "invalid connection",
		dsn:           "fail:test@tcp(localhost:3306)/go_blog_test?charset=utf8&parseTime=True&loc=Local",
		expectedError: true,
	},
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
	for _, t := range testsForDatabaseConnection {
		d, err := ConnectSQL(t.dsn, suite.a)
		if err != nil {
			if t.expectedError == false {
				suite.Fail("did not expect a database connection error")
			}
		}
		assert.ObjectsAreEqualValues(gorm.DB{}, d)
	}
}

func (suite *driverTestSuite) TestNewDatabase() {
	for _, t := range testsForDatabaseConnection {
		d, err := NewDatabase(t.dsn, suite.a)
		if err != nil {
			if t.expectedError == false {
				suite.Fail("did not expect a database connection error")
			}
		}
		assert.ObjectsAreEqualValues(gorm.DB{}, d)
	}
}
