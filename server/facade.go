package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ResourceAPI/Core/config"
	"github.com/ResourceAPI/REST/nodes"
	"github.com/Vilsol/GoLib"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type RESTFacade struct {
	router http.Handler
}

// Initialize the facade.
func (facade RESTFacade) Initialize() error {
	router := mux.NewRouter()
	router.NotFoundHandler = GoLib.LoggerHandler(GoLib.NotFoundHandler())

	v1 := GoLib.RouteHandler(router, "/v1")
	nodes.RegisterResourceRoutes(v1)

	CORSHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	facade.router = GoLib.LoggerHandler(facade.router)
	facade.router = handlers.CompressHandler(facade.router)
	facade.router = handlers.ProxyHeaders(facade.router)
	facade.router = CORSHandler(facade.router)

	return nil
}

// Start the facade. Must be a blocking call.
func (facade RESTFacade) Start() error {
	fmt.Printf("REST server listening on port %d\n", config.Get().Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Get().Port), facade.router))
	return nil
}

// Graceful stopping of the facade with a 30s timeout.
func (facade RESTFacade) Stop() error {
	return nil // TODO
}