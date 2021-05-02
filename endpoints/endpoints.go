package endpoints

import (
	"github.com/go-chi/chi"
	"net/http"
)

const Version = "V1"


// Will set up all of the developer-bot supplied endpoints.
func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/endpoints", developer)

	return router
}


func developer(w http.ResponseWriter, r *http.Request) {
}