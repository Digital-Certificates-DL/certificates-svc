package pdf

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/handlers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/qr"
)

type Container struct {
	users   []*helpers.User
	number  int
	status  bool
	log     *logan.Entry
	client  *data.Client
	config  config.Config
	masterQ data.MasterQ
}

func (c *Container) run() error {
	var filesCert []handlers.FilesBytes
	for _, user := range c.users {
		qrData := qr.NewQR(user, helpers.Config(r))
		_, img, name, err := qrData.GenerateQR([]byte(req.Data.Address))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to generate qrData")
			ape.Render(w, problems.InternalError())
			return
		}

		hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))
		if hash != "" {
			helpers.Log(r).Info(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		c.log.Info(user.TxHash, " tx")
		if user.TxHash != "" {
			user.SetDataHash(user.TxHash)
			c.log.Debug(user.DataHash)
			c.log.Debug(user.TxHash)
		}

		req := DefaultTemplateTall
		c.log.Debug("user", user)
		certificate := NewPDF(req.High, req.Width)
		certificate.SetSerialNumber(req.SerialNumber.X, req.SerialNumber.Y, req.SerialNumber.FontSize, req.SerialNumber.Font)

		pdfData := NewData("", user.CourseTitle, "", "", user.SerialNumber, "", nil, "", "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, c.config, c.masterQ.TemplateQ(), user.ImageCertificate, c.client.ID)
		if err != nil {
			return errors.Wrap(err, "failed to generate pdf")
		}

		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, handlers.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
	}

	users, err := handlers.Drive(googleClient, Log(r), filesCert, users, SENDCERTIFICATE, Config(r).Google().PdfPath)
	if err != nil {
		Log(r).WithError(err).Error("failed to send date to drive")
		ape.Render(w, problems.InternalError())
		return
	}

	Log(r).Info("creating table")
	errs := googleClient.SetRes(users, req.Data.Url)
	if errs != nil {
		Log(r).Error("failed to send date to drive: Errors: ", errs)
		ape.Render(w, problems.InternalError())
		return
	}

	return nil
}
