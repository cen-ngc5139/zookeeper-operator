package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("GET", "/status", getZkStatus, nil)
	register("GET", "/client", getZkClient, nil)
	register("GET", "/runok", getZkRunok, nil)
	register("GET", "/health", Health, nil)
	register("GET", "/get", getMember, nil)
	register("POST", "/add", addMember, nil)
	register("POST", "/del", delMember, nil)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		r := router.Methods(route.Method).
			Path(route.Pattern)
		if route.Middleware != nil {
			r.Handler(route.Middleware(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
