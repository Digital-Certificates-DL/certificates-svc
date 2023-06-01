package handlers

import (
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
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
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}

	client, err := helpers.ClientQ(r).GetByName(req.Data.Name) //todo use relationship
	helpers.Log(r).Debug("user ", client)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to get user"))
		ape.Render(w, problems.InternalError())
		return
	}
	if client == nil {
		helpers.Log(r).Error(errors.Wrap(err, "user is not found"))
		ape.Render(w, problems.NotFound())
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
		helpers.Log(r).Info(user.TxHash, " tx")
		if user.TxHash != "" {
			user.SetDataHash(user.TxHash)
			helpers.Log(r).Info(user.DataHash)
			helpers.Log(r).Info(user.TxHash)
		}

		req := pdf.DefaultTemplateTall
		helpers.Log(r).Info("user", user)
		certificate := pdf.NewPDF(req.High, req.Width)
		certificate.SetSerialNumber(req.SerialNumber.X, req.SerialNumber.Y, req.SerialNumber.FontSize, req.SerialNumber.Font)

		pdfData := pdf.NewData("", user.CourseTitle, "", "", user.SerialNumber, "", nil, "", "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, helpers.Config(r), helpers.TemplateQ(r), user.ImageCertificate, client.ID)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to create pdf")
			ape.Render(w, problems.BadRequest(err))
			return
		}
		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, handlers.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
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
