package qr

import (
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"helper/internal/config"
	"helper/internal/data"
	"helper/internal/service/signature"
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
	sign signature.Signature
	cfg  config.Config
}

func NewQR(user *data.User, cfg config.Config, sign signature.Signature) QR {
	return QR{
		user: user,
		sign: sign,
		cfg:  cfg,
	}
}

func (q QR) GenerateQR() (string, string, string, error) {

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
		return "", "", "", errors.Wrap(err, "failed to create file by path")
	}
	s := svg.New(fi)
	aggregatedStr := fmt.Sprintf("%s %s %s", q.user.Date, q.user.Participant, q.user.CourseTitle)

	signedMsg, _, address, err := q.sign.Sign(aggregatedStr)
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to sign msg")
	}

	qrCode, _ := qr.Encode(q.PrepareMsgForQR(aggregatedStr, address, signedMsg), qr.M, qr.Auto)

	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)
	s.End()
	return path, pathWithSuffix, string(signedMsg), nil
}

func (q QR) PrepareMsgForQR(name string, address, signature []byte) string {
	msg := fmt.Sprintf(sample, name, fmt.Sprintf("%s", address), fmt.Sprintf("%s", signature))
	return msg
}
