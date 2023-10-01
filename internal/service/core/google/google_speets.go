package google

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"google.golang.org/api/sheets/v4"
	"reflect"
	"time"
)

func (g *Google) ParseFromWeb(spreadsheetId, readRange string) ([]*helpers.Certificate, []error) {
	errs := make([]error, 0)
	resp, err := g.sheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		g.log.Info("Unable to retrieve data from sheet: %v", err)
		errs = append(errs, err)
		return nil, errs
	}

	if len(resp.Values) == 0 {
		return nil, nil
	} else {
		g.log.Info(resp.Values)
	}

	users := make([]*helpers.Certificate, 0)
	for id, row := range resp.Values {
		if id < 1 {
			continue
		}
		userInfo := new(helpers.Certificate)

		st := reflect.ValueOf(userInfo)
		st = st.Elem()
		for i, str := range row {
			st.Field(i).SetString(fmt.Sprintf("%v", str))
		}
		users = append(users, userInfo)
	}
	if len(errs) != 0 {
		return nil, append(errs, errors.Wrap(err, "failed to parse xlsx"))
	}
	return users, nil
}

func (g *Google) SetRes(users []*helpers.Certificate, sheetID string) []error {

	errs := make([]error, 0)
	for _, user := range users {

		dataForSend := make([]string, 0)

		dataForSend = append(dataForSend, user.SerialNumber, user.Certificate, user.DataHash, user.TxHash, user.Signature, user.DigitalCertificate)

		if err := g.UpdateTable(fmt.Sprint("Sheet1!F", user.ID+2), dataForSend, sheetID); err != nil {
			errs = append(errs, err)
			continue
		}
		time.Sleep(4 * time.Millisecond)
	}
	return nil
}

func (g *Google) UpdateTable(position string, value []string, spreadsheetId string) error {
	values := stringToInterface(value)
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)
	if _, err := g.sheetSrv.Spreadsheets.Values.Update(spreadsheetId, position, &vr).ValueInputOption("RAW").Do(); err != nil {
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
