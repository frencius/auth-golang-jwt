package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
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
		successResp generated.GetProfileResponse
	)

	authHeader := params.Authorization
	splittedAuth := strings.Split(authHeader, " ")
	if len(splittedAuth) < 2 {
		return sendErrorResponse(ctx, http.StatusBadRequest, errors.New("Auth header is not valid"))
	}

	tknStr := splittedAuth[1]
	claims := &Claims{}

	block, _ := pem.Decode([]byte(PrivateKey))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)

	tkn, err := ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key.Public(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return sendErrorResponse(ctx, http.StatusForbidden, errors.New("User is not authorized"))
		}
		return sendErrorResponse(ctx, http.StatusForbidden, err)
	}

	if !tkn.Valid {
		return sendErrorResponse(ctx, http.StatusForbidden, errors.New("User is not authorized"))
	}

	// get user data by user id
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), claims.UserID)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	successResp.FullName = user.FullName
	successResp.PhoneNumber = user.PhoneNumber

	return ctx.JSON(http.StatusOK, successResp)
}

func (s *Server) UpdateProfile(ctx echo.Context, params generated.UpdateProfileParams) error {
	var (
		successResp generated.UpdateProfileResponse
	)

	authHeader := params.Authorization
	splittedAuth := strings.Split(authHeader, " ")
	if len(splittedAuth) < 2 {
		return sendErrorResponse(ctx, http.StatusBadRequest, errors.New("Auth header is not valid"))
	}

	tknStr := splittedAuth[1]
	claims := &Claims{}

	block, _ := pem.Decode([]byte(PrivateKey))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)

	tkn, err := ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key.Public(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return sendErrorResponse(ctx, http.StatusForbidden, errors.New("User is not authorized"))
		}

		return sendErrorResponse(ctx, http.StatusForbidden, err)
	}

	if !tkn.Valid {
		return sendErrorResponse(ctx, http.StatusForbidden, errors.New("User is not authorized"))
	}

	request := &generated.UpdateProfileRequest{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = validateUpdateProfile(request)
	if err != nil {
		return sendErrorResponse(ctx, http.StatusBadRequest, err)
	}

	// update profile
	registrationData := &repository.User{
		ID:          claims.UserID,
		FullName:    request.FullName,
		PhoneNumber: request.PhoneNumber,
	}
	err = s.Repository.UpdateProfile(ctx.Request().Context(), registrationData)
	if err != nil {
		if err.Error() == "user is not exist" {
			return sendErrorResponse(ctx, http.StatusNotFound, err)
		} else if strings.Contains(err.Error(), "pq: duplicate key value") {
			return sendErrorResponse(ctx, http.StatusConflict, errors.New("phone number conflict"))
		}
		return sendErrorResponse(ctx, http.StatusInternalServerError, err)
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
