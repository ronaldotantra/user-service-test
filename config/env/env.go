package env

import (
	"fmt"
	"os"
	"strconv"
)

// Env .
type Env struct{}

// Environment .
func (e *Env) Environment() string {
	return getStringOrDefault(EnvironmentEnv, "local")
}

// ApplicationName .
func (e *Env) ApplicationName() string {
	return getStringOrDefault(ApplicationName, "SawitPro User Service")
}

// AppPort .
func (e *Env) AppPort() int {
	return getIntOrDefault(AppPort, 1323)
}

// JwtPrivateKey .
func (e *Env) JwtPrivateKey() string {
	return getStringOrDefault(JwtPrivateKey, "randomstring-asdasd")
}

// DatabaseUrl .
func (e *Env) DatabaseUrl() string {
	return getStringOrDefault(DatabaseUrl, "host=127.0.0.1 port=5432 user=postgres "+
		"password=postgres dbname=postgres sslmode=disable")
}

// New .
func New() *Env {
	return &Env{}
}

func getStringOrDefault(key, def string) string {
	return getEnvOrDefault(key, def)
}

func getIntOrDefault(key string, def int) int {
	results := getEnvOrDefault(key, fmt.Sprint(def))
	i, err := strconv.ParseInt(results, 10, 64)
	if err != nil {
		return def
	}
	return int(i)
}

func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
