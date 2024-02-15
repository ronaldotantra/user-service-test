package handler

import (
	"github.com/SawitProRecruitment/UserService/service"
)

type Server struct {
	Service service.ServiceInterface
}

type NewServerOptions struct {
	Service service.ServiceInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Service: opts.Service,
	}
}
