package service

import (
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"helper/internal/config"
	"helper/internal/service/google"
	"helper/internal/service/signature"
	"net/http"
	"os"
	"time"
)

var folderIDList []string
var err error

func Start(cfg config.Config) error {
	log := cfg.Log()
	log.Info("Start")
	os.MkdirAll(cfg.QRCode().QRPath, os.ModePerm)

	users, errs := Parse(cfg.Table().Input, cfg)
	if errs != nil {
		for _, err := range errs {
			cfg.Log().Debug(err)
		}
		return errors.New("failed to parse")
	}
	sendToDrive := cfg.Google().Enable
	var connect *http.Client

	if sendToDrive {
		connect, sendToDrive = google.Connect(cfg.Google().SecretPath, cfg.Google().Code)
		folderIDList, err = google.CreateFolder(connect, cfg.Google().QRPath)
		if err != nil {
			cfg.Log().Debug(err)
			return err
		}
	}
	ctx := context.Background()
	for _, user := range users {
		if user.TxHash != "" || user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" {
			continue
		}
		signature := signature.NewSignature(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle), cfg.Key().Private)
		qr := NewQR(user, cfg, signature)
		hash := signature.Hashing()
		user.SetDataHash(hash)
		var path string
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
			user.DigitalCertificate = <-chLink //todo
		}
	}

	errs = SetRes(users, cfg.Table().Result)
	if errs != nil {
		for _, err := range errs {
			log.Info(err)
		}
		return errors.New("error with creating table")
	}
	return nil
}
