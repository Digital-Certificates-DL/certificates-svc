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
		Log(r).WithError(err).Error("failed to parse request ")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	client, err := MasterQ(r).ClientQ().FilterByName(userName.User).Get()
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

	tmps, err := MasterQ(r).TemplateQ().FilterByUser(client.ID).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to select templates ")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newTemlateListResp(tmps))

}

func newTemlateListResp(tmps []data.Template) resources.TemplateListResponse {
	responseTmpList := make([]resources.Template, 0)
	for _, tmp := range tmps {
		responseTmpList = append(responseTmpList, resources.Template{
			Attributes: resources.TemplateAttributes{
				BackgroundImg: string(tmp.ImgBytes),
				TemplateName:  tmp.Name,
				TemplateId:    tmp.ID,
			},
		})
	}
	return resources.TemplateListResponse{
		Data: responseTmpList,
	}
}
