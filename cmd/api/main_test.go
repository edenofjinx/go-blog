package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func (suite *TestMainPackage) TestSetServerPort() {
	e.setFlag()
	err := e.fs.Set("env", "test")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	err = setServerPort(&cfg)
	suite.Nil(err, "should not have error")

	suite.Equal("test", cfg.env)
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	suite.Nil(err, "incorrect port provided in .env file")
	suite.Equal(port, cfg.port, "incorrect port selected")

	e.setFlag()
	err = e.fs.Set("env", "development")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	defaultAppPort := os.Getenv("APP_PORT")
	os.Setenv("APP_PORT", "test")
	err = setServerPort(&cfg)
	suite.Error(err, "should have error")
	os.Setenv("APP_PORT", defaultAppPort)
}

func (suite *TestMainPackage) TestSetDSN() {
	e.setFlag()
	err := e.fs.Set("env", "test")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	setDSN(&cfg)
	suite.Equal("test", cfg.env)
	eDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	suite.Equal(eDsn, cfg.db.dsn, "incorrect dsn")
}

func (suite *TestMainPackage) TestSetEnvironment() {
	e.setFlag()
	err := e.fs.Set("env", "development")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	suite.Equal("development", cfg.env)

	e.setFlag()
	err = e.fs.Set("env", "test")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	suite.Equal("test", cfg.env)

	e.setFlag()
	err = e.fs.Set("env", "development")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	suite.Equal("development", cfg.env)

	e.setFlag()
	err = e.fs.Set("env", "production")
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	suite.Equal("production", cfg.env)

	e.setFlag()
	suite.Nil(err, "env variable set should not throw an error")
	e.setEnvironment(&cfg)
	suite.Equal("production", cfg.env)
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
	e.setFlag()
	e.fs.Set("migrate", "true")
	e.fs.Set("seed", "true")
	e.parseEnvFlag()
	e.setEnvironment(&cfg)
	setDSN(&cfg)
	err := setServerPort(&cfg)
	suite.Nil(err, "should not contain error")
	setAppCfg()
	_, err = setupDatabase()
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

func (suite *TestMainPackage) TestParseEnvFlag() {
	e.setFlag()
	os.Args = append(os.Args, "envSet")
	os.Args[2] = "-env=test"
	e.parseEnvFlag()

	e.setFlag()
	os.Args = append(os.Args, "envSet")
	os.Args[2] = "-test=new"
	e.parseEnvFlag()
}
