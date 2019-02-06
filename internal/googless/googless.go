package googless

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

type Config struct {
	Credentials string
	Token       string
	Scope       string
}

func NewConfig() *Config {
	return &Config{
		Credentials: os.Getenv("GOOGLE_CREDENTIALS_PATH"),
		Token:       os.Getenv("GOOGLE_TOKEN_PATH"),
		Scope:       "https://www.googleapis.com/auth/spreadsheet",
	}
}

type GoogleSS struct {
	Config  *Config
	Context context.Context
	Service *sheets.Service
}

func Default() (*GoogleSS, error) {
	return New(NewConfig())
}

func New(config *Config) (*GoogleSS, error) {
	b, err := ioutil.ReadFile(config.Credentials)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	conf, err := google.ConfigFromJSON(b, config.Scope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	f, err := os.Open(config.Token)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", config.Token, err)
	}

	tok := &oauth2.Token{}
	if derr := json.NewDecoder(f).Decode(tok); derr != nil {
		return nil, fmt.Errorf("decode token: %v", derr)
	}

	ctx := context.Background()
	httpc := conf.Client(ctx, tok)
	svc, err := sheets.New(httpc)
	if err != nil {
		return nil, fmt.Errorf("new sheets client: %v", err)
	}

	return &GoogleSS{
		Config:  config,
		Context: ctx,
		Service: svc,
	}, nil
}

func (ss *GoogleSS) NewSpreadSheets(name string) (*sheets.Spreadsheet, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title:    name,
			Locale:   "ja_JP",
			TimeZone: "Asia/Tokyo",
		},
	}

	res, err := ss.Service.Spreadsheets.Create(spreadsheet).Context(ss.Context).Do()
	if err != nil {
		return nil, fmt.Errorf("new spreadsheets: %v", err)
	}

	return res, nil
}

func (ss *GoogleSS) Update(sheetsID, range_ string, value *sheets.ValueRange) (*sheets.UpdateValuesResponse, error) {
	return ss.Service.Spreadsheets.Values.Update(sheetsID, range_, value).ValueInputOption("RAW").Do()
}
