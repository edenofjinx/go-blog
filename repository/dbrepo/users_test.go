package dbrepo

import "fmt"

var testsForVerifyApiKeyExists = []struct {
	name             string
	testApiKey       string
	expectedResponse bool
}{
	{
		name:             "verify with existing api key",
		testApiKey:       testUsers[0].ApiKey,
		expectedResponse: true,
	},
	{
		name:             "verify with non existing api key",
		testApiKey:       "asdf-4444-test-123",
		expectedResponse: false,
	},
	{
		name:             "verify with empty api key",
		testApiKey:       "",
		expectedResponse: false,
	},
}

func (suite *databaseRequestTestSuite) TestVerifyApiKeyExists() {
	for _, t := range testsForVerifyApiKeyExists {
		v := suite.testRepo.VerifyApiKeyExists(t.testApiKey)
		suite.Equal(t.expectedResponse, v, fmt.Sprintf("expected and actual responses are not equal in test name %s", t.name))
	}
}
