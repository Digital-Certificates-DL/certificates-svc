package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func GetTemplates(w http.ResponseWriter, r *http.Request) {
	userName, err := requests.NewGetTemplateRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse request ")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	client, err := MasterQ(r).ClientQ().GetByName(userName.User)
	if err != nil {
		Log(r).WithError(err).Debug("failed to get client")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if client == nil {
		Log(r).WithError(err).Debug("client is not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	tmps, err := MasterQ(r).TemplateQ().Select(client.ID)
	if err != nil {
		Log(r).WithError(err).Debug("failed to select templates ")
		ape.RenderErr(w, problems.InternalError())
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
