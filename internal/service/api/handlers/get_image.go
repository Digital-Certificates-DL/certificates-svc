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

func GetImages(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPrepareCertificates(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse request ")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	certificates := request.PrepareCertificates()
	googleClient := google.NewGoogleClient(Config(r))

	link, err := googleClient.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), request.Data.Attributes.Name)
	if len(link) != 0 {
		Log(r).WithError(err).Debug("failed to authorize")
		ape.RenderErr(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "Invalid token",
			Status: "403",
			Meta:   &map[string]interface{}{"auth_link": link}},
		}...)

		return
	}

	if err != nil {
		Log(r).WithError(err).Debug("failed to connect")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	for _, certificate := range certificates {
		if certificate.Certificate != "" {
			file, err := googleClient.Download(certificate.Certificate)
			if err != nil {
				Log(r).WithError(err).Debug("failed to download  file")
				ape.RenderErr(w, problems.InternalError())
				return
			}
			img, err := pdf.NewImageConverter().Convert(file)
			if err != nil {
				Log(r).WithError(err).Debug("failed to convert")
				ape.RenderErr(w, problems.InternalError())
				return
			}

			certificate.ImageCertificate = img
		}

	}

	ape.Render(w, newUserResponse(certificates))
}
