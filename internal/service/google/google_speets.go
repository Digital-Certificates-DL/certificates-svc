package google

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"helper/internal/data"
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
	//todo add handler
	//input := make(chan handlers.File)
	//output := make(chan handlers.File)
	//ctx := context.Background()
	//handler := handlers.NewHandler(input, output, g.cfg.Log(), g, 10, ctx)
	//handler.StartSheetRunner()
	//go handler.insertData(users)
	//users = handler.Read(users)

	errs := make([]error, 0)
	for _, user := range users {

		dataForSend := make([]string, 0)

		dataForSend = append(dataForSend, user.SerialNumber, user.Note, user.Certificate, user.DataHash, user.TxHash, user.Signature, user.DigitalCertificate)

		err := g.UpdateTable(fmt.Sprint("sheet1!D", user.ID+2), dataForSend, sheetID)
		time.Sleep(1 * time.Millisecond)
		if err != nil {
			errs = append(errs, err)
			continue
		}
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
