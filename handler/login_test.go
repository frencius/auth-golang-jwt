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
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)
	srv := Server{
		Repository: mockRepository,
	}
	var param generated.LoginRequest

	mockUser := repository.User{
		FullName: "sadam",
	}
	t.Run("positive", func(t *testing.T) {
		payload := []byte(
			`{
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&mockUser, nil).Times(1)

		tempCompareHashAndPassword := CompareHashAndPassword
		CompareHashAndPassword = func(hashedPassword, password []byte) error {
			return nil
		}

		defer func() {
			CompareHashAndPassword = tempCompareHashAndPassword
		}()

		mockRepository.EXPECT().UpdateLogin(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := srv.Login(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("update login error", func(t *testing.T) {
		payload := []byte(
			`{
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&mockUser, nil).Times(1)

		tempCompareHashAndPassword := CompareHashAndPassword
		CompareHashAndPassword = func(hashedPassword, password []byte) error {
			return nil
		}

		defer func() {
			CompareHashAndPassword = tempCompareHashAndPassword
		}()

		mockRepository.EXPECT().UpdateLogin(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

		err := srv.Login(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("decode request error", func(t *testing.T) {
		payload := []byte(
			`{
				"phone_number": "+622342342322",
				"password":     1
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := srv.Login(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("payload validation error - empty", func(t *testing.T) {
		payload := []byte(
			`{
				"phone_number": "",
				"password":     ""
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := srv.Login(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("CompareHashAndPassword error", func(t *testing.T) {
		payload := []byte(
			`{
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&mockUser, nil).Times(1)

		err := srv.Login(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("GetUser error", func(t *testing.T) {
		payload := []byte(
			`{
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&mockUser, errors.New("error")).Times(1)

		err := srv.Login(c)
		assert.Nil(t, err, "error should be nil")
	})

}
