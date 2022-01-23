package database

import "bitbucket.org/julius_liaudanskis/go-blog/models"

func (suite *seedTestSuite) TestSeedDataWithoutTables() {
	err := SeedData(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedUserGroupsWithoutTable() {
	err := seedUserGroups(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedUsersWithoutTable() {
	err := seedUsers(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedArticlesWithoutTable() {
	err := seedArticles(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedCommentsWithoutTable() {
	err := seedComments(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedUserGroups() {
	err := suite.db.SQL.AutoMigrate(&models.UserGroup{})
	suite.Nil(err, "Should not have error")
	err = seedUserGroups(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedUsers() {
	err := suite.db.SQL.AutoMigrate(&models.User{})
	suite.Nil(err, "Should not have error")
	err = seedUsers(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedArticles() {
	err := suite.db.SQL.AutoMigrate(&models.Article{})
	suite.Nil(err, "Should not have error")
	err = seedArticles(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedComments() {
	err := suite.db.SQL.AutoMigrate(&models.Comment{})
	suite.Nil(err, "Should not have error")
	err = seedComments(suite.db)
	suite.Nil(err, "Should not have error")
}
