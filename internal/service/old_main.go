package service

import (
	"helper/internal/config"
	"helper/internal/service/google"
	"log"
	"os"
)

func Start(cfg config.Config) error {
	log.Println("start")
	//argsWithProg := os.Args
	//os.Mkdir("QRs", os.ModePerm)
	os.MkdirAll("./qr", os.ModePerm)

	users, err := Parse(cfg.Table().Input)
	if err != nil {
		log.Println(err)
		return err
	}
	connect := google.Connect(cfg.Google().SecretPath)

	folderIDList, err := google.CreateFolder(connect)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, user := range users {
		hashing(user)
		GenerateQR(user, cfg.Key().Private, connect, folderIDList)

		//GenerateQR(user, cfg.Key().Private, cfg.Google().Login, cfg.Google().Password, cfg.Google().SecretPath)
	}

	SetRes(users)
	return nil
}
