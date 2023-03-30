package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

type UploadFileToIPFS struct {
	Data resources.IpfsFileUploadRequest //todo replace string to []byte
}

func NewUploadFileToIPFS(r *http.Request) (UploadFileToIPFS, error) {
	response := UploadFileToIPFS{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return UploadFileToIPFS{}, errors.Wrap(err, "failed to decode data")
	}
	return response, err
}
