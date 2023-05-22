package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"net/http"
)

func UpdateToken(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUpdateTokenRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userInfo, err := helpers.ClientQ(r).GetByName(req.Data.Name)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if userInfo == nil {
		helpers.Log(r).WithError(err).Error("user is not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	userInfo.Token = nil
	userInfo.Code = req.Data.Code

	err = helpers.ClientQ(r).Update(userInfo)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	client := google.NewGoogleClient(helpers.Config(r))
	link, err := client.Connect(helpers.Config(r).Google().SecretPath, helpers.ClientQ(r), req.Data.Name)

	if len(link) != 0 {
		ape.Render(w, newLinkResponse(link))
		w.WriteHeader(201)
		return
	}

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect to google")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.RenderErr(w, problems.InternalError())

}
