package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type GenerateTable struct {
	tableID string `json:"table_id"`
}

func NewGenerateTable(r *http.Request) (GenerateTable, error) {

	response := GenerateTable{}
	err := json.NewDecoder(r.Body).Decode(response)
	if err != nil {
		return GenerateTable{}, errors.Wrap(err, "failed to decode data")
	}
	return response, err
}
