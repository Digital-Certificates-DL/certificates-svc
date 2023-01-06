package service

import (
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"helper/internal/config"
	"helper/internal/data"
	"helper/internal/service/signature"
	"log"
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

func GenerateQR(user *data.User, cfg config.Config) (string, error) {

	parsedName := strings.Split(user.Participant, " ")
	path := ""
	if len(parsedName) < 2 {
		path = fmt.Sprintf("certificate_%s_%s_QR_codecreate.svg", parsedName[0], cfg.TemplatesConfig()[user.CourseTitle])
	} else {
		path = fmt.Sprintf("certificate_%s_%s_%s_QR_codecreate.svg", parsedName[0], parsedName[1], cfg.TemplatesConfig()[user.CourseTitle])
	}

	pathWithSuffix := fmt.Sprintf(cfg.QRCode().QRPath, path)

	fi, err := os.Create(pathWithSuffix)
	if err != nil {
		log.Println(err)
		return "", err
	}
	s := svg.New(fi)
	aggregatedStr := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
	signature, _, address, err := signature.Sign(cfg.Key().Private, aggregatedStr)
	if err != nil {
		log.Println(err)
		return "", err
	}

	user.Signature = string(signature)
	qrCode, _ := qr.Encode(PrepareMsgForQR(aggregatedStr, address, signature), qr.M, qr.Auto)

	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	user.CertificatePath = path

	s.End()

	return path, nil
}

func PrepareMsgForQR(name string, address, signature []byte) string {
	msg := fmt.Sprintf(sample, name, fmt.Sprintf("%s", address), fmt.Sprintf("%s", signature))
	return msg
}
