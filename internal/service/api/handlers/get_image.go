package handlers

import (
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
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
		Log(r).Error(errors.Wrap(err, "failed to parse request "))
		ape.Render(w, problems.BadRequest(err))
		return
	}

	certificates := request.PrepareUsers()

	client := google.NewGoogleClient(Config(r))

	link, err := client.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), request.Data.Attributes.Name)
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

	if err != nil {
		Log(r).WithError(err).Error("failed to connect")
		ape.Render(w, problems.InternalError())
		return
	}

	for _, certificate := range certificates {
		if certificate.Certificate != "" {
			file, err := client.Download(certificate.Certificate)
			if err != nil {
				Log(r).Error("failed to ", err)
				ape.Render(w, problems.BadRequest(err))
				return
			}
			img, err := pdf.NewImageConverter().Convert(file)
			if err != nil {
				Log(r).Error("failed to convert", err)
				ape.Render(w, problems.BadRequest(err))
				return
			}

			certificate.ImageCertificate = img
		}

	}

	ape.Render(w, newUserResponse(certificates))
}
