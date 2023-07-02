package handlers

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
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
