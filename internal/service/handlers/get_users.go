package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"net/http"
	"strings"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetUsers(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	client := google.NewGoogleClient(helpers.Config(r))
	link, err := client.Connect(helpers.Config(r).Google().SecretPath, helpers.ClientQ(r), req.Data.Name)

	if len(link) != 0 {
		helpers.Log(r).WithError(err).Error("failed to authorize")
		w.Header().Set("auth_link", link)

		ape.RenderErr(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "Invalid token",
			Status: "403",
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to authorize")
		if strings.Contains(err.Error(), "unable to get client") {
			ape.RenderErr(w, problems.NotFound())
			return
		}
		helpers.Log(r).WithError(err).Error("failed to authorize")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	users, errs := client.ParseFromWeb(req.Data.Url, "A1:K", helpers.Config(r).Log())
	if errs != nil {
		fmt.Println(errs[0].Error())

		if strings.Contains(errs[0].Error(), "400") {
			helpers.Log(r).Error("token expired")
			ape.RenderErr(w, problems.Unauthorized())
			return
		}
		helpers.Log(r).Error("failed to parse table: Errors:", errs)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	readyUsers := make([]*helpers.User, 0)
	for id, user := range users {
		user.ID = id

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
		user.Msg = fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		readyUsers = append(readyUsers, user)
	}
	ape.Render(w, newUserResponse(readyUsers))
}
