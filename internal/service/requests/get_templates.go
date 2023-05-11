package requests

import (
	"github.com/go-chi/chi"
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
	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}
	request.User = chi.URLParam(r, UserPathParam)
	return request, err
}
