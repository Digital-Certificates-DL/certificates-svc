package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetUsersEmpty(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetUsers(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
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
			Code:   "125",
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	users, errs := client.ParseFromWeb(req.Data.Url, "A1:K", helpers.Config(r).Log())
	if errs != nil {
		helpers.Log(r).Error("failed to parse table: Errors:", errs)
		ape.Render(w, problems.BadRequest(err))
		return
	}
	emptyUsers := make([]*helpers.User, 0)
	for id, user := range users {
		user.ID = id
		if user.Certificate != "" {
			helpers.Log(r).Debug("has already")
			continue
		}
		user.Msg = fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		emptyUsers = append(emptyUsers, user)
	}

	ape.Render(w, newUserResponse(emptyUsers)) //todo make better
}

func newUserResponse(users []*helpers.User) resources.UserListResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.USER,
			},
			Attributes: resources.UserAttributes{
				Certificate:  user.Certificate,
				Id:           int64(user.ID),
				Points:       user.Points,
				Participant:  user.Participant,
				Msg:          user.Msg,
				SerialNumber: user.SerialNumber,
				Note:         user.Note,
				DataHash:     user.DataHash,
				Signature:    user.Signature,
				TxHash:       user.TxHash,
				Date:         user.Date,
				CourseTitle:  user.CourseTitle,
			},
		}
		usersData = append(usersData, resp)
	}

	return resources.UserListResponse{
		Data: usersData,
	}

}
