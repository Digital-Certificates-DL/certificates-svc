package google

import (
	"fmt"
	"log"
)

func (g *Google) GetTable(readRange string) error {
	spreadsheetId := "1zpK6lQRucAbtwkGcee7XswMtAJ-nSklfM1iet7rajWc"

	resp, err := g.sheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		log.Println(resp.Values)
	}
	return nil
}

func (g *Google) UpdateTable() error {

	return nil
}
