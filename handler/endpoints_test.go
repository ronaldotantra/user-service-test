package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/lib/errors"
	"github.com/SawitProRecruitment/UserService/lib/jwt"
	"github.com/SawitProRecruitment/UserService/lib/validator"
	"github.com/SawitProRecruitment/UserService/service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type component struct {
	ctx       context.Context
	service   *service.MockServiceInterface
	handler   *Server
	mockedErr error
	jwt       string
}

func setupService(t *testing.T) *component {
	g := gomock.NewController(t)
	service := service.NewMockServiceInterface(g)
	token, _ := jwt.GenerateToken(jwt.User{
		ID:    1,
		Name:  "rotan",
		Phone: "+62123456789",
	})

	return &component{
		ctx: context.Background(),
		handler: NewServer(NewServerOptions{
			Service: service,
		}),
		service:   service,
		mockedErr: fmt.Errorf("mocked error"),
		jwt:       string(*token),
	}
}
func TestServer_GetCurrentUser(t *testing.T) {
	t.Parallel()

	t.Run("error get user because no auth", func(t *testing.T) {
		s := setupService(t)
		req := httptest.NewRequest(http.MethodGet, "/url", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := s.handler.GetCurrentUser(c)
		assert.NotNil(t, err)
		assert.EqualError(t, err, errors.NewForbiddenError("unauthorized").Error())
	})

	t.Run("error get user because invalid token", func(t *testing.T) {
		s := setupService(t)
		req := httptest.NewRequest(http.MethodGet, "/url", nil)
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDc4ODgyOTMsImp0aSI6IjE3MDc4ODgyOTMwMjAzMDcwMDAiLCJpYXQiOjE3MDc4ODgyOTMsImlzcyI6IlNhd2l0UHJvIFVzZXIgU2VydmljZSIsInN1YiI6IkF1dGgiLCJ1c2VyIjp7ImlkIjoxMiwibmFtZSI6IlJvdGFuIiwicGhvbmUiOiIrNjI4MTM4MjUyMDAyMjIifX0.ENeN68G2QsNrD9MIXKQwilSJ9RaoSmO3GHLPBnSw0lU"
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := s.handler.GetCurrentUser(c)
		assert.NotNil(t, err)
		assert.EqualError(t, err, errors.NewForbiddenError("unauthorized").Error())
	})

	t.Run("successfully get user", func(t *testing.T) {
		s := setupService(t)
		req := httptest.NewRequest(http.MethodGet, "/url", nil)
		req.Header.Set("Authorization", "Bearer "+s.jwt)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		s.service.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&service.User{Id: 1, Name: "rotan", Phone: "+62123456789"}, nil)

		err := s.handler.GetCurrentUser(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
func TestServer_Register(t *testing.T) {
	t.Parallel()

	t.Run("success register user", func(t *testing.T) {
		s := setupService(t)
		id := int64(1)
		payload := service.PayloadInsert{
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: "Password123!",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		s.service.EXPECT().InsertUser(gomock.Any(), payload).Return(&id, nil)

		err := s.handler.Register(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid password", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadInsert{
			Name:     "rotan",
			Phone:    "+628123456789",
			Password: "Password123",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.Register(c)
		assert.NotNil(t, err)
	})

	t.Run("invalid name", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadInsert{
			Name:     "r",
			Phone:    "+628123456789",
			Password: "Password123!",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.Register(c)
		assert.NotNil(t, err)
	})

	t.Run("invalid phone", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadInsert{
			Name:     "rotan",
			Phone:    "+28123456789",
			Password: "Password123!",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.Register(c)
		assert.NotNil(t, err)
	})
}

func TestServer_UpdateProfile(t *testing.T) {
	t.Parallel()

	t.Run("success update user", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadUpdate{
			Name:  "rotan",
			Phone: "+628123456789",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", "Bearer "+s.jwt)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		s.service.EXPECT().UpdateProfile(gomock.Any(), payload).Return(nil)

		err := s.handler.UpdateProfile(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid name", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadUpdate{
			Name:  "r",
			Phone: "+628123456789",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", "Bearer "+s.jwt)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.UpdateProfile(c)
		assert.NotNil(t, err)
	})

	t.Run("invalid phone", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadUpdate{
			Name:  "rotan",
			Phone: "+28123456789",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", "Bearer "+s.jwt)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.UpdateProfile(c)
		assert.NotNil(t, err)
	})

}

func TestServer_Login(t *testing.T) {
	t.Parallel()

	t.Run("success login", func(t *testing.T) {
		s := setupService(t)
		id := int64(1)
		payload := service.PayloadLogin{
			Phone:    "+628123456789",
			Password: "Password123!",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		s.service.EXPECT().Login(gomock.Any(), payload).Return(&service.ResponseLogin{
			UserId: id,
			Token:  s.jwt,
		}, nil)

		err := s.handler.Login(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid password", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadLogin{
			Phone:    "+628123456789",
			Password: "Password123",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.Login(c)
		assert.NotNil(t, err)
	})

	t.Run("invalid phone", func(t *testing.T) {
		s := setupService(t)
		payload := service.PayloadLogin{
			Phone:    "+28123456789",
			Password: "Password123!",
		}
		bs, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Echo().Validator = validator.NewValidator()

		err := s.handler.Login(c)
		assert.NotNil(t, err)
	})
}
