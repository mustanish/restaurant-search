package routes

import (
	"search/server/handlers"

	"github.com/go-chi/chi"
)

func restaurant() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", handlers.SearchRestaurants)
	return router
}
