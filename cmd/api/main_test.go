package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func (suite *TestMainPackage) TestSetServerPort() {
	var e envSet
	var testCfg serverConfig
	e.setFlag()
	e.fs.StringVar(&testCfg.env, "env", "test", "testing")
	e.setEnvironment(&testCfg)
	err := setServerPort(&testCfg)
	suite.Nil(err, "should not have error")
	suite.Equal("test", testCfg.env)
	suite.Equal(4000, testCfg.port, "incorrect port selected")

	e.setFlag()
	e.fs.StringVar(&testCfg.env, "env", "development", "development")
	e.setEnvironment(&testCfg)
	os.Setenv("APP_PORT", "test")
	err = setServerPort(&testCfg)
	suite.Error(err, "should have error")
}

func (suite *TestMainPackage) TestSetDSN() {
	var e envSet
	var testCfg serverConfig
	e.setFlag()
	e.fs.StringVar(&testCfg.env, "env", "test", "testing")
	e.setEnvironment(&testCfg)
	setDSN(&testCfg)
	suite.Equal("test", testCfg.env)
	suite.Equal("root:root@tcp(127.0.0.1:3306)/go_blog_test?charset=utf8&parseTime=True&loc=Local", testCfg.db.dsn, "incorrect dsn")
}

func (suite *TestMainPackage) TestSetEnvironment() {
	var e envSet
	var testCfg serverConfig
	e.setFlag()
	e.setEnvironment(&testCfg)
	suite.Equal("development", testCfg.env)
	e.setFlag()
	e.fs.StringVar(&testCfg.env, "env", "test", "testing")
	e.setEnvironment(&testCfg)
	suite.Equal("test", testCfg.env)
	e.setFlag()
	e.fs.StringVar(&testCfg.env, "env", "development", "development")
	e.setEnvironment(&testCfg)
	suite.Equal("development", testCfg.env)
	e.setFlag()
	e.fs.StringVar(&testCfg.env, "env", "production", "production")
	e.setEnvironment(&testCfg)
	suite.Equal("production", testCfg.env)
	e.setFlag()
	e.setEnvironment(&testCfg)
	suite.Equal("production", testCfg.env)
}

func (suite *TestMainPackage) TestLoadEnvFile() {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	err := loadEnvFile("test", false)
	suite.Error(err, "should be an error")
	err = loadEnvFile(root+"/.env.test", false)
	suite.Nil(err, "should be nil")
}

func (suite *TestMainPackage) TestSetEnvCfg() {
	setEnvCfg()
}

func (suite *TestMainPackage) TestSetAppCfg() {
	setAppCfg()
}

func (suite *TestMainPackage) TestSetupDatabase() {
	_, err := setupDatabase()
	suite.Nil(err, "should not contain error")
}

func (suite *TestMainPackage) TestCreateHandlers() {
	setEnvCfg()
	setAppCfg()
	db, err := setupDatabase()
	suite.Nil(err, "should not contain error")
	createHandlers(db)
}

func (suite *TestMainPackage) TestCreateServer() {
	setEnvCfg()
	setAppCfg()
	_ = createServer()
}
