package service

import (
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"helper/internal/data"
	"helper/internal/service/google"
	"log"
	"os"
	"strings"
)

const sample = "message:\n%s\n\naddress:\n%s\n\nsignature:\n%s\n\ncertificate page:\nhttps://dlt-academy.com/certificates"

func GenerateQR(user *data.User, key string, login, password string, secretPath string) {
	parsedName := strings.Split(user.Participant, " ")
	path := ""
	if len(parsedName) < 2 {
		path = fmt.Sprintf("QRs/certificate_%s_QR_codecreate.svg", parsedName[0])
	} else {
		path = fmt.Sprintf("QRs/certificate_%s_%s_QR_codecreate.svg", parsedName[0], parsedName[1])
	}
	fi, _ := os.Create(path)
	s := svg.New(fi)
	aggregatedStr := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
	signature, _, address, err := Sign(key, aggregatedStr)
	if err != nil {
		log.Println(err)
	}

	user.Signature = string(signature)
	qrCode, _ := qr.Encode(PrepareMsgForQR(aggregatedStr, address, signature), qr.M, qr.Auto)
	google.Update(aggregatedStr, secretPath, login, password)
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)
	user.DataCertificatePath = path

	s.End()
}

func PrepareMsgForQR(name string, address, signature []byte) string {
	msg := fmt.Sprintf(sample, name, fmt.Sprintf("%s", address), fmt.Sprintf("%s", signature))
	return msg
}
