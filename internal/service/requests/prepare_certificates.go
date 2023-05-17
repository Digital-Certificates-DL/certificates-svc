package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"net/http"
	"strings"
)

type PrepareCertificates struct {
	Data PdfsCreateRequest //todo update model
}

type PdfsCreateRequest struct {
	Address string `json:"address"`
	Data    []User `json:"users"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

type User struct {
	Certificate        string `json:"certificate"`
	CertificateImg     string `json:"certificateImg"`
	CourseTitle        string `json:"courseTitle"`
	DataHash           string `json:"dataHash"`
	Date               string `json:"date"`
	DigitalCertificate string `json:"digitalCertificate"`
	ID                 int64  `json:"id"`
	Msg                string `json:"msg"`
	Note               string `json:"note"`
	Participant        string `json:"participant"`
	Points             string `json:"points"`
	SerialNumber       string `json:"serialNumber"`
	Signature          string `json:"signature"`
	TxHash             string `json:"txHash"`
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
		//id, err := strconv.ParseInt(user.ID, 16, 64)
		//if err != nil {
		//	return nil
		//}
		img, err := base64toJpg(user.CertificateImg)
		if err != nil {
			img = nil
		}
		resUser := helpers.User{
			ID:                 int(user.ID),
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

func (g *PrepareCertificates) parse() string {
	id := strings.Replace(g.Data.Url, "https://docs.google.com/spreadsheets/d/", "", 1)
	id = strings.Replace(id, "/", "", 1)
	return id
}
