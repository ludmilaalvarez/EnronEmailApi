package controllers

import (
	"EnronEmailApi/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func IndexerEnron(w http.ResponseWriter, r *http.Request) {

	go services.IndexStart()
	w.Write([]byte("Cargando emails..."))
}

func SearchEmails(w http.ResponseWriter, r *http.Request) {

	var text string

	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte un item válido")
	}

	json.Unmarshal(rqBody, &text)

	respuesta := services.SearchEmails(&text)

	json.NewEncoder(w).Encode(respuesta)

}
