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
	req, err := requests.NewGenerateTemplate(r)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to generate template"))
		ape.Render(w, problems.BadRequest(err))
		return
	}
	file := pdf.NewPDF(req.Data.High, req.Data.Width)

	file.SetName(req.Data.Name.X, req.Data.Name.Y, req.Data.Name.Size, req.Data.Name.Font)
	file.SetDate(req.Data.Date.X, req.Data.Date.Y, req.Data.Date.Size, req.Data.Date.Font)
	file.SetCourse(req.Data.Course.X, req.Data.Course.Y, req.Data.Course.Size, req.Data.Course.Font)
	file.SetCredits(req.Data.Credits.X, req.Data.Credits.Y, req.Data.Credits.Size, req.Data.Credits.Font)
	file.SetExam(req.Data.Exam.X, req.Data.Exam.Y, req.Data.Exam.Size, req.Data.Exam.Font)
	file.SetLevel(req.Data.Level.X, req.Data.Level.Y, req.Data.Level.Size, req.Data.Level.Font)
	file.SetNote(req.Data.Note.X, req.Data.Note.Y, req.Data.Note.Size, req.Data.Note.Font)
	file.SetPoints(req.Data.Points.X, req.Data.Points.Y, req.Data.Points.Size, req.Data.Points.Font)
	file.SetNote(req.Data.Note.X, req.Data.Note.Y, req.Data.Note.Size, req.Data.Note.Font)
	file.SetQR(req.Data.QR.X, req.Data.QR.Y, req.Data.QR.Size, req.Data.QR.High, req.Data.Width)

	pdfBytes, _, err := file.Prepare(pdf.DefaultData, helpers.Config(r))
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to prepare pdf"))
		ape.Render(w, problems.InternalError())
		return
	}
	//image, err := file.PDFToImg(pdfBytes)
	//if err != nil {
	//	helpers.Log(r).Error(errors.Wrap(err, "failed to convert pdf to image"))
	//	ape.Render(w, problems.InternalError())
	//	return
	//}
	//ape.Render(w, image)
	ape.Render(w, pdfBytes)

}
