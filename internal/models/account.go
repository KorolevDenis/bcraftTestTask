package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	Id       uint   `json:"id" gorm:"primarykey"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token" ;sql:"-"`
}
