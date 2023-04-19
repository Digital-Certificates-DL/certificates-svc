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
	res := SetSettings{}
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return SetSettings{}, errors.Wrap(err, "failed to decode data")
	}
	return res, nil
}
