package handlers

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetTemplates(w http.ResponseWriter, r *http.Request) {
	userName, err := requests.NewGetTemplateRequest(r)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to parse request "))
		ape.Render(w, problems.BadRequest(err))
		return
	}

	client, err := ClientQ(r).GetByName(userName.User)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to get client"))
		ape.Render(w, problems.InternalError())
		return
	}

	if client == nil {
		Log(r).Error(errors.Wrap(err, "client is not found"))
		ape.RenderErr(w, problems.NotFound())
		return
	}

	tmps, err := TemplateQ(r).Select(client.ID)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to select templates "))
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
				BackgroundImg: string(tmp.ImgBytes),
				TemplateName:  tmp.Name,
			},
		})
	}
	return resources.TemplateListResponse{
		Data: reponseTmpList,
	}
}
