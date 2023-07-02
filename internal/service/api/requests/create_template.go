package requests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"image"
	"image/jpeg"
	"net/http"
	"strings"
)

type GenerateTemplate struct {
	Data resources.Template
}

func NewGenerateTemplate(r *http.Request) (pdf.PDF, []byte, GenerateTemplate, error) {
	request := GenerateTemplate{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode data")
	}
	pdfTemplate := pdf.PDF{}
	err = json.Unmarshal(request.Data.Attributes.Template, &pdfTemplate)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode data")
	}

	str := strings.ReplaceAll(request.Data.Attributes.BackgroundImg, "data:image/jpeg;base64,", "")
	str = strings.ReplaceAll(str, "data:image/png;base64,", "")

	data, err := base64toJpg(str)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode data")

	}
	return pdfTemplate, data, request, err
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
