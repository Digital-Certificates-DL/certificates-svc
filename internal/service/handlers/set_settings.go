package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"net/http"
)

func SetSettings(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewSetSettings(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	settings, err := helpers.ClientQ(r).GetByName(req.Data.Name)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get settings")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	if settings == nil {
		user := data.Client{
			Name: req.Data.Name,
			Code: req.Data.Code,
		}
		_, err := helpers.ClientQ(r).Insert(&user)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get settings")
			ape.Render(w, problems.InternalError())
			return
		}
		w.WriteHeader(204)
		return
	}
	settings.Code = req.Data.Code
	err = helpers.ClientQ(r).Update(settings)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update settings")
		ape.Render(w, problems.InternalError())
		return
	}
	if req.Data.Code != "" {
		settings.Token = nil
		err = helpers.ClientQ(r).Update(settings)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to update settings")
			ape.Render(w, problems.InternalError())
			return
		}
		client := google.NewGoogleClient(helpers.Config(r))
		_, err = client.Connect(helpers.Config(r).Google().SecretPath, helpers.ClientQ(r), req.Data.Name)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to connect")
			ape.Render(w, problems.InternalError())
			return
		}
	}
	w.WriteHeader(204)
	return
}
