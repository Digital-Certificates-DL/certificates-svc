package handlers

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetTemplateByName(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetTemplateByNameRequest(r)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to parse request "))
		ape.Render(w, problems.BadRequest(err))
		return
	}

	client, err := helpers.ClientQ(r).GetByName(request.User)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to get client"))
		ape.Render(w, problems.InternalError())
		return
	}

	if client == nil {
		helpers.Log(r).Error(errors.Wrap(err, "client is not found"))
		ape.RenderErr(w, problems.NotFound())
		return
	}

	tmp, err := helpers.TemplateQ(r).GetByName(request.TemplateName, client.ID)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to select templates "))
		ape.Render(w, problems.InternalError())
		return
	}

	if tmp != nil {
		helpers.Log(r).Error(errors.Wrap(err, "template is not found"))
		ape.Render(w, problems.NotFound())
		return
	}

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
