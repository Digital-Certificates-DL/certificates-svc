package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPrepareCertificates(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}
	users := req.PrepareUsers()

	googleClient := google.NewGoogleClient(Config(r))
	link, err := googleClient.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), req.Data.Attributes.Name)
	if err != nil {
		Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	if len(link) != 0 {
		Log(r).WithError(err).Error("failed to authorize")
		ape.RenderErr(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "Invalid token",
			Status: "403",
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}

	client, err := MasterQ(r).ClientQ().GetByName(req.Data.Attributes.Name)
	Log(r).Debug("user ", client)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to get user"))
		ape.Render(w, problems.InternalError())
		return
	}
	if client == nil {
		Log(r).Error(errors.Wrap(err, "user is not found"))
		ape.Render(w, problems.NotFound())
		return
	}

	id := PdfCreator(r).NewContainer(users, googleClient, req.Data.Attributes.Address, req.Data.Attributes.Url, client, MasterQ(r), pdf.Update)

	ape.Render(w, NewContainerResponse(id))
}

func NewContainerResponse(id int) resources.ContainerResponse {
	return resources.ContainerResponse{
		Data: resources.Container{
			Attributes: resources.ContainerAttributes{
				ContainerId: fmt.Sprintf("%d", id),
			},
		},
	}
}
