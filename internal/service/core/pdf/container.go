package pdf

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/qr"
	"log"
)

type ContainerHandler interface {
	Generate() error
	Update() error
}

type Container struct {
	Users        []*helpers.User
	ID           int
	Status       bool
	log          *logan.Entry
	config       config.Config
	masterQ      data.MasterQ
	googleClient *google.Google
	address      string
	sheetUrl     string
	owner        *data.Client
	process      string
}

const SendQR = "qr"
const sendCertificate = "certificate"

func (c *Container) Generate() error {
	var files []google.FilesBytes
	var filesCert []google.FilesBytes
	for _, user := range c.Users {
		qrData := qr.NewQR(user, c.log, c.config.TemplatesConfig())
		hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))

		if hash != "" {
			c.log.Debug(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		file, img, name, err := qrData.GenerateQR([]byte(c.address))
		if err != nil {
			return errors.Wrap(err, "failed to Generate qrData")
		}

		files = append(files, google.FilesBytes{File: file, Name: name, ID: user.ID, Type: "image/svg+xml"})

		req := DefaultTemplateTall
		log.Println(req)
		log.Println("user", user)

		pdf := PDF{}
		certificate := pdf.SetTemplateData(req)

		pdfData := NewData(user.Participant, user.CourseTitle, "45 hours / 1.5 ECTS Credit", user.Points, user.SerialNumber, user.Date, img, user.Note, "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, NewPDFConfig(c.config), c.masterQ, nil, c.owner.ID)
		if err != nil {
			return errors.Wrap(err, "failed to create pdf")
		}
		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, google.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
	}

	users, err := google.Drive(c.googleClient, c.log, files, c.Users, SendQR, c.config.Google().QRPath)
	if err != nil {
		return errors.Wrap(err, "failed to send date to drive")
	}

	users, err = google.Drive(c.googleClient, c.log, filesCert, c.Users, sendCertificate, c.config.Google().PdfPath)
	if err != nil {
		return errors.Wrap(err, "failed to send date to drive")
	}

	c.log.Debug("creating table")
	errs := c.googleClient.SetRes(users, c.sheetUrl)
	if errs != nil {
		return errors.Wrap(err, "failed to set result on table")
	}
	c.Status = true

	return nil
}

func (c *Container) Update() error {
	var files []google.FilesBytes
	var filesCert []google.FilesBytes
	for _, user := range c.Users {
		qrData := qr.NewQR(user, c.log, c.config.TemplatesConfig())
		hash := user.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))

		if hash != "" {
			c.log.Debug(user.Participant, " hash = ", hash)
		}

		user.SetDataHash(hash)
		file, img, name, err := qrData.GenerateQR([]byte(c.address))
		if err != nil {
			return errors.Wrap(err, "failed to Generate qrData")
		}

		files = append(files, google.FilesBytes{File: file, Name: name, ID: user.ID, Type: "image/svg+xml"})

		req := DefaultTemplateTall
		log.Println(req)
		log.Println("user", user)

		pdf := PDF{}
		certificate := pdf.SetTemplateData(req)

		pdfData := NewData(user.Participant, user.CourseTitle, "45 hours / 1.5 ECTS Credit", user.Points, user.SerialNumber, user.Date, img, user.Note, "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, NewPDFConfig(c.config), c.masterQ, nil, c.owner.ID)
		if err != nil {
			return errors.Wrap(err, "failed to create pdf")
		}

		user.ImageCertificate = certificateImg
		filesCert = append(filesCert, google.FilesBytes{File: fileBytes, Name: name, ID: user.ID, Type: "application/pdf"})
	}

	users, err := google.Drive(c.googleClient, c.log, files, c.Users, SendQR, c.config.Google().QRPath)
	if err != nil {
		return errors.Wrap(err, "failed to send date to drive")
	}

	users, err = google.Drive(c.googleClient, c.log, filesCert, c.Users, sendCertificate, c.config.Google().PdfPath)
	if err != nil {
		return errors.Wrap(err, "failed to send date to drive")
	}

	c.log.Debug("creating table")
	errs := c.googleClient.SetRes(users, c.sheetUrl)
	if errs != nil {
		return errors.Wrap(err, "failed to set result on table")
	}
	c.Status = true

	return nil
}