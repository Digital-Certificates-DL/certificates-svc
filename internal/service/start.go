package service

import (
	"context"
	"gitlab.com/distributed_lab/running"
	"helper/internal/config"
	"helper/internal/service/google"
	"helper/internal/service/signature"
	"os"
	"time"
)

func Start(cfg config.Config) error {
	log := cfg.Log()
	log.Level(4).Info("Start")

	os.MkdirAll(cfg.QRCode().QRPath, os.ModePerm)

	users, err := Parse(cfg.Table().Input)
	if err != nil {
		cfg.Log().Debug(err)
		return err
	}
	var sendToDrive bool
	sendToDrive = cfg.Google().Enable
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
		signature.Hashing(user)
		path, err := GenerateQR(user, cfg)
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
