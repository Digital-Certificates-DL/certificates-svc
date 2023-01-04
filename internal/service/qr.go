package service

import (
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"helper/internal/data"
	"helper/internal/service/google"
	"log"
	"net/http"
	"os"
	"strings"
)

const sample = "message:\n%s\n\naddress:\n%s\n\nsignature:\n%s\n\ncertificate page:\nhttps://dlt-academy.com/certificates"

var shortTitles = map[string]string{
	"Beginner at theoretical aspects blockchain technology": "blockchain",
	"Theory of database organization and basic SQL":         "database",
	"Cryptography and information security theory":          "security",
	"Golang": "golang", //todo rename key
}

func GenerateQR(user *data.User, key string, client *http.Client, folderIDList []string) {

	parsedName := strings.Split(user.Participant, " ")
	path := ""
	if len(parsedName) < 2 {
		path = fmt.Sprintf("certificate_%s_QR_codecreate.svg", parsedName[0])
	} else {
		path = fmt.Sprintf("certificate_%s_%s_QR_codecreate.svg", parsedName[0], parsedName[1])
	}

	pathWithSuffix := fmt.Sprintf("./qr/%s", path)

	fi, err := os.Create(pathWithSuffix)
	if err != nil {
		log.Println(err)
		return
	}
	s := svg.New(fi)
	aggregatedStr := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
	signature, _, address, err := Sign(key, aggregatedStr)
	if err != nil {
		log.Println(err)
		return
	}

	user.Signature = string(signature)
	qrCode, _ := qr.Encode(PrepareMsgForQR(aggregatedStr, address, signature), qr.M, qr.Auto)

	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	user.CertificatePath = path

	s.End()

	link, err := google.Update(path, client, folderIDList)
	if err != nil {
		log.Println(err)
		return
	}
	user.CertificatePath = link

}

func PrepareMsgForQR(name string, address, signature []byte) string {
	msg := fmt.Sprintf(sample, name, fmt.Sprintf("%s", address), fmt.Sprintf("%s", signature))
	return msg
}
