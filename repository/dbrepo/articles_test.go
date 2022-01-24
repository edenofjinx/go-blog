package dbrepo

import (
	"fmt"
	"strconv"
)

var testsForGetArticlesList = []struct {
	name            string
	pagination      testPagination
	expectedLength  int
	expectedTitle   string
	expectedContent string
}{
	{
		name:            "default articles list",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  2,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with limit of 1",
		pagination:      testPagination{limit: "1", page: "", order: ""},
		expectedLength:  1,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with limit of 1000",
		pagination:      testPagination{limit: "1000", page: "", order: ""},
		expectedLength:  2,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with page of 5",
		pagination:      testPagination{limit: "", page: "5", order: ""},
		expectedLength:  0,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "articles list with order of DESC",
		pagination:      testPagination{limit: "", page: "", order: "DESC"},
		expectedLength:  2,
		expectedTitle:   testArticles[1].Title,
		expectedContent: testArticles[1].Content,
	},
	{
		name:            "articles list with mixed pagination",
		pagination:      testPagination{limit: "1", page: "2", order: "DESC"},
		expectedLength:  1,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with incorrect pagination variables",
		pagination:      testPagination{limit: "test", page: "-50", order: "TestOrder"},
		expectedLength:  2,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
}

var testsForGetArticleById = []struct {
	name            string
	requestID       int
	expectedTitle   string
	expectedContent string
	testNoID        bool
}{
	{
		name:            "get article with id 1",
		requestID:       1,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
		testNoID:        false,
	},
	{
		name:            "get article with id 2",
		requestID:       2,
		expectedTitle:   testArticles[1].Title,
		expectedContent: testArticles[1].Content,
		testNoID:        false,
	},
	{
		name:            "get article with non existing id",
		requestID:       99,
		expectedTitle:   "",
		expectedContent: "",
		testNoID:        false,
	},
	{
		name:            "get article with no id",
		requestID:       -500,
		expectedTitle:   "",
		expectedContent: "",
	},
}

func (suite *databaseRequestTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		c := generateNewGETRequest("/v1/articles", t.pagination)
		list, err := suite.testRepo.GetArticlesList(c)
		if err != nil {
			suite.Fail(fmt.Sprintf("failed to load articles in test name %s", t.name))
		}
		l := len(list)
		suite.Equal(
			t.expectedLength,
			l,
			fmt.Sprintf("expected length %d is not equal to actual length %d in test name %s", t.expectedLength, l, t.name),
		)
		if t.expectedLength == 0 && l == t.expectedLength {
			continue
		}
		suite.Equal(
			t.expectedTitle,
			list[0].Title,
			fmt.Sprintf("expected and actual title are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			list[0].Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}

func (suite *databaseRequestTestSuite) TestGetArticleById() {
	for _, t := range testsForGetArticleById {
		ar, err := suite.testRepo.GetArticleById(t.requestID)
		if err != nil {
			if t.requestID == 0 {
				continue
			}
			suite.Fail(fmt.Sprintf("error when loading article by id %d, in test name %s", t.requestID, t.name))
		}
		suite.Equal(
			t.expectedTitle,
			ar.Title,
			fmt.Sprintf("expected and actual title are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			ar.Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}

func (suite *databaseRequestTestSuite) TestGetCommentsByArticleId() {
	for _, t := range testsForGetCommentsByArticleId {
		c := generateNewGETRequest("/v1/article/:id/comments", t.pagination)
		artId, err := strconv.Atoi(t.articleId)
		suite.Nil(err, "should not be an error")
		list, err := suite.testRepo.GetCommentsByArticleId(artId, c.Request)
		if err != nil {
			if t.articleId == "" {
				continue
			}
			suite.Fail(fmt.Sprintf("failed to load articles in test name %s", t.name))
		}
		l := len(list)
		suite.Equal(
			t.expectedLength,
			l,
			fmt.Sprintf("expected length %d is not equal to actual length %d in test name %s", t.expectedLength, l, t.name),
		)
		if t.expectedLength == 0 && l == t.expectedLength {
			continue
		}
		suite.Equal(
			t.expectedUserID,
			list[0].UserID,
			fmt.Sprintf("expected and actual user id are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			list[0].Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}
