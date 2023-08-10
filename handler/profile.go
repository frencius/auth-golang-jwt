package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetProfile(ctx echo.Context, params generated.GetProfileParams) error {
	var (
		errorResp   generated.ErrorResponse
		successResp generated.GetProfileResponse
	)

	authHeader := ctx.Request().Header.Get("Authorization")
	splittedAuth := strings.Split(authHeader, " ")
	if len(splittedAuth) < 2 {
		errorResp.Message = "Auth header is not valid"
		return ctx.JSON(http.StatusBadRequest, errorResp)
	}

	tknStr := splittedAuth[1]
	claims := &Claims{}

	block, _ := pem.Decode([]byte(PrivateKey))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key.Public(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			errorResp.Message = "User is not authorized"
			return ctx.JSON(http.StatusForbidden, errorResp)
		}

		errorResp.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errorResp)
	}

	if !tkn.Valid {
		errorResp.Message = "User is not authorized"
		return ctx.JSON(http.StatusForbidden, errorResp)
	}

	// get user data by user id
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), claims.UserID)
	if err != nil {
		errorResp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errorResp)
	}

	successResp.FullName = user.FullName
	successResp.PhoneNumber = user.PhoneNumber

	return ctx.JSON(http.StatusOK, successResp)
}

func (s *Server) UpdateProfile(ctx echo.Context, params generated.UpdateProfileParams) error {
	var (
		errorResp   generated.ErrorResponse
		successResp generated.UpdateProfileResponse
	)

	authHeader := ctx.Request().Header.Get("Authorization")
	splittedAuth := strings.Split(authHeader, " ")
	if len(splittedAuth) < 2 {
		errorResp.Message = "Auth header is not valid"
		return ctx.JSON(http.StatusBadRequest, errorResp)
	}

	tknStr := splittedAuth[1]
	claims := &Claims{}

	block, _ := pem.Decode([]byte(PrivateKey))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key.Public(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			errorResp.Message = "User is not authorized"
			return ctx.JSON(http.StatusForbidden, errorResp)
		}

		errorResp.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errorResp)
	}

	if !tkn.Valid {
		errorResp.Message = "User is not authorized"
		return ctx.JSON(http.StatusForbidden, errorResp)
	}

	request := &generated.UpdateProfileRequest{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return err
	}

	err = validateUpdateProfile(request)
	if err != nil {
		errorResp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errorResp)
	}

	// update profile
	registrationData := &repository.User{
		ID:          claims.UserID,
		FullName:    request.FullName,
		PhoneNumber: request.PhoneNumber,
	}
	err = s.Repository.UpdateProfile(ctx.Request().Context(), registrationData)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "user is not exist" {
			errorResp.Message = err.Error()
			return ctx.JSON(http.StatusBadRequest, errorResp)
		} else if strings.Contains(err.Error(), "pq: duplicate key value") {
			errorResp.Message = "phone number conflict"
			return ctx.JSON(http.StatusConflict, errorResp)
		}

		return err
	}

	successResp.Result = "update profile success"
	return ctx.JSON(http.StatusOK, successResp)
}

func validateUpdateProfile(request *generated.UpdateProfileRequest) error {
	errStrs := []string{}

	// validate phone number
	err := validatePhoneNumber(request.PhoneNumber)
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}

	// validate full name
	err = validateFullName(request.FullName)
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}

	return helper.ErrStringsToErr(errStrs)
}
