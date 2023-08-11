package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)
	srv := Server{
		Repository: mockRepository,
	}

	mockAuthHeader := "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9"
	profileParams := generated.GetProfileParams{
		Authorization: mockAuthHeader,
	}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	t.Run("positive", func(t *testing.T) {
		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		mockRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&repository.User{}, nil).Times(1)

		err := srv.GetProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("get user by id error", func(t *testing.T) {
		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		mockRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&repository.User{}, errors.New("error")).Times(1)

		err := srv.GetProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("token invalid", func(t *testing.T) {
		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: false,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		err := srv.GetProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("parse with claims error", func(t *testing.T) {
		err := srv.GetProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("auth token invalid", func(t *testing.T) {
		profileParams := generated.GetProfileParams{
			Authorization: "mockAuthHeader",
		}
		err := srv.GetProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("signature invalid", func(t *testing.T) {
		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: false,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, jwt.ErrSignatureInvalid
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		err := srv.GetProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})
}

func TestUpdateProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)
	srv := Server{
		Repository: mockRepository,
	}

	mockAuthHeader := "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9"
	profileParams := generated.UpdateProfileParams{
		Authorization: mockAuthHeader,
	}

	t.Run("positive", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name": "hai",
				"phone_number": "+622342342322"
			}`)

		param := generated.UpdateProfileRequest{}
		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPatch, "/profile/update", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		mockRepository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("update profile error", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name": "hai",
				"phone_number": "+622342342322"
			}`)

		param := generated.UpdateProfileRequest{}
		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPatch, "/profile/update", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		mockRepository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("update profile error - user is not exist", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name": "hai",
				"phone_number": "+622342342322"
			}`)

		param := generated.UpdateProfileRequest{}
		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPatch, "/profile/update", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		mockRepository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(errors.New("user is not exist")).Times(1)

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("update profile error - phone number conflict", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name": "hai",
				"phone_number": "+622342342322"
			}`)

		param := generated.UpdateProfileRequest{}
		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPatch, "/profile/update", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		mockRepository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(errors.New("pq: duplicate key value")).Times(1)

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("error validate payload", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name": "",
				"phone_number": ""
			}`)

		param := generated.UpdateProfileRequest{}
		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPatch, "/profile/update", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("decode request error", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name": "asdas",
				"phone_number": 1
			}`)

		param := generated.UpdateProfileRequest{}
		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPatch, "/profile/update", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: true,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("signature invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/profile/update", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: false,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, jwt.ErrSignatureInvalid
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("token invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/profile/update", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		tempParseWithClaims := ParseWithClaims
		ret := &jwt.Token{
			Valid: false,
		}
		ParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return ret, nil
		}

		defer func() {
			ParseWithClaims = tempParseWithClaims
		}()

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("parse with claims error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/profile/update", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("auth token invalid", func(t *testing.T) {
		profileParams := generated.UpdateProfileParams{
			Authorization: "mockAuthHeader",
		}
		req := httptest.NewRequest(http.MethodPatch, "/profile/update", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		err := srv.UpdateProfile(c, profileParams)
		assert.Nil(t, err, "error should be nil")
	})
}
