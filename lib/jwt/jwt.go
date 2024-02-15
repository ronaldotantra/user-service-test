package jwt

import (
	"fmt"
	"time"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/config/env"

	jwt "github.com/dgrijalva/jwt-go"
)

var cfg *config.Config

func init() {
	envConfig := env.New()
	config.Init(envConfig)
	cfg = config.Load()
}

// MyClaims .
type MyClaims struct {
	jwt.StandardClaims
	User User `json:"user"`
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// GenerateToken .
func GenerateToken(user User) (*string, error) {
	privKey := cfg.JwtPrivateKey()

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cfg.ApplicationName(),
			Subject:   "Auth",
			Id:        fmt.Sprint(time.Now().UnixNano()),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(1 * 24 * time.Hour).Unix(),
		},
		User: user,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(privKey))
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

// GetDataFromToken .
func GetDataFromToken(param string) (*User, error) {
	token, err := jwt.ParseWithClaims(param, &MyClaims{}, func(x *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtPrivateKey()), nil
	})

	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return nil, err
	}
	return &claims.User, nil
}
