package handlers

import (
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"net/http"
)

func PrepareCertificate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPrepareCertificates(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	certificates := req.PrepareCertificates()

	googleClient := google.NewGoogleClient(Config(r))
	link, err := googleClient.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), req.Data.Attributes.Name)
	if err != nil {
		Log(r).WithError(err).Error("failed to connect")
		ape.RenderErr(w, problems.InternalError())
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

	client, err := MasterQ(r).ClientQ().FilterByName(req.Data.Attributes.Name).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get client")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if client == nil {
		Log(r).WithError(err).Error("client is not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	id := PdfCreator(r).NewContainer(certificates, googleClient, req.Data.Attributes.Address, req.Data.Attributes.Url, client, MasterQ(r), pdf.Generate)

	ape.Render(w, NewContainerResponse(id))
}
