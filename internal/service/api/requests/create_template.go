package requests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode raw request ")
	}

	pdfTemplate := pdf.PDF{}
	if err := json.Unmarshal(request.Data.Attributes.Template, &pdfTemplate); err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to decode template data")
	}

	rawImg := strings.ReplaceAll(request.Data.Attributes.BackgroundImg, "data:image/jpeg;base64,", "")
	rawImg = strings.ReplaceAll(rawImg, "data:image/png;base64,", "")

	imgBytes, err := base64toJpg(rawImg)
	if err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to convert image")
	}

	if err = validateTemplateData(request.Data); err != nil {
		return pdf.PDF{}, nil, GenerateTemplate{}, errors.Wrap(err, "failed to validate data")

	}
	return pdfTemplate, imgBytes, request, err
}

func base64toJpg(data string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 75}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func validateTemplateData(template resources.Template) error {
	return MergeErrors(validation.Errors{
		"/attributes/is_completed": validation.Validate(template.Attributes.IsCompleted,
			validation.Required),
		"/attributes/template_name": validation.Validate(template.Attributes.TemplateName,
			validation.Required),
		"/attributes/background_img": validation.Validate(template.Attributes.BackgroundImg,
			validation.Required),
	}).Filter()
}
