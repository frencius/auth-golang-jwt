package handler

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var PrivateKey = ``

type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	jwt.RegisteredClaims
}

var CompareHashAndPassword = bcrypt.CompareHashAndPassword
var GenerateFromPassword = bcrypt.GenerateFromPassword
var ParseWithClaims = jwt.ParseWithClaims
