package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
)

type verifyHandler struct {
	testHandlerRepo *handlers.Repository
	testSQL         *driver.DB
}

var hr verifyHandler

var testsToVerifyApiKey = []struct {
	name         string
	url          string
	apiKey       string
	expectedCode int
}{
	{
		name:         "test status handler without api key",
		url:          "/v1/status",
		apiKey:       "",
		expectedCode: http.StatusAccepted,
	},
	{
		name:         "test status handler with api key",
		url:          "/v1/status",
		apiKey:       testUsers[0].ApiKey,
		expectedCode: http.StatusAccepted,
	},
	{
		name:         "test articles handler without api key",
		url:          "/v1/articles",
		apiKey:       "",
		expectedCode: http.StatusForbidden,
	},
	{
		name:         "test articles handler with incorrect api key",
		url:          "/v1/articles",
		apiKey:       "test123",
		expectedCode: http.StatusForbidden,
	},
	{
		name:         "test articles handler with a correct api key",
		url:          "/v1/articles",
		apiKey:       testUsers[1].ApiKey,
		expectedCode: http.StatusAccepted,
	},
}

func (suite *TestMainPackage) TestVerifyApiKey() {
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
	err := e.fs.Set("env", "test")
	suite.Nil(err, "env flag should not throw an error")
	e.setEnvironment(&cfg)
	setDSN(&cfg)
	err = setServerPort(&cfg)
	suite.Nil(err, "should not have error")
	suite.Equal("test", cfg.env)
	db, err := driver.ConnectSQL(cfg.db.dsn, a)
	if err != nil {
		suite.Fail("could not connect to the test sql")
	}
	handlerRepo := handlers.NewRepo(&a, db)
	handlers.NewHandlers(handlerRepo)
	hr.testHandlerRepo = handlerRepo
	hr.testSQL = db
	err = hr.testSQL.SQL.AutoMigrate(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not create test tables")
	hr.testSQL.SQL.Create(&testUserGroups)
	hr.testSQL.SQL.Create(&testUsers)
	hr.testSQL.SQL.Create(&testArticles)
	hr.testSQL.SQL.Create(&testComments)

	v := routes()

	for _, t := range testsToVerifyApiKey {
		req := httptest.NewRequest(http.MethodGet, t.url, nil)
		req.Header.Set(handlers.AppApiKey, t.apiKey)
		res := httptest.NewRecorder()
		v.ServeHTTP(res, req)
		suite.Equal(t.expectedCode, res.Code, "response codes should be equal")
	}

	err = hr.testSQL.SQL.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	err = os.RemoveAll(root + "/" + hr.testHandlerRepo.App.StaticImages)
	suite.Nil(err, "could not remove test images from folder")
}
