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
)

func Update(name string, client *http.Client) error {

	myFile := drive.File{Name: name}
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Println(err)
		return err
	}
	myQR, err := os.Open("./" + name)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = srv.Files.Create(&myFile).Media(myQR).Do()
	if err != nil {
		log.Println("Couldn't create file ", err)
		return err
	}
	return nil

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

//
//func Update(name, path string, client *http.Client) error {
//
//	myFile := drive.File{Name: filepath.Base(path)}
//	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	createdFile, err := srv.Files.Create(&myFile).Do()
//	if err != nil {
//		log.Println("Couldn't create file ", err)
//		return err
//	}
//
//	myQR, err := os.Open("./" + name)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	updatedFile := drive.File{Name: name}
//
//	//srv.Files.Update(createdFile.Id, &updatedFile)
//
//	_, err = srv.Files.Create (createdFile.Id, &updatedFile).Media(myQR).Do()
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	return nil
//
//}
