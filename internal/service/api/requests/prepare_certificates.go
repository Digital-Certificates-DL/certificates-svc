package requests

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
	"strings"
)

type PrepareCertificates struct {
	Data resources.PrepareCertificates
}

func NewPrepareCertificates(r *http.Request) (PrepareCertificates, error) {
	request := PrepareCertificates{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return PrepareCertificates{}, errors.Wrap(err, "failed to decode data")
	}

	request.Data.Attributes.Url = request.parse()
	return request, nil
}

func (p *PrepareCertificates) PrepareCertificates() []*helpers.Certificate {
	result := make([]*helpers.Certificate, 0)
	for _, user := range p.Data.Attributes.CertificatesData {
		img, err := base64toJpg(bytes.NewBuffer(user.CertificateImg).String())
		if err != nil {
			img = nil
		}
		resUser := helpers.Certificate{
			ID:                 int(user.Id),
			Date:               user.Date,
			CourseTitle:        user.CourseTitle,
			TxHash:             user.TxHash,
			Signature:          user.Signature,
			DataHash:           user.DataHash,
			SerialNumber:       user.SerialNumber,
			Note:               user.Note,
			Msg:                user.Msg,
			Participant:        user.Participant,
			Points:             user.Points,
			Certificate:        user.Certificate,
			DigitalCertificate: user.DigitalCertificate,
			ImageCertificate:   img,
		}
		result = append(result, &resUser)
	}
	return result
}

func (p *PrepareCertificates) parse() string {
	id := strings.Replace(p.Data.Attributes.Url, "https://docs.google.com/spreadsheets/d/", "", 1)

	index := strings.Index(id, "/")
	if index != -1 {
		id = id[:index]
	}

	return id
}
