package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/handlers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/qr"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"log"
	"net/http"
)

const SENDQR = "qr"
const SENDCERTIFICATE = "certificate"

func PrepareCertificate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPrepareCertificates(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse data")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	users := req.PrepareUsers()
	var files []handlers.FilesBytes
	var filesCert []handlers.FilesBytes

	googleClient := google.NewGoogleClient(helpers.Config(r))
	link, err := googleClient.Connect(helpers.Config(r).Google().SecretPath, helpers.ClientQ(r), req.Data.Name)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	if len(link) != 0 {
		helpers.Log(r).WithError(err).Error("failed to authorize")

		ape.RenderErr(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "Invalid token",
			Status: "403",
			Code:   "125",
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}

	client, err := helpers.ClientQ(r).GetByName(req.Data.Name)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get client")
		ape.Render(w, problems.InternalError())
		return
	}
	if client == nil {
		helpers.Log(r).Error(errors.Wrap(err, "client is not found"))
		ape.RenderErr(w, problems.NotFound())
		return
	}

	for _, user := range users {
		qrData := qr.NewQR(user, helpers.Config(r))
		hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))

		if hash != "" {
			helpers.Log(r).Info(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		var file []byte
		name := ""
		file, img, name, err := qrData.GenerateQR([]byte(req.Data.Address))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to generate qrData")
			ape.Render(w, problems.InternalError())
			return
		}

		files = append(files, handlers.FilesBytes{File: file, Name: name, ID: user.ID, Type: "image/svg+xml"})

		req := pdf.DefaultTemplateNormal
		log.Println("user", user)
		certificate := pdf.NewPDF(req.High, req.Width)
		certificate.SetName(req.Name.X, req.Name.Y, req.Name.FontSize, req.Name.Font)
		certificate.SetDate(req.Date.X, req.Date.Y, req.Date.FontSize, req.Date.Font)
		certificate.SetCourse(req.Course.X, req.Course.Y, req.Course.FontSize, req.Course.Font)
		certificate.SetCredits(req.Credits.X, req.Credits.Y, req.Credits.FontSize, req.Credits.Font)
		certificate.SetExam(req.Exam.X, req.Exam.Y, req.Exam.FontSize, req.Exam.Font)
		certificate.SetLevel(req.Level.X, req.Level.Y, req.Level.FontSize, req.Level.Font)
		certificate.SetSerialNumber(req.SerialNumber.X, req.SerialNumber.Y, req.SerialNumber.FontSize, req.SerialNumber.Font)
		certificate.SetPoints(req.Points.X, req.Points.Y, req.Points.FontSize, req.Points.Font)
		certificate.SetQR(req.QR.X, req.QR.Y, req.QR.FontSize, req.QR.High, req.Width)

		pdfData := pdf.NewData(user.Participant, user.CourseTitle, "45 hours / 1.5 ECTS Credit", user.Points, user.SerialNumber, user.Date, img, user.Note, "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, helpers.Config(r), helpers.TemplateQ(r), nil, client.ID)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to create pdf")
			ape.Render(w, problems.BadRequest(err))
			return
		}
		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, handlers.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
	}

	users, err = handlers.Drive(googleClient, helpers.Log(r), files, users, SENDQR, helpers.Config(r).Google().QRPath)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	users, err = handlers.Drive(googleClient, helpers.Log(r), filesCert, users, SENDCERTIFICATE, helpers.Config(r).Google().PdfPath)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	helpers.Log(r).Info("creating table")
	errs := googleClient.SetRes(users, req.Data.Url)
	if errs != nil {
		helpers.Log(r).Error("failed to send date to drive: Errors: ", errs)
		ape.Render(w, problems.InternalError())
		return
	}

	ape.Render(w, newUserWithImgResponse(users))
}

func newUserWithImgResponse(users []*helpers.User) resources.UserListResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.USER,
			},
			Attributes: resources.UserAttributes{
				Participant:        user.Participant,
				Date:               user.Date,
				CourseTitle:        user.CourseTitle,
				CertificateImg:     user.ImageCertificate,
				DigitalCertificate: user.DigitalCertificate,
				Certificate:        user.Certificate,
				Points:             user.Points,
				Note:               user.Note,
				Signature:          user.Signature,
			},
		}
		usersData = append(usersData, resp)
	}

	return resources.UserListResponse{
		Data: usersData,
	}

}

func newLinkResponse(link string) resources.LinkResponse {
	data := resources.LinkResponse{
		Data: resources.Link{
			Attributes: resources.LinkAttributes{
				Link: link,
			},
		},
	}
	return data
}
