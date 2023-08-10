package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
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

	successResp.FullName = claims.FullName
	successResp.PhoneNumber = claims.PhoneNumber

	return ctx.JSON(http.StatusOK, successResp)
}
