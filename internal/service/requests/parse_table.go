package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/resources"
	"net/http"
	"strings"
)

type GetUsers struct {
	Data resources.UsersGetRequest
}

func NewGetUsers(r *http.Request) (GetUsers, error) {
	//https://docs.google.com/spreadsheets/d/1CYqLid0t90bgGm1HPx5j8q-h_RNVVLPVkot9iJZguuo/edit#gid=1988631106
	response := GetUsers{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return GetUsers{}, errors.Wrap(err, "failed to decode data")
	}
	response.Data.Url = response.parse()
	return response, err
}

func (g *GetUsers) parse() string {
	id := strings.Replace(g.Data.Url, "https://docs.google.com/spreadsheets/d/", "", 1)
	id = strings.Replace(id, "/", "", 1)
	return id
}
