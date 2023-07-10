package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
	"regexp"
)

type SetSettings struct {
	Data resources.Settings
}

func NewSetSettings(r *http.Request) (SetSettings, error) {
	request := SetSettings{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return SetSettings{}, errors.Wrap(err, "failed to decode data")
	}

	if err := validateSettingsData(request.Data); err != nil {
		return SetSettings{}, errors.Wrap(err, "failed to validate data")
	}

	return request, nil
}

func validateSettingsData(request resources.Settings) error {
	return MergeErrors(validation.Errors{
		"/attributes/name": validation.Validate(request.Attributes.Name,
			validation.Required, validation.Match(regexp.MustCompile("^([A-Za-z])[A-Za-z\\s]+$"))),
	}).Filter()
}
