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
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return PrepareCertificates{}, errors.Wrap(err, "failed to decode data")
	}
	request.Data.Attributes.Url = request.parse()
	return request, err
}

func (p *PrepareCertificates) PrepareUsers() []*helpers.User {
	result := make([]*helpers.User, 0)
	for _, user := range p.Data.Attributes.CertificatesData {
		//id, err := strconv.ParseInt(user.ID, 16, 64)
		//if err != nil {
		//	return nil
		//}
		img, err := base64toJpg(bytes.NewBuffer(user.CertificateImg).String())
		if err != nil {
			img = nil
		}
		resUser := helpers.User{
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
	id = strings.Replace(id, "/", "", 1)
	return id
}
