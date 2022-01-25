package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
	"strconv"
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

var testsForUpdateComment = []struct {
	name          string
	jsonData      string
	expectedError bool
}{
	{
		name:          "update comment",
		jsonData:      "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"user_id\": 1, \"id\": 1}",
		expectedError: false,
	},
	{
		name:          "update comment with empty data",
		jsonData:      "{\"content\": \"\",\"title\": \"\"}",
		expectedError: true,
	},
	{
		name:          "update article with incorrect data",
		jsonData:      "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"user_id\": 999, \"id\": 1}",
		expectedError: true,
	},
}

var testsForDeleteComment = []struct {
	name          string
	articleId     int
	expectedError bool
}{
	{
		name:          "delete comment with incorrect id",
		articleId:     9999,
		expectedError: false,
	},
	{
		name:          "delete comment with 0 id",
		articleId:     0,
		expectedError: false,
	},
	{
		name:          "delete comment",
		articleId:     1,
		expectedError: false,
	},
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

func (suite *databaseRequestTestSuite) TestSaveComment() {
	for _, t := range testsForInsertComment {
		var tc models.Comment
		tc.ArticleID = t.articleID
		tc.UserID = t.userID
		tc.Content = t.content
		tc.CreatedAt = time.Now()
		tc.UpdatedAt = time.Now()
		err := suite.testRepo.SaveComment(tc)
		if err != nil {
			if t.expectedError {
				continue
			}
			suite.Fail(fmt.Sprintf("error not expected in test name %s", t.name))
		}
	}
}

func (suite *databaseRequestTestSuite) TestUpdateComment() {
	for _, t := range testsForUpdateComment {
		var payload models.CommentPayload
		var comment models.Comment
		err := json.Unmarshal([]byte(t.jsonData), &payload)
		suite.Nil(err, "should not be nil")
		comment.UserID = payload.UserID
		comment.Content = payload.Content
		comment.ID = payload.ID
		err = suite.testRepo.UpdateComment(comment)
		if t.expectedError {
			suite.Error(err, "should be an error")
		} else {
			suite.Nil(err, "should not be a nil")
		}
	}
}

func (suite *databaseRequestTestSuite) TestDeleteComment() {
	for _, t := range testsForDeleteComment {
		err := suite.testRepo.DeleteComment(t.articleId)
		if t.expectedError {
			suite.Error(err, "should be an error")
		} else {
			suite.Nil(err, "should not be a nil")
		}
	}
}
