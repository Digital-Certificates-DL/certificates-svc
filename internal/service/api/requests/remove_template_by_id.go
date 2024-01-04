package requests

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
	"strconv"
)

const (
	TemplateIDPathParam = "template_id"
)

type RemoveTemplateByIDRequest struct {
	TemplateID string `url:"-"`
}

func NewRemoveTemplateByIDRequest(r *http.Request) (*int, error) {
	request := RemoveTemplateByIDRequest{}
	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return nil, errors.Wrap(err, "failed to decode data")
	}
	request.TemplateID = chi.URLParam(r, TemplateIDPathParam)

	id, err := strconv.Atoi(request.TemplateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert template id to int")
	}

	return &id, nil
}
