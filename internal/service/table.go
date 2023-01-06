package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"helper/internal/data"
	"log"
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

var usersTag = map[string]string{ //todo will make better
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
			log.Println(err)
			return errs
		}
	}
	err := f.SaveAs(resultFile)
	if err != nil {
		log.Println(err)
		return errs
	}
	return nil
}

func Parse(pathToFile string) ([]*data.User, error) {
	users := make([]*data.User, 0)
	f, err := excelize.OpenFile(pathToFile)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	list := f.GetSheetList()

	rows, err := f.GetRows(list[0])
	if err != nil {

		log.Println(err)
		return nil, err
	}

	for id, row := range rows {
		if id < 1 {
			continue
		}
		userInfo := new(data.User)
		userInfo.Date = row[0]
		userInfo.Participant = row[1]
		userInfo.CourseTitle = row[2]
		users = append(users, userInfo)
	}
	return users, err
}
