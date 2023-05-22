package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

type UpdateTokenRequest struct {
	Data resources.Settings
}

func NewUpdateTokenRequest(r *http.Request) (UpdateTokenRequest, error) {
	res := UpdateTokenRequest{}
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return UpdateTokenRequest{}, errors.Wrap(err, "failed to decode data")
	}
	return res, nil
}
