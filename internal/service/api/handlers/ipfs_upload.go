package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/api/requests"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/ipfs"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
)

func UploadFileToIpfs(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUploadFileToIPFS(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse data")
		ape.Render(w, problems.InternalError())
		return
	}
	connector := ipfs.NewConnector(Config(r).NetworksConfig())
	imgHash, err := connector.Upload(req.Data.Attributes.Img)
	if err != nil {
		Log(r).WithError(err).Error("failed to upload ")
		ape.Render(w, problems.InternalError())
		return
	}
	jsonHash, err := connector.PrepareJSON(req.Data.Attributes.Name, req.Data.Attributes.Description, Config(r).SbtConfig().ExternalURL, imgHash)
	if err != nil {
		Log(r).WithError(err).Error("failed to prepare json")
		ape.Render(w, problems.InternalError())
		return
	}

	preparedURI, err := connector.Upload(jsonHash)
	if err != nil {
		Log(r).WithError(err).Error("failed to upload")
		ape.Render(w, problems.InternalError())
		return
	}

	ape.Render(w, newIpfsUploadResponse(preparedURI))
}

func newIpfsUploadResponse(uri string) resources.IpfsFileResponse {
	return resources.IpfsFileResponse{
		Data: resources.IpfsFile{
			Key: resources.Key{
				Type: resources.IPFS_FILE_UPLOAD,
			},
			Attributes: resources.IpfsFileAttributes{
				Url: uri,
			},
		},
	}
}
