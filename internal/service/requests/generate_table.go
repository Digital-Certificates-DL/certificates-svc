package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"helper/internal/data"
	"net/http"
)

type PrepareCertificates struct {
	Data []*data.User //todo update model
	//Data    resources.User //todo update model
	Address []byte
	Url     string
}

func NewPrepareCertificates(r *http.Request) (PrepareCertificates, error) {
	response := PrepareCertificates{}
	response.Data = make([]*data.User, 0)
	err := json.NewDecoder(r.Body).Decode(&response.Data)
	if err != nil {
		return PrepareCertificates{}, errors.Wrap(err, "failed to decode data")
	}
	return response, err
}
