package httpcontext

import (
	"github.com/SawitProRecruitment/UserService/lib/jwt"
	"github.com/labstack/echo/v4"
)

const UserKey = "user"

// GetUserJWT get user response from context
func GetUserJWT(c echo.Context) (*jwt.User, bool) {
	resUser, ok := c.Get(UserKey).(*jwt.User)

	return resUser, ok
}
