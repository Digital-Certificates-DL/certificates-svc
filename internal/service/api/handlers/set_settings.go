package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"net/http"
)

func SetSettings(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewSetSettings(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	settings, err := MasterQ(r).ClientQ().GetByName(req.Data.Attributes.Name)
	if err != nil {
		Log(r).WithError(err).Error("failed to get settings")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	if settings == nil {
		user := data.Client{
			Name: req.Data.Attributes.Name,
			Code: req.Data.Attributes.Code,
		}
		_, err := MasterQ(r).ClientQ().Insert(&user)
		if err != nil {
			Log(r).WithError(err).Error("failed to get settings")
			ape.Render(w, problems.InternalError())
			return
		}
		w.WriteHeader(204)
		return
	}
	settings.Code = req.Data.Attributes.Code
	err = MasterQ(r).ClientQ().Update(settings)
	if err != nil {
		Log(r).WithError(err).Error("failed to update settings")
		ape.Render(w, problems.InternalError())
		return
	}
	if req.Data.Attributes.Code != "" {
		client := google.NewGoogleClient(Config(r))
		_, err = client.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), req.Data.Attributes.Name)
		if err != nil {
			Log(r).WithError(err).Error("failed to connect")
			ape.Render(w, problems.InternalError())
			return
		}
	}
	w.WriteHeader(204)
	return
}
