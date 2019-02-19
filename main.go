package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/itsubaki/csv2ss/internal/googless"
	"github.com/urfave/cli"
	sheets "google.golang.org/api/sheets/v4"
)

func main() {
	if err := New().Run(os.Args); err != nil {
		panic(err)
	}
}

func New() *cli.App {
	app := cli.NewApp()

	app.Name = "csv2ss"
	app.Usage = "csv to google spreadsheets"
	app.Action = Action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "spreadsheetname, ssn",
		},
		cli.StringFlag{
			Name:  "sheetname, sn",
			Value: "シート1",
		},
	}

	return app
}

func Action(c *cli.Context) {
	values, err := Read()
	if err != nil {
		fmt.Printf("read: %v\n", err)
		return
	}

	ssname := c.String("spreadsheetname")
	sname := c.String("sheetname")

	if len(ssname) < 1 {
		ssname = uuid.Must(uuid.NewRandom()).String()
	}

	ss, res, err := Write(ssname, sname, values)
	if err != nil {
		fmt.Printf("write: %v\n", err)
		return
	}

	fmt.Println(ss)
	fmt.Println(res)
}

func Read() (*sheets.ValueRange, error) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("stdin: %v", err)
	}
	csv := strings.Split(string(stdin), "\n")

	tmp := [][]string{}
	for _, line := range csv {
		tmp = append(tmp, strings.Split(line, ","))
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

func Write(ssname, sname string, values *sheets.ValueRange) (*sheets.Spreadsheet, *sheets.UpdateValuesResponse, error) {
	gss, derr := googless.Default()
	if derr != nil {
		return nil, nil, fmt.Errorf("new spreadsheets client: %v", derr)
	}

	ss, nerr := gss.NewSpreadSheets(ssname)
	if nerr != nil {
		return nil, nil, fmt.Errorf("new spreadsheets: %v", nerr)
	}

	if sname != "シート1" {
		if _, err := gss.NewSheet(ss.SpreadsheetId, sname); err != nil {
			return nil, nil, fmt.Errorf("new sheet=%s: %v", sname, err)
		}
	}

	res, uerr := gss.Update(ss.SpreadsheetId, sname, values)
	if uerr != nil {
		return ss, nil, fmt.Errorf("update sheet1: %v", uerr)
	}

	return ss, res, nil
}
