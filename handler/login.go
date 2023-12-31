package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func (s *Server) Login(ctx echo.Context) error {
	var (
		successResp generated.LoginResponse
	)

	request := &generated.LoginRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = validateLogin(request)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusBadRequest, err)
	}

	// get user by phone_number
	user, err := s.Repository.GetUser(ctx.Request().Context(), request.PhoneNumber)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusNotFound, errors.New("User is not exist"))
	}

	hashedPassword := []byte(user.Password)
	err = CompareHashAndPassword(hashedPassword, []byte(request.Password))
	if err != nil {
		return sendErrorResponse(ctx, http.StatusForbidden, errors.New("Password is not valid"))
	}

	// create JWT token
	token, err := createJWTToken(user)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	// increment login counter
	err = s.Repository.UpdateLogin(ctx.Request().Context(), user.ID)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	successResp.UserId = user.ID
	successResp.Token = token
	return ctx.JSON(http.StatusOK, successResp)
}

func createJWTToken(user *repository.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	block, _ := pem.Decode([]byte(PrivateKey))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateLogin(request *generated.LoginRequest) error {
	errStrs := []string{}

	// validate phone number
	err := validateLoginPhoneNumber(request.PhoneNumber)
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}

	// validate password
	err = validateLoginPassword(request.Password)
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}

	return helper.ErrStringsToErr(errStrs)
}

func validateLoginPhoneNumber(phoneNumber string) error {
	var errStr []string
	phoneNumberLen := len(phoneNumber)

	if phoneNumberLen < 1 {
		errStr = append(errStr, "can't be empty")
	}

	// numeric validation
	_, err := strconv.ParseInt(phoneNumber, 10, 64)
	if err != nil {
		errStr = append(errStr, "not a phone number")
	}

	if len(errStr) > 0 {
		return errors.New(fmt.Sprintf("phone_number: %s", strings.Join(errStr, ",")))
	}

	return nil
}

func validateLoginPassword(password string) error {
	var errStr []string
	passwordLen := len(password)

	// password length
	if passwordLen < 1 {
		errStr = append(errStr, "can't be empty")
	}

	if len(errStr) > 0 {
		return errors.New(fmt.Sprintf("password: %s", strings.Join(errStr, ",")))
	}

	return nil
}
