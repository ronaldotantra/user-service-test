package config

import (
	"fmt"
)

// Config .
type Config struct {
	c IConfig
}

var defaultConfig = new(Config)

// ErrNoInitializationConfig .
var ErrNoInitializationConfig = fmt.Errorf("config: %s", "No driver is loaded, use blank import to load driver to config.")

// Environment .
func (c *Config) Environment() string {
	return c.c.Environment()
}

// ApplicationName .
func (c *Config) ApplicationName() string {
	return c.c.ApplicationName()
}

// AppPort .
func (c *Config) AppPort() int {
	return c.c.AppPort()
}

// JwtPrivateKey .
func (c *Config) JwtPrivateKey() string {
	return c.c.JwtPrivateKey()
}

// DatabaseUrl .
func (c *Config) DatabaseUrl() string {
	return c.c.DatabaseUrl()
}

// Init .
func Init(c IConfig) {
	defaultConfig.c = c
}

// Load .
func Load() *Config {
	if defaultConfig.c == nil {
		panic(ErrNoInitializationConfig)
	}
	return defaultConfig
}
