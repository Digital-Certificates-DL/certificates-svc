package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/ipfs"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/requests"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func UploadFileToIpfs(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUploadFileToIPFS(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse data")
		ape.Render(w, problems.InternalError())
		return
	}

	//recipientAdd := common.HexToAddress(req.Data.Address)
	connector := ipfs.NewConnector(helpers.Config(r))
	img, err := connector.PrepareImage(req.Data.Img)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to send image to ipfs")
		ape.Render(w, problems.InternalError())
		return
	}
	imgHash, err := connector.Upload(img)
	jsonHash, err := connector.PrepareJSON(req.Data.Name, req.Data.Description, imgHash)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to prepare json")
		ape.Render(w, problems.InternalError())
		return
	}

	preparedURI, err := connector.Upload(jsonHash)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to upload")
		ape.Render(w, problems.InternalError())
		return
	}

	ape.Render(w, newIpfsUploadResponse(preparedURI))
}

func newIpfsUploadResponse(uri string) resources.IpfsFile {

	return resources.IpfsFile{
		Key: resources.Key{
			Type: resources.IPFS,
		},
		Attributes: resources.IpfsFileAttributes{
			Url: uri,
		},
	}

}
