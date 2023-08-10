package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	MaxPhoneNumber       = 13
	MinPhoneNumber       = 10
	IndonesiaCountryCode = "+62"

	MaxFullName = 60
	MinFullName = 3

	MinPassword = 6
	MaxPassword = 64
)

// Register handles user registration
func (s *Server) Register(ctx echo.Context) error {
	var (
		errorResp   generated.ErrorResponse
		successResp generated.RegistrationResponse
	)

	request := &generated.RegistrationRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return err
	}

	err = validate(request)
	if err != nil {
		errorResp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errorResp)
	}

	// hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	request.Password = string(hashedPwd)

	// set user id
	userID := uuid.New().String()

	// store to database
	registrationData := &repository.User{
		ID:          userID,
		FullName:    request.FullName,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	}

	err = s.Repository.StoreRegistration(ctx.Request().Context(), registrationData)
	if err != nil {
		return err
	}

	successResp.UserId = userID
	return ctx.JSON(http.StatusOK, successResp)
}

func validate(request *generated.RegistrationRequest) error {
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

	// validate password
	err = validatePassword(request.Password)
	if err != nil {
		errStrs = append(errStrs, err.Error())
	}

	return helper.ErrStringsToErr(errStrs)
}

func validatePhoneNumber(phoneNumber string) error {
	var errStr []string
	phoneNumberLen := len(phoneNumber)

	// numeric validation
	if phoneNumberLen > 1 {
		phoneNumeric := phoneNumber[1:]
		_, err := strconv.ParseInt(phoneNumeric, 10, 64)
		if err != nil {
			errStr = append(errStr, "not a phone number")
		}
	}

	// character length
	if phoneNumberLen < MinPhoneNumber || phoneNumberLen > MaxPhoneNumber {
		errStr = append(errStr, "must be at minimum 10 characters and maximum 13 characters")
	}

	// country code
	if phoneNumberLen > 3 {
		countryCode := phoneNumber[:3]
		if countryCode != IndonesiaCountryCode {
			errStr = append(errStr, "country code is not valid")
		}
	}

	if len(errStr) > 0 {
		return errors.New(fmt.Sprintf("phone_number: %s", strings.Join(errStr, ",")))
	}

	return nil
}

func validateFullName(fullName string) error {
	var errStr []string
	fullNameLen := len(fullName)

	// character length
	if fullNameLen < MinFullName || fullNameLen > MaxFullName {
		errStr = append(errStr, "must be at minimum 3 characters and maximum 60 characters")
	}

	if len(errStr) > 0 {
		return errors.New(fmt.Sprintf("full_name: %s", strings.Join(errStr, ",")))
	}

	return nil
}

func validatePassword(password string) error {
	var errStr []string
	passwordLen := len(password)
	var number, upper, special bool

	// character length
	if passwordLen < MinPassword || passwordLen > MaxPassword {
		errStr = append(errStr, "must be at minimum 6 characters and maximum 64 characters")
	}

	// password combination
	// refs: https://stackoverflow.com/questions/25837241/password-validation-with-regexp
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}

	if !number || !upper || !special {
		errStr = append(errStr, "must contain at least 1 capital characters, 1 number, and 1 special character")
	}

	if len(errStr) > 0 {
		return errors.New(fmt.Sprintf("password: %s", strings.Join(errStr, ",")))
	}

	return nil
}
