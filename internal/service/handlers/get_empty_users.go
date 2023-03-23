package handlers

import (
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"log"
	"net/http"
)

func GetUsersEmpty(w http.ResponseWriter, r *http.Request) {
	tableID, err := requests.NewGetUsers(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	client := google.NewGoogleClient(helpers.Config(r))
	err = client.Connect(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	if err != nil {
		log.Println(err)
		return
	}

	users, errs := client.ParseFromWeb(tableID, "A1:K", helpers.Config(r).Log())
	if errs != nil {
		helpers.Log(r).Error("failed to parse table: Errors:", errs)
		ape.Render(w, problems.BadRequest(err))
		return
	}
	emptyUsers := make([]*data.User, 0)
	for id, user := range users {
		user.ID = id
		if user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" || user.SerialNumber != "" {
			helpers.Log(r).Debug("has already")
			continue
		}
		user.Msg = fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		emptyUsers = append(emptyUsers, user)
	}

	ape.Render(w, newActionResponse(emptyUsers)) //todo make better
}

func newActionResponse(users []*data.User) resources.UserListResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.USER,
			},
			Attributes: resources.UserAttributes{
				Certificate:  user.Certificate,
				ID:           int64(user.ID),
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
