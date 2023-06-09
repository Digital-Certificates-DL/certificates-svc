package handlers

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"net/http"
)

func CreateTemplate(w http.ResponseWriter, r *http.Request) {
	template, backgroundImg, err := requests.NewGenerateTemplate(r)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to generate template"))
		ape.Render(w, problems.BadRequest(err))
		return
	}
	file := pdf.NewPDF(template.High, template.Width)

	file.SetName(template.Name.X, template.Name.Y, template.Name.Size, template.Name.Font)
	file.SetDate(template.Date.X, template.Date.Y, template.Date.Size, template.Date.Font)
	file.SetCourse(template.Course.X, template.Course.Y, template.Course.Size, template.Course.Font)
	file.SetCredits(template.Credits.X, template.Credits.Y, template.Credits.Size, template.Credits.Font)
	file.SetExam(template.Exam.X, template.Exam.Y, template.Exam.Size, template.Exam.Font)
	file.SetLevel(template.Level.X, template.Level.Y, template.Level.Size, template.Level.Font)
	file.SetSerialNumber(template.SerialNumber.X, template.SerialNumber.Y, template.SerialNumber.Size, template.SerialNumber.Font)
	file.SetPoints(template.Points.X, template.Points.Y, template.Points.Size, template.Points.Font)
	file.SetQR(template.QR.X, template.QR.Y, template.QR.Size, template.QR.High, template.Width)
	d := pdf.DefaultData

	_, _, imgBytes, err := file.Prepare(d, helpers.Config(r), helpers.TemplateQ(r), backgroundImg)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to prepare pdf"))
		ape.Render(w, problems.InternalError())
		return
	}

	ape.Render(w, imgBytes) //todo wrap it

}
