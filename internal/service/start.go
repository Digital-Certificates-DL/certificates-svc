package service

import (
	"context"
	"fmt"
	"gitlab.com/distributed_lab/running"
	"helper/internal/config"
	"helper/internal/service/google"
	"helper/internal/service/signature"
	"os"
	"time"
)

func Start(cfg config.Config) error {
	log := cfg.Log()
	log.Info("Start")

	os.MkdirAll(cfg.QRCode().QRPath, os.ModePerm)

	users, err := Parse(cfg.Table().Input)
	if err != nil {
		cfg.Log().Debug(err)
		return err
	}

	sendToDrive := cfg.Google().Enable
	connect, sendToDrive := google.Connect(cfg.Google().SecretPath, cfg.Google().Code)
	var folderIDList []string
	if sendToDrive {
		folderIDList, err = google.CreateFolder(connect, cfg.Google().QRPath)
		if err != nil {
			cfg.Log().Debug(err)
			return err
		}
	}
	ctx := context.Background()
	for _, user := range users {
		signature := signature.NewSignature(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle), cfg.Key().Private)
		qr := NewQR(user, cfg, signature)
		hash := signature.Hashing()
		user.SetDataHash(hash)
		path := ""
		path, user.DigitalCertificate, user.Signature, err = qr.GenerateQR()
		if err != nil {
			cfg.Log().Debug(err)
			return err
		}

		chLink := make(chan string)

		if sendToDrive {
			go running.UntilSuccess(ctx, log, "test", func(ctx context.Context) (bool, error) {
				link, err := google.Update(path, connect, folderIDList, cfg)
				chLink <- link
				var success bool
				if err == nil {
					success = true
				}
				return success, err
			}, time.Millisecond*100, time.Millisecond*150)

			user.DigitalCertificate = <-chLink

		}
	}

	SetRes(users, cfg.Table().Result)

	return nil
}
