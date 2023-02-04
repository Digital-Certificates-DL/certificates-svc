package service

import (
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"helper/internal/config"
	"helper/internal/handlers"
	"helper/internal/service/signature"
	"os"
)

var err error
var paths []handlers.Path

func Start(cfg config.Config) error {
	log := cfg.Log()
	log.Info("Start")
	os.MkdirAll(cfg.QRCode().QRPath, os.ModePerm)

	users, errs := Parse(cfg.Table().Input, cfg.Log())
	if errs != nil {
		for _, err := range errs {
			log.Error(err)
		}

		err = errors.New("failed to parse")
		log.Error(err)
		return err
	}
	sign, err := signature.NewSignature(cfg.Key().Private)
	if err != nil {
		log.Error(err)
		return err
	}
	for id, user := range users {
		user.ID = id
		if user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" {
			log.Debug("skip")
			continue
		}

		qr := NewQR(user, cfg, sign)
		hash := sign.Hashing(fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle))
		if hash != "" {
			log.Info(user.Participant, " hash = ", hash)
		}
		user.SetDataHash(hash)
		var path string
		path, user.DigitalCertificate, user.Signature, err = qr.GenerateQR()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debug(user.Participant, "local qr path = ", user.DigitalCertificate)
		paths = append(paths, handlers.Path{Path: path, ID: id})
	}

	if sendToDrive := cfg.Google().Enable; sendToDrive {
		users, err = handlers.Drive(cfg, log, paths, users)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	log.Info("creating table")
	errs = SetRes(users, cfg.Table().Result)
	if errs != nil {
		for _, err := range errs {
			log.Error(err)
		}
		err = errors.New("error with creating table")
		log.Error(err)
		return err
	}
	return nil
}
