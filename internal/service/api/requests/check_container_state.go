package requests

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
	"strconv"
)

const (
	containerIDPathParam = "container"
)

type CheckContainerState struct {
	Container string `url:"-"`
}

func NewCheckContainerState(r *http.Request) (int, error) {
	request := CheckContainerState{}
	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return -1, err
	}

	request.Container = chi.URLParam(r, containerIDPathParam)
	containerID, err := strconv.Atoi(request.Container)
	if err != nil {
		return -1, errors.Wrap(err, "failed to convert string to int ")
	}
	return containerID, err
}
