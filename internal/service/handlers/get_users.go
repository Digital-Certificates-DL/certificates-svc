package handlers

import (
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"helper/internal/data"
	"helper/internal/service/google"
	"helper/internal/service/helpers"
	"helper/internal/service/requests"
	"helper/resources"
	"log"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
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

	//header.Add()
	//header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	//header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	ape.Render(w, newActionResponse(users)) //todo make better
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
