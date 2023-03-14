package handlers

import (
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"helper/internal/data"
	"helper/internal/handlers"
	"helper/internal/service/google"
	"helper/internal/service/helpers"
	"helper/internal/service/qr"
	"helper/internal/service/requests"
	"helper/internal/service/signature"
	"log"
	"net/http"
	"os"
)

func GenerateTable(w http.ResponseWriter, r *http.Request) {
	var paths []handlers.FilesBytes
	request, err := requests.NewGenerateTable(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	os.MkdirAll(helpers.Config(r).QRCode().QRPath, os.ModePerm) //todo maybe remove it

	client := google.NewGoogleClient(helpers.Config(r))
	err = client.Connect(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	if err != nil {
		log.Println(err)
		return
	}

	//err = client.ConnectToDrive(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	users, errs := client.ParseFromWeb(request.Id, "A1:K", helpers.Config(r).Log())
	if errs != nil {
		helpers.Log(r).Error("failed to parse table: Errors:", errs)
		ape.Render(w, problems.BadRequest(err))
		return
	}

	//todo make signature in front and use gorutines
	sign, err := signature.NewSignature(helpers.Config(r).Key().Private)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create signature")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	var usersResult []*data.User

	for id, user := range users {
		user.ID = id
		if user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" {
			helpers.Log(r).Debug("has already")
			//todo maybe add render
			continue
		}

		qr := qr.NewQR(user, helpers.Config(r), sign)
		hash := sign.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)) //todo signing in frontend and return it in back

		if hash != "" {
			helpers.Log(r).Info(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		var path string
		path, user.DigitalCertificate, user.Signature, err = qr.GenerateQR()
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to generate qr")
			ape.Render(w, problems.InternalError())
			return
		}
		_ = path
		paths = append(paths, handlers.FilesBytes{File: nil, ID: id}) //todo
		//paths = append(paths, handlers.FilesBytes{File: path, ID: id}) //todo
		usersResult = append(usersResult, user)
	}

	//todo  add new event that handle error with connect to drive
	if sendToDrive := helpers.Config(r).Google().Enable; sendToDrive {
		users, err = handlers.Drive(client, helpers.Config(r), helpers.Log(r), paths, users)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to send date to drive")
			ape.Render(w, problems.InternalError())
			return
		}
	}

	helpers.Log(r).Info("creating table")
	errs = client.SetRes(usersResult, request.Id)
	if errs != nil {
		helpers.Log(r).Error("failed to send date to drive: Errors: ", errs)
		ape.Render(w, problems.InternalError())
		return
	}
	return
}
