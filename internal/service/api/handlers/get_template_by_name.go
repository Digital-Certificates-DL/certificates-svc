package handlers

import (
	"encoding/base64"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetTemplateByName(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetTemplateByNameRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request ")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	client, err := MasterQ(r).ClientQ().FilterByName(request.User).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get client")
		ape.Render(w, problems.InternalError())
		return
	}

	if client == nil {
		Log(r).WithError(err).Error("client is not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	tmp, err := MasterQ(r).TemplateQ().FilterByName(request.TemplateName).FilterByUser(client.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to select templates ")
		ape.Render(w, problems.InternalError())
		return
	}

	if tmp != nil {
		Log(r).WithError(err).Error("template is not found")
		ape.Render(w, problems.NotFound())
		return
	}
	Log(r).Debug("template: ", tmp.Template)

	ape.Render(w, newTemplateResp(tmp))
}

func newTemplateResp(tmp *data.Template) resources.TemplateResponse {
	return resources.TemplateResponse{
		Data: resources.Template{
			Attributes: resources.TemplateAttributes{
				BackgroundImg: base64.StdEncoding.EncodeToString(tmp.ImgBytes),
				TemplateName:  tmp.Name,
				Template:      tmp.Template,
			},
		},
	}
}
