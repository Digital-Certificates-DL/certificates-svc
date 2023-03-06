package google

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"google.golang.org/api/sheets/v4"
	"helper/internal/data"
	"helper/internal/service/helpers"
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

func (g *Google) UpdateTable(position, value, spreadsheetId string) error {
	//	writeRange := "A1" // or "sheet1:A1" if you have a different sheet
	values := []interface{}{value}
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)
	_, err := g.sheetSrv.Spreadsheets.Values.Update(spreadsheetId, position, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "failed to update sheet")
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
	//input := make(chan handlers.Path)
	//output := make(chan handlers.Path)
	//ctx := context.Background()
	//handler := handlers.NewHandler(input, output, g.cfg.Log(), g, 10, ctx)
	//handler.StartSheetRunner()
	//go handler.insertData(users)
	//users = handler.Read(users)

	errs := make([]error, 0)
	for _, user := range users {
		t := reflect.ValueOf(*user)
		counter := 0
		for key, val := range helpers.UsersTag {
			if counter < 3 {
				counter++
			}
			err := g.UpdateTable(fmt.Sprintf("%s%d", strings.ToUpper(key), user.ID+2), t.FieldByName(val).String(), sheetID)
			time.Sleep(1 * time.Millisecond)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}
	return nil
}
