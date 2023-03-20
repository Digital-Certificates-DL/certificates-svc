package table

//
//import (
//	"fmt"
//	"github.com/xuri/excelize/v2"
//	"gitlab.com/distributed_lab/logan/v3"
//	"gitlab.com/distributed_lab/logan/v3/errors"
//	"helper/internal/data"
//	"helper/internal/service/helpers"
//	"reflect"
//	"strings"
//)
//
//func SetRes(users []*data.User, resultFile string) []error {
//
//	errs := make([]error, 0)
//	f := excelize.NewFile()
//	defer f.Close()
//	sheepList := f.GetSheetList()
//	for id, user := range users {
//		if id == 0 {
//			for key, val := range helpers.Titles {
//				err := f.SetCellValue(sheepList[0], fmt.Sprintf("%s%d", strings.ToUpper(key), id+1), strings.ToTitle(val))
//				if err != nil {
//					errs = append(errs, err)
//					continue
//				}
//			}
//		}
//		t := reflect.ValueOf(*user)
//		for key, val := range helpers.UsersTag {
//			err := f.SetCellValue(sheepList[0], fmt.Sprintf("%s%d", strings.ToUpper(key), id+2), t.FieldByName(val).String())
//			if err != nil {
//				errs = append(errs, err)
//				continue
//			}
//		}
//	}
//	if resultFile == "" {
//		err := f.SaveAs("result.xlsx")
//		if err != nil {
//			return errs
//		}
//	}
//	err := f.SaveAs(resultFile)
//	if err != nil {
//		return errs
//	}
//	return nil
//}
//
//func Parse(pathToFile string, log *logan.Entry) ([]*data.User, []error) {
//	users := make([]*data.User, 0)
//	f, err := excelize.OpenFile(pathToFile)
//	defer f.Close()
//	errs := make([]error, 0)
//	if err != nil {
//		return nil, append(errs, errors.Wrap(err, "failed to open file"))
//	}
//
//	list := f.GetSheetList()
//
//	rows, err := f.GetRows(list[0])
//	if err != nil {
//		return nil, append(errs, errors.Wrap(err, "failed to get rows from file"))
//	}
//	for id, row := range rows {
//		if id < 1 {
//			continue
//		}
//		userInfo := new(data.User)
//
//		st := reflect.ValueOf(userInfo)
//		st = st.Elem()
//		for i, str := range row {
//			st.Field(i).SetString(str)
//			log.Debug(str, i)
//			if err != nil {
//				errs = append(errs, err)
//				continue
//			}
//		}
//		users = append(users, userInfo)
//
//	}
//	if len(errs) != 0 {
//		return nil, append(errs, errors.Wrap(err, "failed to parse xlsx"))
//	}
//	return users, nil
//}
