package service

import (
	"fmt"
	"helper/internal/config"
	"helper/internal/data"
	"helper/internal/service/google"
	"helper/internal/service/signature"
	"net/http"
	"os"
	"sync"
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

	connect, sendToDrive := google.Connect(cfg.Google().SecretPath, cfg.Google().Code)
	var folderIDList []string
	if sendToDrive {
		folderIDList, err = google.CreateFolder(connect, cfg.Google().QRPath)
		if err != nil {
			cfg.Log().Debug(err)
			return err
		}
	}

	wg := new(sync.WaitGroup)
	for _, user := range users {
		wg.Add(1)
		go deciding(user, cfg, wg, connect, folderIDList, sendToDrive)
	}
	wg.Wait()
	SetRes(users, cfg.Table().Result)

	return nil
}

func deciding(user *data.User, cfg config.Config, wg *sync.WaitGroup, client *http.Client, folderIDList []string, sendToDrive bool) {
	defer wg.Done()
	if user.Signature != "" && user.SerialNumber != "" {
		return
	}
	signature.Hashing(user)
	path, err := GenerateQR(user, cfg)
	if err != nil {
		cfg.Log().Debug(err)
		return
	}
	if sendToDrive {
		link, err := google.Update(path, client, folderIDList, cfg)
		if err != nil {
			for {
				time.Sleep(5 * time.Microsecond)
				_, err := google.Update(path, client, folderIDList, cfg)
				if err == nil {
					break
				}
			}

		}
		user.CertificatePath = link
		return
	}
	user.CertificatePath = fmt.Sprintf("%s%s", cfg.QRCode().QRPath, path)
	return
}
