package database

func (suite *seedTestSuite) TestMigrateData() {
	err := MigrateData(suite.db)
	suite.Nil(err, "Should not have error")
}
