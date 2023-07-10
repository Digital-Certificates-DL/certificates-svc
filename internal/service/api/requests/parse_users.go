package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
	"regexp"
	"strings"
)

type ParseUsers struct {
	Data resources.ParseUsers
}

func NewGetUsers(r *http.Request) (ParseUsers, error) {
	request := ParseUsers{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return ParseUsers{}, errors.Wrap(err, "failed to decode data")
	}

	if err := validateParseData(request.Data); err != nil {
		return ParseUsers{}, errors.Wrap(err, "failed to validate data")
	}

	request.Data.Attributes.Url = request.parse()

	return request, nil
}

func (g *ParseUsers) parse() string {
	id := strings.Replace(g.Data.Attributes.Url, "https://docs.google.com/spreadsheets/d/", "", 1)
	id = strings.Replace(id, "/", "", 1)

	return id
}

func validateParseData(request resources.ParseUsers) error {
	return MergeErrors(validation.Errors{
		"/attributes/url": validation.Validate(request.Attributes.Url,
			validation.Required),
		"/attributes/name": validation.Validate(request.Attributes.Name,
			validation.Required, validation.Match(regexp.MustCompile("^([A-Za-z])[A-Za-z\\s]+$"))),
	}).Filter()
}
