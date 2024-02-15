package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SawitProRecruitment/UserService/handler/httpcontext"
	"github.com/SawitProRecruitment/UserService/handler/middleware"
	_ "github.com/SawitProRecruitment/UserService/lib/errors"
	"github.com/SawitProRecruitment/UserService/service"
	"github.com/labstack/echo/v4"
)

func bindAndValidate(c echo.Context, r any) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

// @Summary Get user
// @Description Get user data
// @Router /v1/user [get]
// @Produce json
// @Param Authorization header string true "Bearer"
// @Success 200 {object} responseWithData
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security ApiKeyAuth
func (s *Server) GetCurrentUser(c echo.Context) error {
	err := middleware.Auth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	userJwt, ok := httpcontext.GetUserJWT(c)
	if !ok {
		return fmt.Errorf("cannot get user from context")
	}
	u, err := s.Service.GetByID(ctx, userJwt.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newSuccessGetByIDResponse(u))
}

// @Summary Register user
// @Description Register new user
// @Router /v1/users [post]
// @Produce json
// @Param name body string true "Name"
// @Param phone body string true "Phone Number"
// @Param password body string true "Password"
// @Success 200 {object} responseWithData
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
func (s *Server) Register(c echo.Context) error {
	ctx := c.Request().Context()
	var payload service.PayloadInsert
	if err := bindAndValidate(c, &payload); err != nil {
		return err
	}
	id, err := s.Service.InsertUser(ctx, payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newSuccessRegisterResponse(id))
}

// @Summary Update user
// @Description Update user data
// @Router /v1/user [patch]
// @Produce json
// @Param Authorization header string true "Bearer"
// @Param name body string true "Name"
// @Param phone body string true "Phone Number"
// @Success 200 {object} baseResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
func (s *Server) UpdateProfile(c echo.Context) error {
	err := middleware.Auth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var payload service.PayloadUpdate
	payload.Id = id
	if err := bindAndValidate(c, &payload); err != nil {
		return err
	}
	err = s.Service.UpdateProfile(ctx, payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newBaseResponse("Successfully update user!"))
}

// @Summary Login user
// @Description Login user
// @Router /v1/users/login [post]
// @Produce json
// @Param phone body string true "Phone Number"
// @Param password body string true "Password"
// @Success 200 {object} responseWithData
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
func (s *Server) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var payload service.PayloadLogin
	if err := bindAndValidate(c, &payload); err != nil {
		return err
	}
	res, err := s.Service.Login(ctx, payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newSuccessLogin(res))
}
