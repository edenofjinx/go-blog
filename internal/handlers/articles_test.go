package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
)

var testsForGetArticlesList = []struct {
	name          string
	expectedCode  int
	expectedTitle string
	pagination    testPagination
}{
	{
		name:          "default articles list",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
		pagination: testPagination{
			limit: "",
			page:  "",
			order: "",
		},
	},
	{
		name:          "articles list with pagination",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
		pagination: testPagination{
			limit: "1",
			page:  "2",
			order: "DESC",
		},
	},
	{
		name:          "articles list with incorrect pagination",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
		pagination: testPagination{
			limit: "9999",
			page:  "-50",
			order: "test",
		},
	},
}

var testsForGetArticleById = []struct {
	name            string
	requestID       string
	expectedCode    int
	expectedTitle   string
	expectedContent string
}{
	{
		name:          "default article by id",
		requestID:     "1",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
	},
	{
		name:            "get article with non existing id",
		requestID:       "99",
		expectedCode:    http.StatusAccepted,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with non existing id string",
		requestID:       "test",
		expectedCode:    http.StatusBadRequest,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with no id",
		expectedCode:    http.StatusBadRequest,
		requestID:       "",
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with no string as id",
		expectedCode:    http.StatusBadRequest,
		requestID:       "test",
		expectedTitle:   "",
		expectedContent: "",
	},
}

type testsForSaveArticle struct {
	Tests []testForSaveArticle `json:"tests"`
}

type testForSaveArticle struct {
	Name         string                `json:"name"`
	JsonData     models.ArticlePayload `json:"json_data"`
	ExpectedCode int                   `json:"expected_code"`
}

func (suite *handlersTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		m := make(map[string][]models.Article)
		c, rr := generateNewGETRequest("/v1/articles/", t.pagination)
		suite.testHandlerRepo.GetArticlesList(c)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		err := json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		suite.Equal(t.expectedTitle, m["data"][0].Title, fmt.Sprintf("status code should be %s but got %s", t.expectedTitle, m["data"][0].Title))
	}
}

func (suite *handlersTestSuite) TestGetArticleById() {
	var p testPagination
	for _, t := range testsForGetArticleById {
		m := make(map[string]models.Article)
		c, rr := generateNewGETRequest("/v1/article/:id", p)
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: t.requestID,
			},
		}
		suite.testHandlerRepo.GetArticleById(c)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		if t.expectedCode == status && t.expectedContent == "" {
			continue
		}
		err := json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		suite.Equal(t.expectedTitle, m["data"].Title, fmt.Sprintf("status code should be %s but got %s", t.expectedTitle, m["data"].Title))
	}
}

func (suite *handlersTestSuite) TestSaveArticle() {
	jsonFile, err := os.Open("./testData/articles_test_save_article.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tests testsForSaveArticle
	err = json.Unmarshal(byteValue, &tests)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(tests.Tests); i++ {
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		b, err := json.Marshal(tests.Tests[i].JsonData)
		log.Println(b)
		if err != nil {
			fmt.Println(err)
			return
		}
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Path: "/v1/article/save",
			},
			Body: io.NopCloser(strings.NewReader(string(b))),
		}
		c.Request = req
		suite.testHandlerRepo.SaveArticle(c)
		status := rr.Code
		suite.Equal(tests.Tests[i].ExpectedCode, status, fmt.Sprintf("status code should be %d but got %d", tests.Tests[i].ExpectedCode, status))
	}
}
