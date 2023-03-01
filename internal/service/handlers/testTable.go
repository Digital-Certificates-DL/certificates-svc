package handlers

import (
	"helper/internal/service/google"
	"helper/internal/service/helpers"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	client := google.NewGoogleClient(helpers.Config(r))
	err := client.ConnectTOSheet(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	if err != nil {
		log.Println(err)
		return
	}

	readRange := "A1:C"

	client.GetTable(readRange)

}
