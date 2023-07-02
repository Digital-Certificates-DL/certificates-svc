package qr

import (
	"bytes"
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	qrcode "github.com/skip2/go-qrcode"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"image/color"
	"os"

	"strings"
)

const sample = "message:\n%s\n\naddress:\n%s\n\nsignature:\n%s\n\ncertificate page:\nhttps://dlt-academy.com/certificates"

var shortTitles = map[string]string{
	"Beginner at theoretical aspects blockchain technology": "blockchain",
	"Theory of database organization and basic SQL":         "database",
	"Cryptography and information security theory":          "security",
	"Біткоїн та криптовалюти":                               "bitcoin",
	"Basic Level in Decentralized Technologies":             "decentralize_technologies",
	"Golang programming on fundamental aspects":             "golang",
	"Blockchain and Distributed Systems":                    "distributed_system",
}

type QR struct {
	user      *helpers.User
	templates map[string]string
	log       *logan.Entry
}

func NewQR(user *helpers.User, log *logan.Entry, templates map[string]string) QR {
	return QR{
		user:      user,
		log:       log,
		templates: templates,
	}
}

func (q QR) GenerateQR(address []byte) ([]byte, []byte, string, error) {
	parsedName := strings.Split(q.user.Participant, " ")
	path := ""
	q.log.Debug(parsedName)
	if len(parsedName) < 2 {
		path = fmt.Sprintf("certificate_%s_%s_QR_codecreate.svg", parsedName[0], q.templates[q.user.CourseTitle])
	} else {
		path = fmt.Sprintf("certificate_%s_%s_%s_QR_codecreate.svg", parsedName[0], parsedName[1], q.templates[q.user.CourseTitle])
	}

	bf := new(bytes.Buffer)
	s := svg.New(bf)
	msg := q.PrepareMsgForQR(q.user.Msg, address, []byte(q.user.Signature))

	qrCode, _ := qr.Encode(msg, qr.M, qr.Auto)
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	err := qs.WriteQrSVG(s)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "failed write qr in svg")
	}
	s.End()

	img, err := q.pngQR(msg)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "failed to generate jpeg")
	}

	return bf.Bytes(), img, path, nil
}

func (q QR) PrepareMsgForQR(name string, address, signature []byte) string {
	msg := fmt.Sprintf(sample, name, fmt.Sprintf("%s", address), fmt.Sprintf("%s", signature))
	return msg
}

func (q QR) pngQR(msg string) ([]byte, error) {
	back := color.RGBA{
		R: 0,
		G: 18,
		B: 54,
		A: 255,
	}
	front := color.RGBA{
		R: 63,
		G: 151,
		B: 255,
		A: 255,
	}
	err := qrcode.WriteColorFile(msg, qrcode.Highest, 400, back, front, "testqr.png")
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile("testqr.png")
	if err != nil {
		return nil, err
	}
	return file, nil
}
