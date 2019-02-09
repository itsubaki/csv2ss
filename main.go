package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/itsubaki/csv2ss/internal/googless"
	sheets "google.golang.org/api/sheets/v4"
)

func main() {
	values, err := Read()
	if err != nil {
		fmt.Printf("read: %v\n", err)
		return
	}

	ss, res, err := Write(values)
	if err != nil {
		fmt.Printf("write: %v\n", err)
		return
	}

	fmt.Println(ss)
	fmt.Println(res)
	return
}

func Read() (*sheets.ValueRange, error) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("stdin: %v", err)
	}
	csv := strings.Split(string(stdin), "\n")

	tmp := [][]string{}
	for _, line := range csv {
		tmp = append(tmp, strings.Split(line, ", "))
	}

	// string -> interface{}
	val := make([][]interface{}, len(tmp))
	for i := range tmp {
		val[i] = make([]interface{}, len(tmp[i]))
		for j := range tmp[i] {
			val[i][j] = tmp[i][j]
		}
	}

	return &sheets.ValueRange{
		Values: val,
	}, nil
}

func Write(values *sheets.ValueRange) (*sheets.Spreadsheet, *sheets.UpdateValuesResponse, error) {
	gss, derr := googless.Default()
	if derr != nil {
		return nil, nil, fmt.Errorf("new spreadsheets client: %v", derr)
	}

	id := uuid.Must(uuid.NewRandom())
	ss, nerr := gss.NewSpreadSheets(id.String())
	if nerr != nil {
		return nil, nil, fmt.Errorf("new spreadsheets: %v", nerr)
	}

	res, uerr := gss.Update(ss.SpreadsheetId, "シート1", values)
	if uerr != nil {
		return ss, nil, fmt.Errorf("update sheet1: %v", uerr)
	}

	return ss, res, nil
}
