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
	client       *data.Client
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
		certificate := c.setTemplateData(req)

		pdfData := NewData(user.Participant, user.CourseTitle, "45 hours / 1.5 ECTS Credit", user.Points, user.SerialNumber, user.Date, img, user.Note, "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, c.config, c.masterQ, nil, c.owner.ID)
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

		certificate := c.setTemplateData(req)

		pdfData := NewData(user.Participant, user.CourseTitle, "45 hours / 1.5 ECTS Credit", user.Points, user.SerialNumber, user.Date, img, user.Note, "", "")
		fileBytes, name, certificateImg, err := certificate.Prepare(pdfData, c.config, c.masterQ, nil, c.owner.ID)
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

func (c *Container) setTemplateData(template PDF) *PDF {
	certificate := NewPDF(template.High, template.Width)
	certificate.SetName(template.Name.X, template.Name.Y, template.Name.FontSize, template.Name.Font)
	certificate.SetDate(template.Date.X, template.Date.Y, template.Date.FontSize, template.Date.Font)
	certificate.SetCourse(template.Course.X, template.Course.Y, template.Course.FontSize, template.Course.Font)
	certificate.SetCredits(template.Credits.X, template.Credits.Y, template.Credits.FontSize, template.Credits.Font)
	certificate.SetExam(template.Exam.X, template.Exam.Y, template.Exam.FontSize, template.Exam.Font)
	certificate.SetLevel(template.Level.X, template.Level.Y, template.Level.FontSize, template.Level.Font)
	certificate.SetSerialNumber(template.SerialNumber.X, template.SerialNumber.Y, template.SerialNumber.FontSize, template.SerialNumber.Font)
	certificate.SetPoints(template.Points.X, template.Points.Y, template.Points.FontSize, template.Points.Font)
	certificate.SetQR(template.QR.X, template.QR.Y, template.QR.FontSize, template.QR.High, template.Width)
	return certificate
}
