package handlers

import (
	"encoding/json"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"net/http"
	"strings"
)

func CreateTemplate(w http.ResponseWriter, r *http.Request) {
	template, _, req, err := requests.NewGenerateTemplate(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to generate template")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	client, err := MasterQ(r).ClientQ().FilterByName(req.Data.Relationships.User.Data.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get client")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if client == nil {
		Log(r).WithError(err).Error("client is not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if req.Data.Attributes.IsCompleted {
		templateBytes, err := json.Marshal(template)
		if err != nil {
			Log(r).WithError(err).Error("failed to marshal")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		err = MasterQ(r).TemplateQ().Insert(&data.Template{
			Template:  templateBytes,
			ImgBytes:  []byte(strings.Replace(req.Data.Attributes.BackgroundImg, "data:image/png;base64,", "", 1)),
			Name:      req.Data.Attributes.TemplateName,
			ShortName: req.Data.Attributes.TemplateShortName,
			UserID:    client.ID,
		})
		if err != nil {
			Log(r).WithError(err).Error("failed to insert template")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
