package google

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"time"
)

const template = "https://drive.google.com/file/d/%s/view"

var srv *drive.Service //todo rewrite it

func (g *Google) Update(path string) (string, bool, error) {
	var err error
	if srv == nil { //todo wrap error
		srv, err = drive.NewService(context.Background(), option.WithHTTPClient(g.client))
	}

	if err != nil {
		return "", false, err
	}
	myQR, err := os.Open(g.cfg.QRCode().QRPath + path)
	if err != nil {
		return "", false, err
	}

	myFile := drive.File{Name: path, Parents: g.folderIDList, MimeType: "image/svg+xml"}

	file, err := srv.Files.Create(&myFile).Fields().SupportsAllDrives(true).Media(myQR).Do()
	if err != nil {
		return "", false, err
	}

	return g.createLink(file.Id), true, nil

}
func (g *Google) createLink(id string) string {
	return fmt.Sprintf(template, id)

}

func (g *Google) CreateFolder(folderPath string) error {
	var err error
	if srv == nil {
		srv, err = drive.NewService(context.Background(), option.WithHTTPClient(g.client))
	}

	if err != nil {
		return err //todo will make wrap
	}

	createFolder, err := srv.Files.Create(&drive.File{Name: folderPath + " " + time.Now().String(), MimeType: "application/vnd.google-apps.folder"}).Do()
	if err != nil {
		log.Fatalf("Unable to create folder: %v", err)
	}

	var folderIDList []string
	folderIDList = append(folderIDList, createFolder.Id)
	g.folderIDList = folderIDList
	return nil
}

func (g *Google) GetFiles(client *http.Client) ([]*drive.File, error) {

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {

		return nil, err
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {

		return nil, err //todo will make wrap
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {

		return nil, errors.New("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	return r.Files, nil
}
