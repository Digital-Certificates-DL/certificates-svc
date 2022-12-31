package service

import (
	"helper/internal/config"
	"log"
)

func Start(cfg config.Config) error {
	log.Println("start")
	//argsWithProg := os.Args
	//os.Mkdir("QRs", os.ModePerm)

	users, err := Parse(cfg.Table().Input)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, user := range users {
		hashing(user)
		GenerateQR(user, cfg.Key().Private, cfg.Google().Login, cfg.Google().Password, cfg.Google().SecretPath)
	}

	SetRes(users)
	return nil
}
