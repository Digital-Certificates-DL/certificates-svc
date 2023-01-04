package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"helper/internal/data"
	"log"
)

func SetRes(users []*data.User, resultFile string) {

	f := excelize.NewFile()

	defer f.Close()
	sheepList := f.GetSheetList()
	for id, user := range users {
		if id == 0 {
			err := f.SetCellValue(sheepList[0], fmt.Sprintf("A%d", id+1), "Date")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("B%d", id+1), "Participant")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("C%d", id+1), "Course Title")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("D%d", id+1), "Serial Number")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("E%d", id+1), "Note")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("F%d", id+1), "Certificate")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("G%d", id+1), "Data Hash")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("H%d", id+1), "Transaction Hash")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("I%d", id+1), "Signature")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("J%d", id+1), "Digital Certificate")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
		}

		err := f.SetCellValue(sheepList[0], fmt.Sprintf("A%d", id+2), user.Date)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("B%d", id+2), user.Participant)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("C%d", id+2), user.CourseTitle)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("D%d", id+2), user.SerialNumber)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("E%d", id+2), user.Note)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("F%d", id+2), user.Certificate)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("G%d", id+2), user.DataHash)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("H%d", id+2), user.TxHash)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("I%d", id+2), user.Signature)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}
		err = f.SetCellValue(sheepList[0], fmt.Sprintf("J%d", id+2), user.CertificatePath)
		if err != nil {
			log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
			continue
		}

	}
	if resultFile == "" {
		err := f.SaveAs("result.xlsx")
		if err != nil {
			log.Println(err)
			return
		}
	}

	err := f.SaveAs(resultFile)
	if err != nil {
		log.Println(err)
		return
	}

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
