package main

import (
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"log"
	"os"
	"strings"
)

func GenerateQR(user *user, key string) {
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

	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)
	user.DataCertificatePath = path

	s.End()
}
