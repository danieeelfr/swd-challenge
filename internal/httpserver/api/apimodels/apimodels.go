package api

import "github.com/labstack/echo/v4"

// Route contains info to handler the requests
type Route struct {
	Method   string
	Endpoint string
	Handler  func(c echo.Context) error
}
