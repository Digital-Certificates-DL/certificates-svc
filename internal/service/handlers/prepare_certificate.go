package handlers

import (
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"helper/internal/data"
	"helper/internal/handlers"
	"helper/internal/service/google"
	"helper/internal/service/helpers"
	"helper/internal/service/pdf"
	"helper/internal/service/qr"
	"helper/internal/service/requests"
	"log"
	"net/http"
	"os"
)

const SENDQR = "qr"
const SENDCERTIFICATE = "certificate"

func PrepareCertificate(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPrepareCertificates(r)

	var usersResult []*data.User
	var files []handlers.FilesBytes
	var filesCert []handlers.FilesBytes

	client := google.NewGoogleClient(helpers.Config(r))
	err = client.Connect(helpers.Config(r).Google().SecretPath, helpers.Config(r).Google().Code)
	if err != nil {
		log.Println(err)
		return
	}

	os.MkdirAll(helpers.Config(r).QRCode().QRPath, os.ModePerm)
	defer os.RemoveAll(helpers.Config(r).QRCode().QRPath)
	for id, user := range request.Data {
		user.ID = id
		if user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" || user.SerialNumber != "" {
			helpers.Log(r).Debug("has already")
			//todo maybe add render
			continue
		}
		log.Println(user)
		qr := qr.NewQR(user, helpers.Config(r))
		hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))

		if hash != "" {
			helpers.Log(r).Info(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		var file []byte
		name := ""
		file, name, err = qr.GenerateQR(request.Address)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to generate qr")
			ape.Render(w, problems.InternalError())
			return
		}

		files = append(files, handlers.FilesBytes{File: file, Name: name, ID: id, Type: "image/svg+xml"})

		req := pdf.DefaultTemplate
		log.Println("user", user)
		certificate := pdf.NewPDF(req.High, req.Width)
		certificate.SetName(req.Name.X, req.Name.Y, req.Name.Size, req.Name.Font)
		certificate.SetDate(req.Date.X, req.Date.Y, req.Date.Size, req.Date.Font)
		certificate.SetCourse(req.Course.X, req.Course.Y, req.Course.Size, req.Course.Font)
		certificate.SetCredits(req.Credits.X, req.Credits.Y, req.Credits.Size, req.Credits.Font)
		certificate.SetExam(req.Exam.X, req.Exam.Y, req.Exam.Size, req.Exam.Font)
		certificate.SetLevel(req.Level.X, req.Level.Y, req.Level.Size, req.Level.Font)
		certificate.SetNote(req.Note.X, req.Note.Y, req.Note.Size, req.Note.Font)
		certificate.SetPoints(req.Points.X, req.Points.Y, req.Points.Size, req.Points.Font)
		certificate.SetNote(req.Note.X, req.Note.Y, req.Note.Size, req.Note.Font)
		certificate.SetQR(req.QR.X, req.QR.Y, req.QR.Size, req.QR.High, req.Width)

		credits, point := certificate.ParsePoints(user.Points)
		log.Println(credits, point)
		data := pdf.NewData(user.Participant, user.CourseTitle, credits, point, user.SerialNumber, user.Date, user.DigitalCertificate, user.Note, "", "")
		fileBytes, name, err := certificate.Prepare(data, helpers.Config(r))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to create pdf")
			ape.Render(w, problems.BadRequest(err))
			return
		}
		filesCert = append(filesCert, handlers.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
		usersResult = append(usersResult, user)
	}

	request.Data, err = handlers.Drive(client, helpers.Config(r), helpers.Log(r), files, request.Data, SENDQR)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	request.Data, err = handlers.Drive(client, helpers.Config(r), helpers.Log(r), filesCert, request.Data, SENDCERTIFICATE)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	helpers.Log(r).Info("creating table")
	errs := client.SetRes(usersResult, request.Url)
	if errs != nil {
		helpers.Log(r).Error("failed to send date to drive: Errors: ", errs)
		ape.Render(w, problems.InternalError())
		return
	}
	return
}
