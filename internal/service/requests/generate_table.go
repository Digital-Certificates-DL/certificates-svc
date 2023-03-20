package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"helper/internal/data"
	"helper/resources"
	"net/http"
)

type PrepareCertificates struct {
	Pdf resources.PdfsCreateRequest //todo update model
}

func NewPrepareCertificates(r *http.Request) (PrepareCertificates, error) {
	response := PrepareCertificates{}
	response.Pdf.Data = make([]*resources.User, 0)
	err := json.NewDecoder(r.Body).Decode(&response.Pdf)
	if err != nil {
		return PrepareCertificates{}, errors.Wrap(err, "failed to decode data")
	}
	return response, err
}

func (p PrepareCertificates) PrepareUsers() []*data.User {
	result := make([]*data.User, 0)
	for _, user := range p.Pdf.Data {
		resUser := data.User{
			ID:                 int(user.Attributes.ID),
			Date:               user.Attributes.Date,
			CourseTitle:        user.Attributes.CourseTitle,
			TxHash:             user.Attributes.TxHash,
			Signature:          user.Attributes.Signature,
			DataHash:           user.Attributes.DataHash,
			SerialNumber:       user.Attributes.SerialNumber,
			Note:               user.Attributes.Note,
			Msg:                user.Attributes.Msg,
			Participant:        user.Attributes.Participant,
			Points:             user.Attributes.Points,
			Certificate:        user.Attributes.Certificate,
			DigitalCertificate: user.Attributes.DigitalCertificate,
		}
		result = append(result, &resUser)
	}
	return result
}
