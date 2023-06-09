package google

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
	"google.golang.org/api/sheets/v4"
	"log"
	"reflect"
	"strings"
	"time"
)

func (g *Google) GetTable(readRange, spreadsheetId string) error {

	resp, err := g.sheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {

	} else {
		log.Println(resp.Values)
	}
	return nil
}

func (g *Google) ParseFromWeb(spreadsheetId, readRange string, log *logan.Entry) ([]*helpers.User, []error, bool) {
	errs := make([]error, 0)
	resp, err := g.sheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Info("Unable to retrieve data from sheet: %v", err)
		errs = append(errs, err)
		return nil, errs, strings.Contains(err.Error(), "invalid_grant")

	}

	if len(resp.Values) == 0 {
		return nil, nil, false
	} else {
		log.Info(resp.Values)
	}

	users := make([]*helpers.User, 0)
	if err != nil {
		return nil, append(errs, errors.Wrap(err, "failed to open file")), false
	}

	for id, row := range resp.Values {
		if id < 1 {
			continue
		}
		userInfo := new(helpers.User)

		st := reflect.ValueOf(userInfo)
		st = st.Elem()
		for i, str := range row {
			st.Field(i).SetString(fmt.Sprintf("%v", str))
		}
		users = append(users, userInfo)
	}
	if len(errs) != 0 {
		return nil, append(errs, errors.Wrap(err, "failed to parse xlsx")), false
	}
	return users, nil, false
}

func (g *Google) SetRes(users []*helpers.User, sheetID string) []error {

	errs := make([]error, 0)
	for _, user := range users {

		dataForSend := make([]string, 0)

		dataForSend = append(dataForSend, user.SerialNumber, user.Certificate, user.DataHash, user.TxHash, user.Signature, user.DigitalCertificate)

		err := g.UpdateTable(fmt.Sprint("Sheet1!F", user.ID+2), dataForSend, sheetID)
		time.Sleep(4 * time.Millisecond)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return nil
}

func (g *Google) UpdateTable(position string, value []string, spreadsheetId string) error {
	values := stringToInterface(value)
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)
	_, err := g.sheetSrv.Spreadsheets.Values.Update(spreadsheetId, position, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "failed to update sheet")
	}

	return nil
}

func stringToInterface(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, v := range strs {
		res[i] = v
	}
	return res
}

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
