package main

import (
	"fmt"
	"log"
	"os"
)

var sample = "message:\n%s\n\naddress:\n%s\n\nsignature:\n%s\n\ncertificate page:\nhttps://dlt-academy.com/certificates"

func main() {
	log.Println("start")
	argsWithProg := os.Args
	os.Mkdir("QRs", os.ModePerm)

	users, err := Parse(argsWithProg[1])
	if err != nil {
		log.Println(err)
		return
	}

	for _, user := range users {
		hashing(user)
		GenerateQR(user, argsWithProg[2])
	}

	SetRes(users)
}

type user struct {
	Date                string
	Participant         string
	CourseTitle         string
	SerialNumber        string
	Note                string
	Certificate         string
	DataHash            string
	TxHash              string
	Signature           string
	DataCertificatePath string
}

func PrepareMsgForQR(name string, address, signature []byte) string {
	msg := fmt.Sprintf(sample, name, fmt.Sprintf("%x", address), fmt.Sprintf("%x", signature))
	return msg
}
