package service

import (
	"helper/internal/config"
	"helper/internal/service/google"
	"helper/internal/service/signature"
	"log"
	"os"
)

func Start(cfg config.Config) error {
	log.Println("start")
	os.MkdirAll(cfg.QRCode().QRPath, os.ModePerm)

	users, err := Parse(cfg.Table().Input)
	if err != nil {
		log.Println(err)
		return err
	}
	connect := google.Connect(cfg.Google().SecretPath, cfg.Google().Code)

	folderIDList, err := google.CreateFolder(connect, cfg.Google().QRPath)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, user := range users {
		if user.Signature != "" && user.SerialNumber != "" {
			continue
		}
		signature.Hashing(user)
		GenerateQR(user, cfg.Key().Private, connect, folderIDList)
	}

	SetRes(users, cfg.Table().Result)
	return nil
}
