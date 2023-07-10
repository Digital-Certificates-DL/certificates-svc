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
	"time"

	"strings"
)

var shortTitles = map[string]string{
	"Beginner at theoretical aspects blockchain technology": "blockchain",
	"Theory of database organization and basic SQL":         "database",
	"Cryptography and information security theory":          "security",
	"Біткоїн та криптовалюти":                               "bitcoin",
	"Basic Level in Decentralized Technologies":             "decentralize_technologies",
	"Golang programming on fundamental aspects":             "golang",
	"Blockchain and Distributed Systems":                    "distributed_system",
}

type QRCreator interface {
	GenerateQR(address []byte) ([]byte, []byte, string, error)
	PrepareMsgForQR(name string, address, signature []byte) string
	pngQR(msg string) ([]byte, error)
}

type QR struct {
	user              *helpers.Certificate
	templates         map[string]string
	log               *logan.Entry
	qrMessageTemplate string
}

func NewQR(user *helpers.Certificate, log *logan.Entry, templates map[string]string, qrMessageTemplate string) QRCreator {
	return QR{
		user:              user,
		log:               log,
		templates:         templates,
		qrMessageTemplate: qrMessageTemplate,
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
	return fmt.Sprintf(q.qrMessageTemplate, name, fmt.Sprintf("%s", address), fmt.Sprintf("%s", signature))

}

func (q QR) pngQR(msg string) ([]byte, error) {
	unixTime := time.Now().Unix()

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
	if err := qrcode.WriteColorFile(msg, qrcode.Highest, 400, back, front, fmt.Sprintf("testqr%d.png", unixTime)); err != nil {
		return nil, errors.Wrap(err, "failed to write color file")
	}
	defer os.Remove(fmt.Sprintf("testqr%d.png", unixTime))

	file, err := os.ReadFile(fmt.Sprintf("testqr%d.png", unixTime))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return file, nil
}
