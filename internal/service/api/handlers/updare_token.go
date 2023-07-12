package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"net/http"
)

func UpdateToken(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUpdateTokenRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userInfo, err := MasterQ(r).ClientQ().GetByName(req.Data.Attributes.Name)
	if err != nil {
		Log(r).WithError(err).Debug("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if userInfo == nil {
		Log(r).WithError(err).Debug("user is not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	userInfo.Token = nil
	userInfo.Code = req.Data.Attributes.Code

	err = MasterQ(r).ClientQ().Update(userInfo)
	if err != nil {
		Log(r).WithError(err).Debug("failed to update user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	client := google.NewGoogleClient(Config(r))
	link, err := client.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), req.Data.Attributes.Name)

	if len(link) != 0 {
		ape.Render(w, newLinkResponse(link))
		w.WriteHeader(204)
		return
	}

	if err != nil {
		Log(r).WithError(err).Debug("failed to connect to google")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.RenderErr(w, problems.InternalError())

}
