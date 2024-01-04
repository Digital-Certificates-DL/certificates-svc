package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"net/http"
)

func RemoveTemplateByID(w http.ResponseWriter, r *http.Request) {
	templateID, err := requests.NewRemoveTemplateByIDRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request ")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	templateID64 := int64(*templateID)
	err = MasterQ(r).TemplateQ().FilterByID(templateID64).Delete()
	if err != nil {
		Log(r).WithError(err).Error("failed to remove template")
		ape.Render(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
