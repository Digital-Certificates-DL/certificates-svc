package handlers

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func PrepareCertificate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPrepareCertificates(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse data")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	users := req.PrepareUsers()

	googleClient := google.NewGoogleClient(Config(r))
	link, err := googleClient.Connect(Config(r).Google().SecretPath, MasterQ(r).ClientQ(), req.Data.Name)
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

	client, err := MasterQ(r).ClientQ().GetByName(req.Data.Name)
	if err != nil {
		Log(r).WithError(err).Error("failed to get client")
		ape.Render(w, problems.InternalError())
		return
	}
	if client == nil {
		Log(r).Error(errors.Wrap(err, "client is not found"))
		ape.RenderErr(w, problems.NotFound())
		return
	}

	PdfCreator(r).NewContainer(users, googleClient, req.Data.Address, req.Data.Url, client, MasterQ(r), pdf.Generate)
	//todo add data to runner

	ape.Render(w, newUserWithImgResponse(users))
}

func newUserWithImgResponse(users []*helpers.User) resources.UserListResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.USER,
			},
			Attributes: resources.UserAttributes{
				Participant:        user.Participant,
				Date:               user.Date,
				CourseTitle:        user.CourseTitle,
				CertificateImg:     user.ImageCertificate,
				DigitalCertificate: user.DigitalCertificate,
				Certificate:        user.Certificate,
				Points:             user.Points,
				Note:               user.Note,
				Signature:          user.Signature,
			},
		}
		usersData = append(usersData, resp)
	}

	return resources.UserListResponse{
		Data: usersData,
	}

}

func newLinkResponse(link string) resources.LinkResponse {
	data := resources.LinkResponse{
		Data: resources.Link{
			Attributes: resources.LinkAttributes{
				Link: link,
			},
		},
	}
	return data
}
