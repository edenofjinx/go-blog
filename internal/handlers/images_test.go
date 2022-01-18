package handlers

import "log"

func (suite *handlersTestSuite) TestParseImageTags() {
	fooString := "\xff"
	barString := string([]rune(fooString))
	_, err := suite.testHandlerRepo.parseImageTags(barString)
	log.Println("*******")
	log.Println(err)
}

func (suite *handlersTestSuite) TestSaveImages() {

}

//func (suite *handlersTestSuite) TestOpenImageFile() {
//	for _, t := range testsForOpenImageFile {
//		suite.testHandlerRepo.App.StaticImages = t.imageBasePath
//		_, err := suite.testHandlerRepo.openImageFile(t.imageType, "labas")
//		suite.Equal(t.expectedError, err, "should both be equal")
//		suite.testHandlerRepo.App.StaticImages = "static/test/images/"
//	}
//}

//func (suite *handlersTestSuite) TestGetStaticImageDir() {
//	for _, t := range testsForGetStaticImageDir {
//		suite.testHandlerRepo.App.StaticImages = t.imageBasePath
//		_, err := suite.testHandlerRepo.getStaticImageDir()
//		suite.Equal(t.expectedError, err, "should both be equal")
//	}
//}
