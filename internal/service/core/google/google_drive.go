package google

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"io"
	"net/http"
	"strings"
	"time"
)

const template = "https://drive.google.com/file/d/%s/view"

func (g *Google) Update(name string, encodedFile []byte, mimeType string) (string, error) {
	myFile := drive.File{Name: name, Parents: g.folderIDList, MimeType: mimeType}
	file, err := g.driveSrv.Files.Create(&myFile).Media(bytes.NewReader(encodedFile)).Do()
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload file to drive")
	}
	return g.createLink(file.Id), nil

}
func (g *Google) createLink(id string) string {
	return fmt.Sprintf(template, id)

}

func (g *Google) CreateFolder(folderPath string) error {
	createFolder, err := g.driveSrv.Files.Create(&drive.File{Name: folderPath + " " + time.Now().String(), MimeType: "application/vnd.google-apps.folder"}).Do()
	if err != nil {
		return errors.Wrap(err, "Unable to create folder")
	}
	var folderIDList []string
	folderIDList = append(folderIDList, createFolder.Id)
	g.folderIDList = folderIDList
	return nil
}

func (g *Google) GetFiles() ([]*drive.File, error) {

	r, err := g.driveSrv.Files.List().PageSize(10).
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

func (g *Google) GetFile(url string) (*drive.File, error) {
	fileId := g.googleParseURl(url)
	f, err := g.driveSrv.Files.Get(fileId).Do()
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred ")
	}

	return f, nil
}

func (g *Google) Download(url string) ([]byte, error) {
	//f, err := g.driveSrv.Files.Export(g.googleParseURl(url), "application/pdf").Download()

	// Get the file metadata to retrieve the file name and download URL
	file, err := g.driveSrv.Files.Get(g.googleParseURl(url)).Fields("name, webContentLink").Do()
	if err != nil {
		fmt.Println("Unable to retrieve file metadata: ", err)
		return nil, err
	}

	// Download the file content
	resp, err := http.Get(file.WebContentLink)
	if err != nil {
		fmt.Println("Unable to download file: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	out := new(bytes.Buffer)
	if err != nil {
		fmt.Println("Unable to create file: ", err)
		return nil, err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Unable to save file: ", err)
		return nil, err
	}

	return out.Bytes(), nil
}

func (g *Google) googleParseURl(url string) string {
	id := strings.Replace(url, "https://drive.google.com/file/d/", "", 1)
	id = strings.Replace(id, "/view", "", 1)
	return id
}
