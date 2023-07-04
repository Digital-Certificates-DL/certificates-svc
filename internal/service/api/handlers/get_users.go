package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"net/http"
	"strings"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetUsers(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	client := google.NewGoogleClient(Config(r))
	link, err := client.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), req.Data.Name)

	if len(link) != 0 {
		Log(r).WithError(err).Error("failed to authorize")

		ape.RenderErr(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "Invalid token",
			Status: "403",
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}
	if err != nil {
		Log(r).WithError(err).Error("failed to authorize")
		if strings.Contains(err.Error(), "unable to get client") {
			ape.RenderErr(w, problems.NotFound())
			return
		}
		if strings.Contains(err.Error(), "Token has been expired or revoked") {
			ape.RenderErr(w, problems.Unauthorized())
			return
		}
		Log(r).WithError(err).Error("failed to authorize")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	users, errs := client.ParseFromWeb(req.Data.Url, "A1:K", Config(r).Log())
	if errs != nil {
		fmt.Println(errs[0].Error())

		if strings.Contains(errs[0].Error(), "400") {
			Log(r).Error("token expired")
			ape.RenderErr(w, problems.Unauthorized())
			return
		}

		Log(r).Error("failed to parse table: Errors:", errs)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	readyUsers := make([]*helpers.Certificate, 0)
	for id, user := range users {
		user.ID = id
		user.ShortCourseName = Config(r).TemplatesConfig()[user.CourseTitle]
		user.Msg = fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		readyUsers = append(readyUsers, user)
	}
	ape.Render(w, newUserResponse(readyUsers))
}
