package service

import (
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"helper/internal/config"
	"helper/internal/handlers"
	"helper/internal/service/signature"
	"os"
)

var folderIDList []string
var err error
var paths []handlers.Path

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

	for id, user := range users {
		user.ID = id
		if len(user.TxHash) > 1 || user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" {
			log.Debug("skip")
			continue
		}
		signature := signature.NewSignature(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle), cfg.Key().Private)
		qr := NewQR(user, cfg, signature)
		hash := signature.Hashing()
		if hash != "" {
			log.Info(user.Participant, "hash = ", hash)
		}
		user.SetDataHash(hash)
		var path string
		path, user.DigitalCertificate, user.Signature, err = qr.GenerateQR()
		if err != nil {
			cfg.Log().Debug(err)
			return err
		}
		log.Info(user.Participant, "local qr path = ", user.DigitalCertificate)
		paths = append(paths, handlers.Path{Path: path, ID: id})
	}

	if sendToDrive := cfg.Google().Enable; sendToDrive {
		users, err = handlers.Drive(cfg, log, paths, users)
		if err != nil {
			return err
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
