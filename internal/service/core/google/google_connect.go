package google

import (
	"bytes"
	"context"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"os"
)

type Google struct {
	client       *http.Client
	folderIDList []string
	cfg          config.Config
	prefixPath   string

	sheetID  string
	log      *logan.Entry
	driveSrv *drive.Service
	sheetSrv *sheets.Service
}

func NewGoogleClient(cfg config.Config) *Google {
	g := Google{
		cfg: cfg,
		log: cfg.Log(),
	}
	return &g
}

func NewGoogleClientTest(prefixPath string) *Google {
	return &Google{
		prefixPath: prefixPath,
	}
}

func (g *Google) getClient(config *oauth2.Config, clientQ data.ClientQ, name string) (*http.Client, string, error) {
	client, err := clientQ.FilterByName(name).Get()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get client")
	}
	if client == nil {
		return nil, "", errors.New("user not found")
	}

	tok := &oauth2.Token{}
	if len(client.Token) == 0 {
		tok, link, err := g.getTokenFromWeb(config, client.Code)
		if err != nil {
			return nil, "", errors.Wrap(err, "you have to update config ")
		}
		if tok == nil {
			return nil, link, nil
		}

		bf := new(bytes.Buffer)
		if err = json.NewEncoder(bf).Encode(tok); err != nil {
			return nil, "", errors.Wrap(err, "failed to decode token")
		}

		client.Token = bf.Bytes()
		if err = clientQ.FilterByID(client.ID).Update(client); err != nil {
			return nil, "", errors.Wrap(err, "failed to update")
		}
	} else {
		g.log.Debug("token already exist")
		if err = json.Unmarshal(client.Token, tok); err != nil {
			return nil, "", errors.Wrap(err, "failed to encode token")
		}
	}

	return config.Client(context.Background(), tok), "", nil
}

// Request a token from the web, then returns the retrieved token.
func (g *Google) getTokenFromWeb(config *oauth2.Config, code string) (*oauth2.Token, string, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if code == "" {
		return nil, authURL, nil
	}

	g.log.Debug(code)
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}
	return tok, "", nil
}

func (g *Google) Connect(path string, clientQ data.ClientQ, name string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "unable to read client secret file")
	}
	link := ""
	googleConfig, err := google.ConfigFromJSON(b, drive.DriveScope, sheets.SpreadsheetsScope)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse client secret file to config")
	}

	g.client, link, err = g.getClient(googleConfig, clientQ, name)
	if err != nil {
		return link, errors.Wrap(err, "unable to get client")
	}
	if g.client == nil {
		return link, errors.New("nil client")
	}
	if len(link) != 0 {
		return link, nil
	}

	g.driveSrv, err = drive.NewService(context.Background(), option.WithHTTPClient(g.client))
	if err != nil {
		return "", errors.Wrap(err, "failed to create new drive service")
	}

	g.sheetSrv, err = sheets.NewService(context.Background(), option.WithHTTPClient(g.client))
	if err != nil {
		return "", errors.Wrap(err, "failed to create new sheet service")
	}

	return "", nil
}
