package service

import (
	"bufio"
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"helper/internal/config"
	"helper/internal/service/google"
	"helper/internal/service/signature"
	"os"
	"strings"
	"sync"
)

var folderIDList []string
var err error
var paths []Paths

func Start(cfg config.Config) error {
	log := cfg.Log()
	log.Info("Start")
	os.MkdirAll(cfg.QRCode().QRPath, os.ModePerm)

	users, errs := Parse(cfg.Table().Input, cfg.Log())
	if errs != nil {
		for _, err := range errs {
			cfg.Log().Debug(err)
		}
		return errors.New("failed to parse")
	}
	sendToDrive := cfg.Google().Enable

	var googleClient *google.Google
	if sendToDrive {
		googleClient = google.NewGoogleClient(cfg)
		err = googleClient.Connect(cfg.Google().SecretPath, cfg.Google().Code)
		if err != nil {
			log.Info("Could you continue to work without google drive? (y)")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if strings.ToLower(text) != "y\n" {
				sendToDrive = true
			}
		}
	}
	if sendToDrive {
		err = googleClient.CreateFolder(cfg.Google().QRPath)
		if err != nil {
			cfg.Log().Debug(err)
			return err
		}
	}

	ctx := context.Background()
	for id, user := range users {
		user.ID = id
		if user.TxHash != "" || user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" {
			log.Debug("skip")
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
		paths = append(paths, Paths{path, id})
	}
	handlers := make([]Handler, 0)
	wg := new(sync.WaitGroup)
	if sendToDrive {
		for id, path := range paths {
			log.Debug(id)
			if id < 10 { //todo move to config
				log.Debug("start ", id)
				handlers = append(handlers, NewHandler(ctx, make(chan Paths), log, fmt.Sprintf("test-%d", id), googleClient))
				handlers[id].SetData(path)
				wg.Add(1)
				go handlers[id].StartRunner(wg)
				continue
			}
			handlers[id%10].SetData(path)
		}

	}
	wg.Wait()
	paths := make([]Paths, 0)
	log.Debug("move out go")
	for _, handler := range handlers {
		path, _ := handler.Result()
		paths = append(paths, path...)
	}

	for _, path := range paths {
		users[path.ID].DigitalCertificate = path.Path
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
