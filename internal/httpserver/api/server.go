package api

import (
	"context"

	"github.com/danieeelfr/swd-challenge/internal/config"
	requestvalidator "github.com/danieeelfr/swd-challenge/internal/httpserver/utils/requestvalidator"
	"github.com/danieeelfr/swd-challenge/pkg/wait"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("package", "api.httpserver")
)

// Server holds the http server implementation
type Server struct {
	conf   *config.HTTPServerConfig
	e      *echo.Echo
	router *router
	wait   *wait.Wait
}

// New http server or error
func New(cfg *config.Config, wg *wait.Wait) (*Server, error) {
	srv := new(Server)
	srv.conf = cfg.HTTPServerConfig
	srv.wait = wg
	srv.e = echo.New()
	srv.e.HideBanner = false
	srv.e.Validator = &requestvalidator.CustomValidator{Validator: validator.New()}

	var err error
	srv.router, err = newRouter(srv.e, cfg, wg)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// Start the http server
func (srv *Server) Start() error {
	log.Infof("starting http server on host:[%s]", srv.conf.HTTPServerHost)
	srv.wait.Add()

	srv.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `method=${method}, uri=${uri}, status=${status}` + "\n",
	}))

	srv.e.Use(middleware.Recover())

	srv.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // TODO: Security improve needed
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	srv.router.build()

	go func() {
		if err := srv.e.Start(":" + srv.conf.HTTPServerHost); err != nil {
			if !srv.wait.IsBlock() {
				log.WithError(err).Fatalf("error starting http server. error: [%#v]", err)
			}

		}
	}()

	return nil
}

// Shutdown the http server
func (srv *Server) Shutdown() {
	defer srv.wait.Done()
	if err := srv.e.Shutdown(context.Background()); err != nil {
		log.WithError(err).Error("failed to shutdown http server.")
	}
}
