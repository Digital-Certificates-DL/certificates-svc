package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"log"
	"net/http"
	"strings"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetUsers(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	client := google.NewGoogleClient(helpers.Config(r))
	link, err := client.Connect(helpers.Config(r).Google().SecretPath, helpers.ClientQ(r), req.Data.Name)
	if err != nil { //todo
		log.Println(err)
		return
	}

	if len(link) != 0 {
		helpers.Log(r).WithError(err).Error("failed to authorize")
		ape.Render(w, newLinkResponse(link))
		w.WriteHeader(403)

		return
	}

	users, errs := client.ParseFromWeb(req.Data.Url, "A1:K", helpers.Config(r).Log())
	if errs != nil {
		helpers.Log(r).Error("failed to parse table: Errors:", errs)
		ape.Render(w, problems.BadRequest(err))
		return
	}
	readyUsers := make([]*helpers.User, 0)
	for id, user := range users {
		user.ID = id
		if user.Certificate == "" {
			continue
		}
		if !strings.Contains(user.Certificate, "http") {
			continue
		}
		//file, err := client.Download(user.Certificate)
		//if err != nil {
		//	helpers.Log(r).Error("failed to ", err)
		//	ape.Render(w, problems.BadRequest(err))
		//	return
		//}
		//img, err := pdf.Convert("png", file)
		//if err != nil {
		//	helpers.Log(r).Error("failed to convert", err)
		//	ape.Render(w, problems.BadRequest(err))
		//	return
		//}
		//user.ImageCertificate = img
		readyUsers = append(readyUsers, user)
	}
	ape.Render(w, newUserResponse(readyUsers))
}
