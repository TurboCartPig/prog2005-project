package endpoints

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Version of the API
const Version = "v1"

// getPort from environment variable
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return port
}

// Serve the REST API over HTTP
func Serve() {
	port := getPort()
	router := routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(handler2 http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging error: %s\n", err.Error())
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

// routes sets up routes
func routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Post("/endpoints", developer)

	return router
}

func developer(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

