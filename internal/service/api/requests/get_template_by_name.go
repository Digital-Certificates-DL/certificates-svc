package requests

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

const (
	NamePathParam = "name"
)

type GetTemplateByNameRequest struct {
	User         string `url:"-"`
	TemplateName string `name:"-"`
}

func NewGetTemplateByNameRequest(r *http.Request) (GetTemplateByNameRequest, error) {
	request := GetTemplateByNameRequest{}
	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return request, errors.Wrap(err, "failed to decode data")
	}
	request.User = chi.URLParam(r, UserPathParam)
	request.TemplateName = chi.URLParam(r, NamePathParam)

	return request, nil
}
