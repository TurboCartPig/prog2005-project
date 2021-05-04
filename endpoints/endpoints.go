package endpoints

import (
	"developer-bot/discord"
	"developer-bot/endpoints/firestore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"developer-bot/endpoints/types"
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
		port = "80"
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

	// Create a new firestore client in the firestore package. 
	firestore.NewFirestoreClient()

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
	var newWebhook types.WebhookData
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, "Working")

	firestore.SaveWebhookToFirestore(&newWebhook)
	discord.SendMessage("833465870872608788", "Hei")

	w.WriteHeader(http.StatusOK)
}
