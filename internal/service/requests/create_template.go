package requests

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strings"
)

type GenerateTemplate struct {
	Data resources.Template
}

func NewGenerateTemplate(r *http.Request) (pdf.PDF, []byte, GenerateTemplate, error) {
	response := GenerateTemplate{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode data")
	}
	pdfTemplate := pdf.PDF{}
	err = json.Unmarshal(response.Data.Attributes.Template, &pdfTemplate)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode data")
	}

	str := strings.ReplaceAll(response.Data.Attributes.BackgroundImg, "data:image/jpeg;base64,", "")
	str = strings.ReplaceAll(str, "data:image/png;base64,", "")

	data, err := base64toJpg(str)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode data")

	}
	return pdfTemplate, data, response, err
}

// Given a base64 string of a JPEG, encodes it into an JPEG image test.jpg
func base64toJpg(data string) ([]byte, error) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	//Encode from image format to writer
	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Gets base64 string of an existing JPEG file
func getJPEGbase64(fileName string) string {

	imgFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	//fmt.Println("Base64 string is:", imgBase64Str)
	return imgBase64Str

}
