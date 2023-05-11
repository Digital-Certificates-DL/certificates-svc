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

func GetTemplates(w http.ResponseWriter, r *http.Request) {
	userName, err := requests.NewGetTemplateRequest(r)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to parse request "))
		ape.Render(w, problems.BadRequest(err))
		return
	}

	tmps, err := helpers.TemplateQ(r).FilterByUser(userName.User).Select()
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to select templates "))
		ape.Render(w, problems.BadRequest(err))
		return
	}
	ape.Render(w, newTemlateListResp(tmps))

}

func newTemlateListResp(tmps []data.Template) resources.TemplateListResponse {
	var reponseTmpList []resources.Template
	for _, tmp := range tmps {
		reponseTmpList = append(reponseTmpList, resources.Template{
			Attributes: resources.TemplateAttributes{
				BackgroundImg: base64.StdEncoding.EncodeToString(tmp.ImgBytes),
				TemplateName:  tmp.Name,
			},
		})
	}
	return resources.TemplateListResponse{
		Data: reponseTmpList,
	}
}
