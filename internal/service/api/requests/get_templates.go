package requests

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

const (
	UserPathParam = "user"
)

type GetTemplateRequest struct {
	User string `url:"-"`
}

func NewGetTemplateRequest(r *http.Request) (GetTemplateRequest, error) {
	request := GetTemplateRequest{}
	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return request, errors.Wrap(err, "failed to decode query")
	}
	request.User = chi.URLParam(r, UserPathParam)

	return request, nil
}
