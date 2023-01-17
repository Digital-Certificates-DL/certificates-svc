package google

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"os"
	"time"
)

const template = "https://drive.google.com/file/d/%s/view"

func (g *Google) Update(path string) (string, error) {

	myQR, err := os.Open(g.cfg.QRCode().QRPath + path)
	if err != nil {
		return "", errors.Wrap(err, "Failed to open file")
	}

	myFile := drive.File{Name: path, Parents: g.folderIDList, MimeType: "image/svg+xml"}

	file, err := g.srv.Files.Create(&myFile).Media(myQR).Do()
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload file to drive")
	}

	return g.createLink(file.Id), nil

}
func (g *Google) createLink(id string) string {
	return fmt.Sprintf(template, id)

}

func (g *Google) CreateFolder(folderPath string) error {
	createFolder, err := g.srv.Files.Create(&drive.File{Name: folderPath + " " + time.Now().String(), MimeType: "application/vnd.google-apps.folder"}).Do()
	if err != nil {
		return errors.Wrap(err, "Unable to create folder")
	}
	var folderIDList []string
	folderIDList = append(folderIDList, createFolder.Id)
	g.folderIDList = folderIDList
	return nil
}

func (g *Google) GetFiles() ([]*drive.File, error) {

	r, err := g.srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get files")
	}
	g.cfg.Log().Info("Files:")
	if len(r.Files) == 0 {

		return nil, errors.New("No files found.")
	} else {
		for _, i := range r.Files {
			g.cfg.Log().Info("%s (%s)\n", i.Name, i.Id)
		}
	}
	return r.Files, nil
}
