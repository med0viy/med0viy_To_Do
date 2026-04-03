package core_http_server

import (
	"net/http"

	core_http_maddleware "github.com/med0viy/practika/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	middleware []core_http_maddleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_maddleware.ChainMiddleWare(
		r.Handler,
		r.middleware...,
	)
}
