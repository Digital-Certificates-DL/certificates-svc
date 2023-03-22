package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/pdf"
	"net/http"
)

type GenerateTemplate struct {
	Data pdf.PDF
}

func NewGenerateTemplate(r *http.Request) (GenerateTemplate, error) {
	response := GenerateTemplate{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return GenerateTemplate{}, errors.Wrap(err, "failed to decode data")
	}
	return response, err
}
