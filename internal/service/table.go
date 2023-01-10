package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"helper/internal/data"
	"reflect"
	"strings"
)

var titles = map[string]string{
	"A": "Date",
	"B": "Participant",
	"C": "Course Title",
	"D": "Serial Number",
	"E": "Note",
	"F": "Certificate",
	"G": "Data Hash",
	"H": "Transaction Hash",
	"I": "Signature",
	"J": "Digital Certificate",
}

var usersTag = map[string]string{
	"A": "Date",
	"B": "Participant",
	"C": "CourseTitle",
	"D": "SerialNumber",
	"E": "Note",
	"F": "Certificate",
	"G": "DataHash",
	"H": "TxHash",
	"I": "Signature",
	"J": "DigitalCertificate",
}

func SetRes(users []*data.User, resultFile string) []error {
	errs := make([]error, 0)
	f := excelize.NewFile()
	defer f.Close()
	sheepList := f.GetSheetList()
	for id, user := range users {
		if id == 0 {
			for key, val := range titles {
				err := f.SetCellValue(sheepList[0], fmt.Sprintf("%s%d", strings.ToUpper(key), id+1), strings.ToTitle(val))
				if err != nil {
					errs = append(errs, err)
					continue
				}
			}
		}
		t := reflect.ValueOf(*user)
		for key, val := range usersTag {
			err := f.SetCellValue(sheepList[0], fmt.Sprintf("%s%d", strings.ToUpper(key), id+2), t.FieldByName(val).String())
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}
	if resultFile == "" {
		err := f.SaveAs("result.xlsx")
		if err != nil {
			return errs
		}
	}
	err := f.SaveAs(resultFile)
	if err != nil {
		return errs
	}
	return nil
}

func Parse(pathToFile string, log *logan.Entry) ([]*data.User, []error) {
	users := make([]*data.User, 0)
	f, err := excelize.OpenFile(pathToFile)
	defer f.Close()
	errs := make([]error, 0)
	if err != nil {
		return nil, append(errs, errors.Wrap(err, "failed to open file"))
	}

	list := f.GetSheetList()

	rows, err := f.GetRows(list[0])
	if err != nil {
		return nil, append(errs, errors.Wrap(err, "failed to get rows from file"))
	}

	for id, row := range rows {
		if id < 1 {
			continue
		}
		userInfo := new(data.User)

		st := reflect.ValueOf(userInfo)
		st = st.Elem()

		for i, str := range row {

			st.Field(i).SetString(str)
			log.Debug(str, i)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
		users = append(users, userInfo)

	}
	if len(errs) != 0 {
		return nil, append(errs, errors.Wrap(err, "failed to parse xlsx"))
	}
	return users, nil
}
