package requests

import (
	"encoding/json"
	"github.com/pkg/errors"
	"helper/internal/data"
	"net/http"
)

//type User struct {
//	Id   string `json:"id"`
//	Msg  string `json:"msg"`
//	Sign string `json:"sign"`
//}

func NewUsers(r *http.Request) ([]*data.User, error) {
	response := make([]*data.User, 0)
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return []*data.User{}, errors.Wrap(err, "failed to decode data")
	}
	return response, err
}
