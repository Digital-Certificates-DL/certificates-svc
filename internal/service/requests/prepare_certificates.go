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
	Certificate        string `json:"Certificate"`
	CertificateImg     []byte `json:"CertificateImg"`
	CourseTitle        string `json:"CourseTitle"`
	DataHash           string `json:"DataHash"`
	Date               string `json:"Date"`
	DigitalCertificate string `json:"DigitalCertificate"`
	ID                 int64  `json:"UserID"`
	Msg                string `json:"Msg"`
	Note               string `json:"Note"`
	Participant        string `json:"Participant"`
	Points             string `json:"Points"`
	SerialNumber       string `json:"SerialNumber"`
	Signature          string `json:"Signature"`
	TxHash             string `json:"TxHash"`
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
