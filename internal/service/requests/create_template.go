package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

type GenerateTemplate struct {
	Data resources.Template
}

func NewGenerateTemplate(r *http.Request) (pdf.PDF, []byte, error) {
	response := GenerateTemplate{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return pdf.PDF{}, nil, errors.Wrap(err, "failed to decode data")
	}
	pdfTemplate := pdf.PDF{}
	err = json.Unmarshal(response.Data.Attributes.Template, &pdfTemplate)
	if err != nil {
		return pdf.PDF{}, nil, errors.Wrap(err, "failed to decode data")
	}
	return pdfTemplate, response.Data.Attributes.BackgroundImg, err
}
