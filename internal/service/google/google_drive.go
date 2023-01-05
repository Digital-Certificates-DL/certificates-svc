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

func Update(name string, client *http.Client, folderIDList []string) (string, error) {

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Println(err)
		return "", err
	}
	myQR, err := os.Open("./qr/" + name)
	if err != nil {
		log.Println(err)
		return "", err
	}

	myFile := drive.File{Name: name, Parents: folderIDList, MimeType: "image/svg+xml"}

	file, err := srv.Files.Create(&myFile).Fields().SupportsAllDrives(true).Media(myQR).Do()
	if err != nil {
		log.Println("Couldn't create file ", err)
		return "", err
	}

	return createLink(file.Id), nil

}
func createLink(id string) string {
	return fmt.Sprintf(template, id)

}

func CreateFolder(client *http.Client, folderPath string) ([]string, error) {
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	createFolder, err := srv.Files.Create(&drive.File{Name: folderPath + " " + time.Now().String(), MimeType: "application/vnd.google-apps.folder"}).Do()
	if err != nil {
		log.Fatalf("Unable to create folder: %v", err)
	}

	var folderIDList []string
	folderIDList = append(folderIDList, createFolder.Id)
	return folderIDList, nil
}

func GetFiles(client *http.Client) ([]*drive.File, error) {

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {

		return nil, err
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {

		return nil, err
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
