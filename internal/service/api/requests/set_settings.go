package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

type SetSettings struct {
	Data resources.Settings
}

func NewSetSettings(r *http.Request) (SetSettings, error) {
	request := SetSettings{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return SetSettings{}, errors.Wrap(err, "failed to decode data")
	}
	return request, nil
}
