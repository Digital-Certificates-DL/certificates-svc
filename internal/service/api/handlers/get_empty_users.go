package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetUsersEmpty(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetUsers(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
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
		Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	users, errs := client.ParseFromWeb(req.Data.Url, "A1:K", Config(r).Log())
	if errs != nil {
		Log(r).Error("failed to parse table: Errors:", errs)
		ape.Render(w, problems.BadRequest(err))
		return
	}
	emptyUsers := make([]*helpers.User, 0)
	for id, user := range users {
		user.ID = id
		user.ShortCourseName = Config(r).TemplatesConfig()[user.CourseTitle]
		if user.Certificate != "" {
			Log(r).Debug("has already")
			continue
		}
		user.Msg = fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		emptyUsers = append(emptyUsers, user)
	}

	ape.Render(w, newUserResponse(emptyUsers))
}

func newUserResponse(users []*helpers.User) resources.UserListResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.USER,
			},
			Attributes: resources.UserBlob{
				Certificate:    user.Certificate,
				Id:             int64(user.ID),
				Points:         user.Points,
				Participant:    user.Participant,
				Msg:            user.Msg,
				SerialNumber:   user.SerialNumber,
				Note:           user.Note,
				DataHash:       user.DataHash,
				Signature:      user.Signature,
				TxHash:         user.TxHash,
				Date:           user.Date,
				CertificateImg: user.ImageCertificate,
				CourseTitle:    user.CourseTitle,
				ShortCourse:    user.ShortCourseName,
			},
		}
		usersData = append(usersData, resp)
	}

	return resources.UserListResponse{
		Data: usersData,
	}

}
