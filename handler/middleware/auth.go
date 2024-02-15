package middleware

import (
	"strings"

	"github.com/SawitProRecruitment/UserService/handler/httpcontext"
	"github.com/SawitProRecruitment/UserService/lib/errors"
	"github.com/SawitProRecruitment/UserService/lib/jwt"
	"github.com/labstack/echo/v4"
)

func Auth(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) < 2 {
		return errors.NewForbiddenError("unauthorized")
	}

	bearer := strings.Trim(splitToken[1], " ")
	user, err := jwt.GetDataFromToken(bearer)
	if err != nil {
		return errors.NewForbiddenError("unauthorized")
	}
	c.Set(httpcontext.UserKey, user)
	return nil
}
