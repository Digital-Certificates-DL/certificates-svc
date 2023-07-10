package handlers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func CreateTemplate(w http.ResponseWriter, r *http.Request) {
	template, backgroundImg, req, err := requests.NewGenerateTemplate(r)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to generate template"))
		ape.Render(w, problems.BadRequest(err))
		return
	}
	defaultData := pdf.DefaultData
	client, err := MasterQ(r).ClientQ().GetByName(req.Data.Relationships.User)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to get client"))
		ape.Render(w, problems.InternalError())
		return
	}
	if client == nil {
		Log(r).Error(errors.Wrap(err, "client is not found"))
		ape.Render(w, problems.NotFound())
		return
	}

	if template.Width == 0 || template.High == 0 {
		tp := pdf.DefaultTemplateTall
		_, _, imgBytes, err := tp.Prepare(defaultData, pdf.NewPDFConfig(Config(r)), MasterQ(r), backgroundImg, client.ID)
		if err != nil {
			Log(r).Error(errors.Wrap(err, "failed to prepare pdf"))
			ape.Render(w, problems.InternalError())
			return
		}
		ape.Render(w, newTemplateImageResp(imgBytes))
		return
	}

	file := pdf.NewPDF(template.High, template.Width)

	file.SetName(template.Name.X, template.Name.Y, template.Name.FontSize, template.Name.Font)
	file.SetDate(template.Date.X, template.Date.Y, template.Date.FontSize, template.Date.Font)
	file.SetCourse(template.Course.X, template.Course.Y, template.Course.FontSize, template.Course.Font)
	file.SetCredits(template.Credits.X, template.Credits.Y, template.Credits.FontSize, template.Credits.Font)
	file.SetExam(template.Exam.X, template.Exam.Y, template.Exam.FontSize, template.Exam.Font)
	file.SetLevel(template.Level.X, template.Level.Y, template.Level.FontSize, template.Level.Font)
	file.SetSerialNumber(template.SerialNumber.X, template.SerialNumber.Y, template.SerialNumber.FontSize, template.SerialNumber.Font)
	file.SetPoints(template.Points.X, template.Points.Y, template.Points.FontSize, template.Points.Font)
	file.SetQR(template.QR.X, template.QR.Y, template.QR.FontSize, template.QR.High, template.Width)
	_, _, imgBytes, err := template.Prepare(defaultData, pdf.NewPDFConfig(Config(r)), MasterQ(r), backgroundImg, client.ID)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to prepare pdf"))
		ape.Render(w, problems.InternalError())
		return
	}
	if req.Data.Attributes.IsCompleted {
		templateBytes, err := json.Marshal(template)
		if err != nil {
			Log(r).Error(errors.Wrap(err, "failed to marshal"))
			ape.Render(w, problems.InternalError())
			return
		}

		_, err = MasterQ(r).TemplateQ().Insert(&data.Template{
			Template: templateBytes,
			//ImgBytes: backgroundImg,
			Name:   req.Data.Attributes.TemplateName,
			UserID: client.ID,
		})
		if err != nil {
			Log(r).Error(errors.Wrap(err, "failed to insert template"))
			ape.Render(w, problems.InternalError())
			return
		}
	}
	ape.Render(w, newTemplateImageResp(imgBytes))
	return
}

func newTemplateImageResp(img []byte) resources.TemplateResponse {
	return resources.TemplateResponse{
		Data: resources.Template{
			Attributes: resources.TemplateAttributes{
				BackgroundImg: base64.StdEncoding.EncodeToString(img),
			},
		},
	}
}
