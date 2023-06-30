package main

import (
	"EnronEmailApi/controllers"
	"EnronEmailApi/zinc"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func main() {
	//services.Init()
	zinc.EstablishConnection()

	r := chi.NewRouter()
	r.Post("/indexer", controllers.IndexerEnron)
	r.Post("/search", controllers.SearchEmails)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	// Agregar el middleware CORS a todas las rutas
	handler := corsOptions.Handler(r)

	http.ListenAndServe(":3000", handler)

	//http.ListenAndServe(":3000", r)

}
