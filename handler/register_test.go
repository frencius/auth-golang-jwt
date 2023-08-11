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

func TestRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepositoryInterface(mockCtrl)
	srv := Server{
		Repository: mockRepository,
	}

	var param generated.RegistrationRequest
	t.Run("positive", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name":    "sadam 2",
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().StoreRegistration(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := srv.Register(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("store registration error", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name":    "sadam 2",
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().StoreRegistration(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

		err := srv.Register(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("store registration error - phone number conflict", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name":    "sadam 2",
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		mockRepository.EXPECT().StoreRegistration(gomock.Any(), gomock.Any()).Return(errors.New("pq: duplicate key value")).Times(1)

		err := srv.Register(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("generate password error", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name":    "sadam 2",
				"phone_number": "+622342342322",
				"password":     "AAAAAAAAA1a^1"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		tempGenerateFromPassword := GenerateFromPassword
		GenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
			return []byte{}, errors.New("error")
		}

		defer func() {
			GenerateFromPassword = tempGenerateFromPassword
		}()

		err := srv.Register(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("payload invalid", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name":    "",
				"phone_number": "+asd",
				"password":     "A"
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := srv.Register(c)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("decode request error", func(t *testing.T) {
		payload := []byte(
			`{
				"full_name":    "",
				"phone_number": "+6213",
				"password":     1
			}`)

		paramBytes := payload
		_ = json.Unmarshal(paramBytes, &param)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(paramBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		err := srv.Register(c)
		assert.Nil(t, err, "error should be nil")
	})
}
