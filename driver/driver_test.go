package driver

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var testDsn = "root:root@tcp(localhost:3306)/go_blog_test?charset=utf8&parseTime=True&loc=Local"
var a config.AppConfig

func TestConnectSQL(t *testing.T) {
	d, err := ConnectSQL(testDsn, a)
	assert.Nil(t, err)
	assert.ObjectsAreEqualValues(gorm.DB{}, d)
}

func TestNewDatabase(t *testing.T) {
	d, err := NewDatabase(testDsn, a)
	assert.Nil(t, err)
	assert.ObjectsAreEqualValues(gorm.DB{}, d)
}
