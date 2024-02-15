package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/config/env"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/lib/errors"
	"github.com/SawitProRecruitment/UserService/lib/validator"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var cfg *config.Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envConfig := env.New()
	config.Init(envConfig)
	cfg = config.Load()
}

// @contact.name Ronaldo Tantra
// @contact.email ronaldotantra@gmail.com
func main() {
	e := echo.New()
	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler
	e.Validator = validator.NewValidator()

	var server generated.ServerInterface = newServer(cfg)
	generated.RegisterHandlers(e, server)

	e.Logger.Fatal(e.Start(fmt.Sprint(":", cfg.AppPort())))
}

func newServer(config *config.Config) *handler.Server {
	dbDsn := config.DatabaseUrl()
	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		panic(err)
	}
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Db: db,
	})
	var service service.ServiceInterface = service.NewService(service.NewServiceOption{
		UserRepository: repo,
	})
	opts := handler.NewServerOptions{
		Service: service,
	}
	return handler.NewServer(opts)
}
