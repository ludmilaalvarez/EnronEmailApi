package controllers

import (
	"EnronEmailApi/services"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func IndexerEnron(w http.ResponseWriter, r *http.Request) {

	go services.IndexStart()

	w.Write([]byte("Cargando emails..."))
}

func SearchEmails(w http.ResponseWriter, r *http.Request) {
	text := chi.URLParam(r, "text")

	fmt.Println(text)

	respuesta := services.SearchEmails(&text)

	json.NewEncoder(w).Encode(respuesta.Hits.Hits)
}
