package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type GenerateTable struct {
	Table string `json:"table"`
	Id    string `json:"id"`
}

func NewGenerateTable(r *http.Request) (GenerateTable, error) {
	//https://docs.google.com/spreadsheets/d/1CYqLid0t90bgGm1HPx5j8q-h_RNVVLPVkot9iJZguuo/edit#gid=1988631106
	response := GenerateTable{}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return GenerateTable{}, errors.Wrap(err, "failed to decode data")
	}
	response.parse()
	return response, err
}

func (g *GenerateTable) parse() {
	g.Id = strings.Replace(g.Table, "https://docs.google.com/spreadsheets/d/", "", 1)
	g.Id = strings.Replace(g.Id, "/", "", 1)
}
