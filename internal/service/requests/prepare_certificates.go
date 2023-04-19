package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
	"strings"
)

type PrepareCertificates struct {
	Data resources.PdfsCreateRequest //todo update model
}

func NewPrepareCertificates(r *http.Request) (PrepareCertificates, error) {
	response := PrepareCertificates{}
	//response.Data.Data = make([]*resources.User, 0)
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return PrepareCertificates{}, errors.Wrap(err, "failed to decode data")
	}
	response.Data.Url = response.parse()
	return response, err
}

func (p PrepareCertificates) PrepareUsers() []*helpers.User {
	result := make([]*helpers.User, 0)
	for _, user := range p.Data.Data {
		resUser := helpers.User{
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

func (g *PrepareCertificates) parse() string {
	id := strings.Replace(g.Data.Url, "https://docs.google.com/spreadsheets/d/", "", 1)
	id = strings.Replace(id, "/", "", 1)
	return id
}
