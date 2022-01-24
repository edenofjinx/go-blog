package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"fmt"
	"time"
)

var testsForInsertComment = []struct {
	name          string
	articleID     int
	userID        int
	content       string
	expectedError bool
}{
	{
		name:          "insert comment with correct data",
		articleID:     1,
		userID:        1,
		content:       "some new content 1",
		expectedError: false,
	},
	{
		name:          "insert comment with incorrect article id",
		articleID:     99,
		userID:        1,
		content:       "some new content 99",
		expectedError: true,
	},
	{
		name:          "insert comment with incorrect user id",
		articleID:     1,
		userID:        99,
		content:       "some new content user 99",
		expectedError: true,
	},
	{
		name:          "insert comment with empty content",
		articleID:     2,
		userID:        2,
		content:       "",
		expectedError: false,
	},
	{
		name:          "insert comment with empty article id",
		userID:        2,
		content:       "empty article id",
		expectedError: true,
	},
	{
		name:          "insert comment with empty user id",
		articleID:     1,
		content:       "empty user id",
		expectedError: true,
	},
}

var testsForGetCommentsByArticleId = []struct {
	name            string
	articleId       string
	pagination      testPagination
	expectedLength  int
	expectedContent string
	expectedUserID  int
}{
	{
		name:            "default comments list",
		articleId:       "1",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  2,
		expectedContent: testComments[0].Content,
		expectedUserID:  testComments[0].UserID,
	},
	{
		name:            "comments list with limit of 1",
		articleId:       "1",
		pagination:      testPagination{limit: "1", page: "", order: ""},
		expectedLength:  1,
		expectedContent: testComments[0].Content,
		expectedUserID:  testComments[0].UserID,
	},
	{
		name:            "comments list with page of 5",
		articleId:       "1",
		pagination:      testPagination{limit: "", page: "5", order: ""},
		expectedLength:  0,
		expectedContent: "",
	},
	{
		name:            "comments list with order of DESC",
		articleId:       "2",
		pagination:      testPagination{limit: "", page: "", order: "DESC"},
		expectedLength:  2,
		expectedContent: testComments[3].Content,
		expectedUserID:  testComments[3].UserID,
	},
	{
		name:            "comments list with mixed pagination",
		articleId:       "2",
		pagination:      testPagination{limit: "1", page: "2", order: "DESC"},
		expectedLength:  1,
		expectedContent: testComments[2].Content,
		expectedUserID:  testComments[2].UserID,
	},
	{
		name:            "comments list with incorrect pagination variables",
		articleId:       "2",
		pagination:      testPagination{limit: "test", page: "-50", order: "TestOrder"},
		expectedLength:  2,
		expectedContent: testComments[2].Content,
		expectedUserID:  testComments[2].UserID,
	},
	{
		name:            "comments list with non existing article id",
		articleId:       "99",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  0,
		expectedContent: "",
	},
	{
		name:            "comments list with incorrect article id",
		articleId:       "-500",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  0,
		expectedContent: "",
	},
}

func (suite *databaseRequestTestSuite) TestInsertComment() {
	for _, t := range testsForInsertComment {
		var tc models.Comment
		tc.ArticleID = t.articleID
		tc.UserID = t.userID
		tc.Content = t.content
		tc.CreatedAt = time.Now()
		tc.UpdatedAt = time.Now()
		err := suite.testRepo.InsertComment(tc)
		if err != nil {
			if t.expectedError {
				continue
			}
			suite.Fail(fmt.Sprintf("error not expected in test name %s", t.name))
		}
	}
}
