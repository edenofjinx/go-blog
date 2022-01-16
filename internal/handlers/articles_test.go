package handlers

import "log"

var testsForGetArticlesList = []struct {
	name string
}{
	{},
}

var testsForGetArticleById = []struct {
}{
	{},
}

func (suite *handlersTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		log.Println(t)
	}
}

func (suite *handlersTestSuite) TestGetArticleById() {
	for _, t := range testsForGetArticleById {
		log.Println(t)
	}
}
