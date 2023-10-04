package handlers

import (
	"fmt"
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
		Log(r).WithError(err).Error("failed to generate template")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	container := PdfCreator(r).CheckContainerState(containerID)
	if container == nil {
		Log(r).WithError(err).Error("user not found")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	Log(r).Debug("container on  handler: ", container)

	ape.Render(w, newUserWithImgResponse(container.Certificates, container.ID, container.Status))
}

func newUserWithImgResponse(users []*helpers.Certificate, id int, status string) resources.ContainerResponse {
	usersData := make([]resources.User, 0)
	for _, user := range users {
		resp := resources.User{
			Key: resources.Key{
				ID:   fmt.Sprintf("%x", user.ID),
				Type: resources.PARSE_USERS,
			},
			Attributes: resources.UserBlob{ //todo make better
				Participant:        user.Participant,
				Date:               user.Date,
				CourseTitle:        user.CourseTitle,
				CertificateImg:     user.ImageCertificate,
				DigitalCertificate: user.DigitalCertificate,
				Certificate:        user.Certificate,
				Points:             user.Points,
				Note:               user.Note,
				Signature:          user.Signature,
				Msg:                user.Msg,
			},
		}
		usersData = append(usersData, resp)
	}

	return resources.ContainerResponse{
		Data: resources.Container{
			Attributes: resources.ContainerAttributes{
				ContainerId:  fmt.Sprintf("%d", id),
				Certificates: usersData,
				Status:       status,
			},
		},
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
