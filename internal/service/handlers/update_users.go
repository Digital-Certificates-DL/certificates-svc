package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/handlers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"net/http"
)

func UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPrepareCertificates(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}
	users := req.PrepareUsers()

	var filesCert []handlers.FilesBytes

	client := google.NewGoogleClient(helpers.Config(r))
	link, err := client.Connect(helpers.Config(r).Google().SecretPath, helpers.ClientQ(r), req.Data.Name)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	if len(link) != 0 {
		helpers.Log(r).WithError(err).Error("failed to authorize")
		ape.Render(w, newLinkResponse(link))
		w.WriteHeader(403)
		return
	}

	for _, user := range users {
		//qrData := qr.NewQR(user, helpers.Config(r))
		//_, img, name, err := qrData.GenerateQR([]byte(req.Data.Address))
		//if err != nil {
		//	helpers.Log(r).WithError(err).Error("failed to generate qrData")
		//	ape.Render(w, problems.InternalError())
		//	return
		//}

		//hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))
		//if hash != "" {
		//	helpers.Log(r).Info(user.Participant, " hash = ", hash)
		//}
		//
		//user.SetDataHash(hash)

		if user.TxHash != "" {
			user.SetDataHash(user.TxHash)
		}

		req := pdf.DefaultTemplateTall
		helpers.Log(r).Info("user", user)
		certificate := pdf.NewPDF(req.High, req.Width)
		//certificate.SetName(req.Name.X, req.Name.Y, req.Name.Size, req.Name.Font)
		//certificate.SetDate(req.Date.X, req.Date.Y, req.Date.Size, req.Date.Font)
		//certificate.SetCourse(req.Course.X, req.Course.Y, req.Course.Size, req.Course.Font)
		//certificate.SetCredits(req.Credits.X, req.Credits.Y, req.Credits.Size, req.Credits.Font)
		//certificate.SetExam(req.Exam.X, req.Exam.Y, req.Exam.Size, req.Exam.Font)
		//certificate.SetLevel(req.Level.X, req.Level.Y, req.Level.Size, req.Level.Font)
		certificate.SetSerialNumber(req.SerialNumber.X, req.SerialNumber.Y, req.SerialNumber.Size, req.SerialNumber.Font)
		//certificate.SetPoints(req.Points.X, req.Points.Y, req.Points.Size, req.Points.Font)
		//certificate.SetQR(req.QR.X, req.QR.Y, req.QR.Size, req.QR.High, req.Width)

		pdfData := pdf.NewData("", "", "", "", user.SerialNumber, "", nil, "", "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, helpers.Config(r), helpers.TemplateQ(r), user.ImageCertificate)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to create pdf")
			ape.Render(w, problems.BadRequest(err))
			return
		}
		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, handlers.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
	}

	users, err = handlers.Drive(client, helpers.Log(r), filesCert, users, SENDCERTIFICATE, helpers.Config(r).Google().PdfPath)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	helpers.Log(r).Info("creating table")
	errs := client.SetRes(users, req.Data.Url)
	if errs != nil {
		helpers.Log(r).Error("failed to send date to drive: Errors: ", errs)
		ape.Render(w, problems.InternalError())
		return
	}

	ape.Render(w, newUserWithImgResponse(users))
}
