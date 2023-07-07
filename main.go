package main

import (
	"EnronEmailApi/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func main() {

	r := chi.NewRouter()
	r.Post("/indexer", controllers.IndexerEnron)
	r.Post("/search/{text}", controllers.SearchEmails)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	// Agregar el middleware CORS a todas las rutas
	handler := corsOptions.Handler(r)

	http.ListenAndServe(":3000", handler)
}
