package handlers

import (
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/handlers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/qr"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"log"
	"net/http"
	"os"
)

const SENDQR = "qr"
const SENDCERTIFICATE = "certificate"

func PrepareCertificate(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPrepareCertificates(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}
	users := request.PrepareUsers()
	var files []handlers.FilesBytes
	var filesCert []handlers.FilesBytes

	client := google.NewGoogleClient(helpers.Config(r))
	err = client.Connect(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	os.MkdirAll(helpers.Config(r).QRCode().QRPath, os.ModePerm)
	defer os.RemoveAll(helpers.Config(r).QRCode().QRPath)
	for _, user := range users {
		qr := qr.NewQR(user, helpers.Config(r))
		hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))

		if hash != "" {
			helpers.Log(r).Info(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		var file []byte
		name := ""
		file, img, name, err := qr.GenerateQR([]byte(request.Data.Address))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to generate qr")
			ape.Render(w, problems.InternalError())
			return
		}

		files = append(files, handlers.FilesBytes{File: file, Name: name, ID: user.ID, Type: "image/svg+xml"})

		req := pdf.DefaultTemplateNormal
		log.Println("user", user)
		certificate := pdf.NewPDF(req.High, req.Width)
		certificate.SetName(req.Name.X, req.Name.Y, req.Name.Size, req.Name.Font)
		certificate.SetDate(req.Date.X, req.Date.Y, req.Date.Size, req.Date.Font)
		certificate.SetCourse(req.Course.X, req.Course.Y, req.Course.Size, req.Course.Font)
		certificate.SetCredits(req.Credits.X, req.Credits.Y, req.Credits.Size, req.Credits.Font)
		certificate.SetExam(req.Exam.X, req.Exam.Y, req.Exam.Size, req.Exam.Font)
		certificate.SetLevel(req.Level.X, req.Level.Y, req.Level.Size, req.Level.Font)
		//certificate.SetNote(req.Note.X, req.Note.Y, req.Note.Size, req.Note.Font)
		certificate.SetSerialNumber(req.SerialNumber.X, req.SerialNumber.Y, req.SerialNumber.Size, req.SerialNumber.Font)
		certificate.SetPoints(req.Points.X, req.Points.Y, req.Points.Size, req.Points.Font)
		certificate.SetQR(req.QR.X, req.QR.Y, req.QR.Size, req.QR.High, req.Width)

		//credits, point := certificate.ParsePoints(user.Points)

		data := pdf.NewData(user.Participant, user.CourseTitle, "45 hours / 1.5 ECTS Credit", user.Points, user.SerialNumber, user.Date, img, user.Note, "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(data, helpers.Config(r))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to create pdf")
			ape.Render(w, problems.BadRequest(err))
			return
		}
		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, handlers.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
	}

	users, err = handlers.Drive(client, helpers.Log(r), files, users, SENDQR, helpers.Config(r).Google().QRPath)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	users, err = handlers.Drive(client, helpers.Log(r), filesCert, users, SENDCERTIFICATE, helpers.Config(r).Google().PdfPath)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	helpers.Log(r).Info("creating table")
	errs := client.SetRes(users, request.Data.Url)
	if errs != nil {
		helpers.Log(r).Error("failed to send date to drive: Errors: ", errs)
		ape.Render(w, problems.InternalError())
		return
	}

	ape.Render(w, newUserWithImgResponse(users))
}

func newUserWithImgResponse(users []*data.User) resources.UserListResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.USER,
			},
			Attributes: resources.UserAttributes{
				Participant:    user.Participant,
				Date:           user.Date,
				CourseTitle:    user.CourseTitle,
				CertificateImg: user.ImageCertificate,
			},
		}
		usersData = append(usersData, resp)
	}

	return resources.UserListResponse{
		Data: usersData,
	}

}
