package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

type IpfsFileUpload struct {
	Data resources.IpfsFileUpload //todo replace string to []byte
}

func NewUploadFileToIPFS(r *http.Request) (IpfsFileUpload, error) {
	request := IpfsFileUpload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return IpfsFileUpload{}, errors.Wrap(err, "failed to decode data")
	}

	if err := validateIpfsData(request.Data); err != nil {
		return IpfsFileUpload{}, errors.Wrap(err, "failed to validate data")
	}

	return request, nil
}

func validateIpfsData(request resources.IpfsFileUpload) error {
	return MergeErrors(validation.Errors{
		"/attributes/description": validation.Validate(request.Attributes.Description,
			validation.Required),
		"/attributes/img": validation.Validate(request.Attributes.Img,
			validation.Required),
	}).Filter()
}
