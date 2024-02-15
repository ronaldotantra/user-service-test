package config

// IConfig we want to make sure the driver of env should implement IConfig to avoid
// missing implementation
type IConfig interface {
	Environment() string
	ApplicationName() string
	AppPort() int
	DatabaseUrl() string
	JwtPrivateKey() string
}
