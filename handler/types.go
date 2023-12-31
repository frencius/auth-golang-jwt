package handler

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var PrivateKey = os.Getenv("PRIVATE_KEY")

type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	jwt.RegisteredClaims
}

var CompareHashAndPassword = bcrypt.CompareHashAndPassword
var GenerateFromPassword = bcrypt.GenerateFromPassword
var ParseWithClaims = jwt.ParseWithClaims
