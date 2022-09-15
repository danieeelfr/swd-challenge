package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/danieeelfr/swd-challenge/internal/config"
	apimodels "github.com/danieeelfr/swd-challenge/internal/httpserver/api/apimodels"
	handler "github.com/danieeelfr/swd-challenge/internal/httpserver/api/handler"

	"github.com/danieeelfr/swd-challenge/pkg/wait"
)

type router struct {
	e      *echo.Echo
	routes []*apimodels.Route
}

func newRouter(e *echo.Echo, cfg *config.Config, wg *wait.Wait) (*router, error) {
	t, err := handler.NewTaskHandler(cfg, e, wg)
	if err != nil {
		return nil, err
	}

	l, err := handler.NewLoginHandler(cfg, e, wg)
	if err != nil {
		return nil, err
	}

	allRoutes := make([]*apimodels.Route, 0)
	allRoutes = append(allRoutes, t.Routes...)
	allRoutes = append(allRoutes, l.Routes...)
	allRoutes = append(allRoutes, getHealthzRoute()...)

	return &router{e: e, routes: allRoutes}, nil
}

func getHealthzRoute() []*apimodels.Route {

	return []*apimodels.Route{
		{
			Method:   http.MethodGet,
			Endpoint: "/healthz",
			Handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			},
		},
	}
}

func (rtr *router) build() {

	for _, route := range rtr.routes {
		rtr.setRoute(route)
	}
}

func (rtr *router) setRoute(r *apimodels.Route) {
	switch r.Method {
	case http.MethodGet:
		rtr.e.GET(r.Endpoint, r.Handler)
	case http.MethodPost:
		rtr.e.POST(r.Endpoint, r.Handler)
	// case http.MethodPut:
	// 	rtr.e.PUT(r.Endpoint, r.Handler)
	// case http.MethodDelete:
	// 	rtr.e.DELETE(r.Endpoint, r.Handler)
	default:
		log.Errorf("http [%s] method not implemented", r.Method)
	}
}
