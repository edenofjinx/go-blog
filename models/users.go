package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// User struct for users
type User struct {
	gorm.Model
	Email         string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password      string    `gorm:"type:varchar(250);not null" json:"-"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	ApiKey        string    `gorm:"type:varchar(60);index" json:"api_key"`
	PasswordReset time.Time `gorm:"default:NULL"`
	GroupID       int
	Group         UserGroup `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// UserPayload struct for user payload
type UserPayload struct {
	ID       uint   `json:"-"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserUpdatePayload struct {
	UserPayload
	Password string `json:"-"`
	ID       uint   `json:"id"`
}

type UserPasswordPayload struct {
	UserPayload
	Email string `json:"-"`
	Name  string `json:"-"`
	ID    uint   `json:"id"`
}

type UserGroupPayload struct {
	UserPayload
	Email    string `json:"-"`
	Name     string `json:"-"`
	Password string `json:"-"`
	ID       uint   `json:"id"`
	GroupID  int    `json:"group_id"`
}

type UserLoginPayload struct {
	UserPayload
	Name string `json:"-"`
}

type UserLoginResponse struct {
	ApiKey string `json:"api_key"`
}

type UserKeyPayload struct {
	UserPayload
	Email    string `json:"-"`
	Name     string `json:"-"`
	Password string `json:"-"`
	ID       uint   `json:"id"`
}

type UserKeyUpdateResponse struct {
	ApiKey string `json:"api_key"`
}

// TODO move hashes to separate package
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
