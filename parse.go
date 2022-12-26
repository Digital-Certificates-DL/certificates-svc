package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func SetRes(users []*user) {

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
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("H%d", id+1), "Signature")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("I%d", id+1), "Digital Certificate")
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
		} else {
			err := f.SetCellValue(sheepList[0], fmt.Sprintf("A%d", id+1), user.Date)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("B%d", id+1), user.Participant)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("C%d", id+1), user.CourseTitle)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("D%d", id+1), user.SerialNumber)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("E%d", id+1), user.Note)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("F%d", id+1), user.Certificate)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("G%d", id+1), user.DataHash)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("H%d", id+1), user.TxHash)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("H%d", id+1), user.Signature)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
			err = f.SetCellValue(sheepList[0], fmt.Sprintf("H%d", id+1), user.DataCertificatePath)
			if err != nil {
				log.Println(fmt.Sprintf("error with %s: %s", user.Participant, err))
				continue
			}
		}

	}

	err := f.SaveAs("result.xlsx")
	if err != nil {
		log.Println(err)
		return
	}
}

func Parse(pathToFile string) ([]*user, error) {
	users := make([]*user, 0)
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
		if id <= 1 {
			continue
		}
		userInfo := new(user)
		userInfo.Date = row[0]
		userInfo.Participant = row[1]
		userInfo.CourseTitle = row[2]
		users = append(users, userInfo)
	}
	return users, err
}
