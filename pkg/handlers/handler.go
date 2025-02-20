package handlers

import (
	"ducky/http/pkg/request"
)

type Handler interface {
	Handle(*request.Request) error
}
