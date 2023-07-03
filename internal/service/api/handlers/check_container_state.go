package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func CheckContainerState(w http.ResponseWriter, r *http.Request) {
	containerID, err := requests.NewCheckContainerState(r)
	if err != nil {
		Log(r).Error(errors.Wrap(err, "failed to generate template"))
		ape.Render(w, problems.BadRequest(err))
		return
	}

	container := PdfCreator(r).CheckContainerState(containerID)
	ape.Render(w, newUserWithImgResponse(container.Users))
	return
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
