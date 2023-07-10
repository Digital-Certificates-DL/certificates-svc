package ipfs

import (
	"bytes"
	"encoding/json"
	shell "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type ConnectorHandler interface {
	Upload(data []byte) (string, error)
	PrepareJSON(tokenName, tokenDescription, link, imagePath string) ([]byte, error)
	PrepareImagePath(imagePath string) ([]byte, error)
}

type Connector struct {
	cfg *config.NetworksConfig
}

func NewConnector(cfg *config.NetworksConfig) *Connector {
	return &Connector{cfg: cfg}
}

func (i Connector) Upload(data []byte) (string, error) {
	ipfs := shell.NewShellWithClient(i.cfg.IpfsEndpoint, NewClient(i.cfg.IpfsPrId, i.cfg.IpfsPrKey))

	fileHash, err := ipfs.Add(bytes.NewReader(data))
	if err != nil {
		return "", errors.Wrap(err, "failed to upload file to ipfs")
	}
	return fileHash, nil

}

func (i Connector) PrepareJSON(tokenName, tokenDescription, link, imagePath string) ([]byte, error) {
	erc721 := ERC721json{
		Name:        tokenName,
		Description: tokenDescription,
		Image:       i.cfg.IpfsDisplayFileDomen + imagePath,
		ExternalUrl: link,
	}

	erc721JSON, err := json.Marshal(erc721)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal erc721")
	}

	return erc721JSON, nil
}

func (i Connector) PrepareImagePath(imagePath string) ([]byte, error) {
	path, err := filepath.Abs("main.go") //todo  make better
	if err != nil {
		return nil, errors.Wrap(err, "failed to get absolute path")
	}
	path = strings.ReplaceAll(path, "main.go", "")

	infile, err := os.Open(path + imagePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open image")
	}

	img, err := png.Decode(infile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}

	buf := new(bytes.Buffer)

	if err = png.Encode(buf, img); err != nil {
		return nil, errors.Wrap(err, "failed to decode image to []byte")
	}
	return buf.Bytes(), nil
}
