package handlers

import (
	"fmt"
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
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}
	if len(link) != 0 {
		helpers.Log(r).WithError(err).Error("failed to authorize")
		w.WriteHeader(http.StatusForbidden)

		ape.Render(w, newLinkResponse(link))
		helpers.Log(r).Info(w.Header())
		return
	}

	users, errs, isTokenExpired := client.ParseFromWeb(req.Data.Url, "A1:K", helpers.Config(r).Log())
	if errs != nil {
		helpers.Log(r).Error("failed to parse table: Errors:", errs)
		ape.Render(w, problems.BadRequest(err))
		return
	}
	if isTokenExpired {
		helpers.Log(r).Error("failed to connect: token expired")
		w.WriteHeader(http.StatusForbidden)
		ape.Render(w, newTokenExpired())
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
func newTokenExpired() resources.ExpiredTokenErrorResponse {
	return resources.ExpiredTokenErrorResponse{
		Data: resources.ExpiredTokenError{
			Attributes: resources.ExpiredTokenErrorAttributes{
				Error: true,
			},
		},
	}
}
