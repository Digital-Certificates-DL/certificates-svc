package qr

import (
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"image/color"

	"github.com/boombuler/barcode/qr"
	qrcode "github.com/skip2/go-qrcode"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
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
	user *data.User
	cfg  config.Config
}

func NewQR(user *data.User, cfg config.Config) QR {
	return QR{
		user: user,
		cfg:  cfg,
	}
}

func (q QR) GenerateQR(address []byte) ([]byte, []byte, string, error) {
	parsedName := strings.Split(q.user.Participant, " ")
	path := ""
	q.cfg.Log().Debug(parsedName)
	if len(parsedName) < 2 {
		path = fmt.Sprintf("certificate_%s_%s_QR_codecreate.svg", parsedName[0], q.cfg.TemplatesConfig()[q.user.CourseTitle])
	} else {
		path = fmt.Sprintf("certificate_%s_%s_%s_QR_codecreate.svg", parsedName[0], parsedName[1], q.cfg.TemplatesConfig()[q.user.CourseTitle])
	}
	pathWithSuffix := fmt.Sprintf(q.cfg.QRCode().QRPath + path)
	fi, err := os.Create(pathWithSuffix)

	if err != nil {
		return nil, nil, "", errors.Wrap(err, "failed to create file by path")
	}
	s := svg.New(fi)
	msg := q.PrepareMsgForQR(q.user.Msg, []byte("1BooKnbm48Eabw3FdPgTSudt9u4YTWKBvf"), []byte(q.user.Signature))

	qrCode, _ := qr.Encode(msg, qr.M, qr.Auto)
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	s.End()

	img, err := q.pngQR(msg)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "failed to generate jpeg")
	}
	file, err := os.ReadFile(pathWithSuffix)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "failed to read file")
	}
	err = os.Remove(pathWithSuffix)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "failed to remove file")
	}
	return file, img, path, nil
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
