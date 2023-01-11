package service

import (
	"bufio"
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
var paths []Path

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
		googleClient, err = google.NewGoogleClient(cfg)

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

	for id, user := range users {
		user.ID = id
		if user.DataHash != "" || user.Signature != "" || user.DigitalCertificate != "" || user.Certificate != "" {
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
		paths = append(paths, Path{path, id})
	}
	input := make(chan Path, 10)
	output := make(chan Path, len(paths))
	handlers := make([]Handler, 0)
	wg := new(sync.WaitGroup)
	//out := make([]Path, 0)

	if sendToDrive {

		for i := 0; i < 10; i++ {
			log.Debug("start ", i)
			handlers = append(handlers, NewHandler(input, output, log, fmt.Sprintf("test-%d", i), googleClient))
			wg.Add(1)
			go handlers[i].StartRunner(wg)
		}
		for id, path := range paths {
			log.Debug(id)
			input <- path
			log.Debug("input chanel = ")
		}

	}
	log.Debug("move out go")

	//func() {
	//	for {
	//		select {
	//		case <-output:
	//			out = append(out, <-output)
	//			if len(paths) >= len(out) {
	//				log.Debug("return ")
	//				return

	//			}
	//		default:
	//			log.Debug("Wait  ========= ")
	//		}
	//	}
	//}()
	wg.Wait()
	close(input)
	for range paths {
		path, ok := <-output
		if !ok {
			break
		}
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
