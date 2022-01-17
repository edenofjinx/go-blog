package database

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSuiteSeedData(t *testing.T) {
	suite.Run(t, new(seedTestSuite))
}

func (suite *seedTestSuite) TestSeedData() {
	err := SeedData(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedUserGroups() {
	err := seedUserGroups(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedUsers() {
	err := seedUsers(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedArticles() {
	err := seedArticles(suite.db)
	suite.Nil(err, "Should not have error")
}

func (suite *seedTestSuite) TestSeedComments() {
	err := seedComments(suite.db)
	suite.Nil(err, "Should not have error")
}
