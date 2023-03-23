package google

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"google.golang.org/api/sheets/v4"
	"log"
	"reflect"
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

func (g *Google) ParseFromWeb(spreadsheetId, readRange string, log *logan.Entry) ([]*data.User, []error) {
	resp, err := g.sheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		return nil, nil
	} else {
		log.Info(resp.Values)
	}
	errs := make([]error, 0)
	users := make([]*data.User, 0)
	if err != nil {
		return nil, append(errs, errors.Wrap(err, "failed to open file"))
	}

	for id, row := range resp.Values {
		if id < 1 {
			continue
		}
		userInfo := new(data.User)

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

func (g *Google) SetRes(users []*data.User, sheetID string) []error {

	errs := make([]error, 0)
	for _, user := range users {

		dataForSend := make([]string, 0)

		dataForSend = append(dataForSend, user.SerialNumber, user.Note, user.Certificate, user.DataHash, user.TxHash, user.Signature, user.DigitalCertificate)

		err := g.UpdateTable(fmt.Sprint("Sheet1!E", user.ID+2), dataForSend, sheetID)
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
