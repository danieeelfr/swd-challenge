package api

import (
	"errors"
	"net/http"

	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	apimodels "github.com/danieeelfr/swd-challenge/internal/httpserver/api/apimodels"
	responses "github.com/danieeelfr/swd-challenge/internal/httpserver/api/apimodels/response"
	"github.com/danieeelfr/swd-challenge/internal/models"
	loginservice "github.com/danieeelfr/swd-challenge/internal/service/login"
	"github.com/danieeelfr/swd-challenge/pkg/wait"
)

func init() {
	middleware.ErrJWTMissing.Code = 401
	middleware.ErrJWTMissing.Message = "Unauthorized - ErrJWTMissing"
	middleware.ErrJWTInvalid.Message = "Unauthorized - ErrJWTInvalid"
}

// LoginHandler holds the login handler implementation
type LoginHandler struct {
	loginService *loginservice.Login
	Routes       []*apimodels.Route
	wait         *wait.Wait
}

// NewLoginHandler returns a new handler or error
func NewLoginHandler(cfg *config.Config, e *echo.Echo, wg *wait.Wait) (*LoginHandler, error) {
	h := new(LoginHandler)
	h.wait = wg

	ls, err := loginservice.NewLoginService(cfg)
	if err != nil {
		return nil, err
	}
	h.loginService = ls

	h.Routes = append(h.Routes, h.getPublicRoutes()...)

	g := e.Group("")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	config := middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	}

	g.Use(middleware.JWTWithConfig(config))

	return h, nil

}

func (h *LoginHandler) getPublicRoutes() []*apimodels.Route {

	return []*apimodels.Route{
		{
			Method:   http.MethodPost,
			Endpoint: "login/authorize",
			Handler:  h.Authorize,
		},
	}
}

func (h *LoginHandler) Authorize(ctx echo.Context) error {

	login := new(models.Login)

	if err := ctx.Bind(login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user, err := h.loginService.GetUser(ctx.Request().Context(), login.User, login.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("invalid user and/or password"))
	}

	token, err := h.loginService.GetToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(http.StatusOK, responses.LoginResponse{
		Message: "authorized with success!",
		Token:   token,
	})
}
