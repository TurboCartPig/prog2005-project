package endpoints

import (
	"developer-bot/discord"
	"encoding/json"
	"fmt"
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

//This should not be a global var, change it soon
var webhooks []WebhookData

func developer(w http.ResponseWriter, r *http.Request) {
	var newWebhook WebhookData
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	webhooks = append(webhooks, newWebhook)
	fmt.Fprint(w, "")

	discord.SendMessage("833465870872608788", "Hei")

	w.WriteHeader(http.StatusOK)
}
