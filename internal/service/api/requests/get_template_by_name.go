package requests

import (
	"github.com/go-chi/chi"
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
	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}
	request.User = chi.URLParam(r, UserPathParam)
	request.TemplateName = chi.URLParam(r, NamePathParam)
	return request, err
}
