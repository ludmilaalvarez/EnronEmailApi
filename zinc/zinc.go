package zinc

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Zinc_url           string
	Bas64encoded_creds string
)

func EstablishConnection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	var (
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
	)

	auth := user + ":" + password
	Bas64encoded_creds = base64.StdEncoding.EncodeToString([]byte(auth))

	zinc_host := "http://localhost:4080"
	Zinc_url = zinc_host + "/api/_bulkv2"

}
