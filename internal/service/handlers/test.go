package handlers

import (
	"helper/internal/service/google"
	"helper/internal/service/helpers"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {

	client := google.NewGoogleClient(helpers.Config(r))
	err := client.Connect(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	if err != nil {
		log.Println(err)
		return
	}

	id := "1CYqLid0t90bgGm1HPx5j8q-h_RNVVLPVkot9iJZguuo"
	client.UpdateTable("sheet1!A5", []string{""}, id)

}
