package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/StratoAPI/REST/nodes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type RESTFacade struct {
	echo   *echo.Echo
	router http.Handler
	server *http.Server
}

// Initialize the facade.
func (facade *RESTFacade) Initialize() error {

	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Pre(middleware.RemoveTrailingSlash())

	nodes.RegisterResourceRoutes(e.Group("/v1"))

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())

	facade.echo = e

	return nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%s", he.Message)
	}

	c.JSON(code, nodes.Response{
		Success: false,
		Error: &nodes.Error{
			Code:    -1,
			Message: message,
		},
	})
}

// Start the facade. Must be a blocking call.
func (facade *RESTFacade) Start() error {
	// TODO Per-Plugin Configs
	fmt.Printf("REST server listening on port %d\n", 5020)
	facade.echo.Start(fmt.Sprintf(":%d", 5020))
	return nil
}

// Graceful stopping of the facade with a 30s timeout.
func (facade *RESTFacade) Stop() error {
	return facade.echo.Shutdown(context.Background())
}
