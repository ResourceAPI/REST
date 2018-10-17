package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/StratoAPI/REST/nodes"
	"github.com/Vilsol/GoLib"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type RESTFacade struct {
	router http.Handler
}

// Initialize the facade.
func (facade *RESTFacade) Initialize() error {
	router := mux.NewRouter()
	router.NotFoundHandler = GoLib.LoggerHandler(GoLib.NotFoundHandler())

	v1 := GoLib.RouteHandler(router, "/v1")
	nodes.RegisterResourceRoutes(v1)

	CORSHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	finalRouter := GoLib.LoggerHandler(router)
	finalRouter = handlers.CompressHandler(finalRouter)
	finalRouter = handlers.ProxyHeaders(finalRouter)
	finalRouter = CORSHandler(finalRouter)

	facade.router = finalRouter

	return nil
}

// Start the facade. Must be a blocking call.
func (facade *RESTFacade) Start() error {
	// TODO Per-Plugin Configs
	fmt.Printf("REST server listening on port %d\n", 5020)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 5020), facade.router))
	return nil
}

// Graceful stopping of the facade with a 30s timeout.
func (facade *RESTFacade) Stop() error {
	return nil // TODO
}
