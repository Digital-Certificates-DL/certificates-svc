package handlers

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetTemplates(w http.ResponseWriter, r *http.Request) {
	userName, err := requests.NewGetTemplateRequest(r)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to parse request "))
		ape.Render(w, problems.BadRequest(err))
		return
	}

	client, err := helpers.ClientQ(r).GetByName(userName.User)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to get client"))
		ape.Render(w, problems.InternalError())
		return
	}

	if client == nil {
		helpers.Log(r).Error(errors.Wrap(err, "client is not found"))
		ape.Render(w, problems.NotFound())
		return
	}

	tmps, err := helpers.TemplateQ(r).Select(client.ID)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to select templates "))
		ape.Render(w, problems.InternalError())
		return
	}
	ape.Render(w, newTemlateListResp(tmps))

}

func newTemlateListResp(tmps []data.Template) resources.TemplateListResponse {
	var reponseTmpList []resources.Template
	for _, tmp := range tmps {
		reponseTmpList = append(reponseTmpList, resources.Template{
			Attributes: resources.TemplateAttributes{
				BackgroundImg: tmp.ImgBytes,
				TemplateName:  tmp.Name,
			},
		})
	}
	return resources.TemplateListResponse{
		Data: reponseTmpList,
	}
}
