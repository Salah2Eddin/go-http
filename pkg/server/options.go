package server

import (
	"github.com/Salah2Eddin/go-http/pkg/readers"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/router"
	"github.com/Salah2Eddin/go-http/pkg/serializers"
)

type Option func(*Server)

// WithAddress sets the server address (IP and Port)
func WithAddress(addr Address) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithIP sets the server IP address
func WithIP(ip string) Option {
	return func(s *Server) {
		s.addr.IP = ip
	}
}

// WithPort sets the server port
func WithPort(port string) Option {
	return func(s *Server) {
		s.addr.Port = port
	}
}

// WithRouter sets a custom router implementation
func WithRouter(r router.IRouter) Option {
	return func(s *Server) {
		s.router = r
	}
}

// WithReader sets a custom request reader implementation
func WithReader(r readers.IRequestReader) Option {
	return func(s *Server) {
		s.reader = r
	}
}

// WithSerializer sets a custom response serializer implementation
func WithSerializer(ser serializers.ISerializer[*response.Response]) Option {
	return func(s *Server) {
		s.serializer = ser
	}
}

// WithErrorResponder sets a custom error responder implementation
func WithErrorResponder(er IErrorResponder) Option {
	return func(s *Server) {
		s.errorResponder = er
	}
}
