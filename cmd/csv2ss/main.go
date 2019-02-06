package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsubaki/csv2ss/internal/googless"
	sheets "google.golang.org/api/sheets/v4"
)

func main() {
	gss, derr := googless.Default()
	if derr != nil {
		fmt.Println(fmt.Errorf("new spreadsheets client: %v", derr))
		return
	}

	id := uuid.Must(uuid.NewRandom())
	ss, nerr := gss.NewSpreadSheets(id.String())
	if nerr != nil {
		fmt.Println(fmt.Errorf("new spreadsheets: %v", nerr))
		return
	}

	value := &sheets.ValueRange{
		Values: [][]interface{}{
			[]interface{}{1, 2, 3, 4, 5},
			[]interface{}{6, 7, 8, 9, 0},
		},
	}

	res, uerr := gss.Update(ss.SpreadsheetId, "シート1", value)
	if uerr != nil {
		fmt.Println(fmt.Errorf("update sheet1: %v", uerr))
		return
	}

	fmt.Println(ss)
	fmt.Println(res)
	return
}
