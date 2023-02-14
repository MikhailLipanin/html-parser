package google_sheets

import (
	"context"
	"errors"
	"fmt"
	"github.com/MikhailLipanin/html-parser/pkg/parsing"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"os"
)

type cell struct {
	value string
	row   int
}

type GoogleSheet struct {
	sheet *spreadsheet.Sheet
	data  map[string]*cell
}

func New() (*GoogleSheet, error) {
	scrt, err := os.ReadFile("configs/client-secret.json")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error in reading secrets: %s", err.Error()))
	}

	conf, err := google.JWTConfigFromJSON(scrt, spreadsheet.Scope)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error in authorizing with Google API: %s", err.Error()))
	}

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(viper.GetString("spread-sheet-id"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error in fetching the spreadsheet by the id: %s", err.Error()))
	}

	sheet, err := spreadsheet.SheetByIndex(0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error in getting Google Sheet by id: %s", err.Error()))
	}

	return &GoogleSheet{
		sheet: sheet,
		data:  make(map[string]*cell),
	}, nil
}

func (sh *GoogleSheet) ReadAllContent() []parsing.ErrorType {
	var res []parsing.ErrorType
	for i, _ := range sh.sheet.Rows {
		if len(sh.sheet.Rows[i]) < 2 {
			break
		}
		res = append(res, parsing.ErrorType{
			Id:      sh.sheet.Rows[i][0].Value,
			Message: sh.sheet.Rows[i][1].Value,
		})
	}
	return res
}

func (sh *GoogleSheet) IsPresent(id string) bool {
	_, ok := sh.data[id]
	return ok
}

func (sh *GoogleSheet) GetValById(id string) string {
	val, _ := sh.data[id]
	return val.value
}

func (sh *GoogleSheet) UpdateValById(id, newVal string) error {
	_, ok := sh.data[id]
	if !ok {
		return errors.New(fmt.Sprintf("Error in Updating Sheet Cell at id: %s --- this id is absent", id))
	}
	sh.data[id].value = newVal

	// Update cell content
	sh.sheet.Update(sh.data[id].row, 1, newVal)

	// Make sure call Synchronize to reflect the changes
	err := sh.sheet.Synchronize()
	if err != nil {
		return errors.New(fmt.Sprintf("Error in reflecting the changes of the sheet: %s", err.Error()))
	}

	return nil
}

func (sh *GoogleSheet) AddValById(id, newVal string) error {
	sz := len(sh.data)

	sh.data[id].value = newVal
	sh.data[id].row = sz

	err := sh.UpdateValById(id, newVal)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in Adding Sheet Cell: %s", err.Error()))
	}

	return nil
}

func (sh *GoogleSheet) DeleteById(id string) error {
	row := sh.data[id].row
	err := sh.sheet.DeleteRows(row, row+1)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in Deleting Sheet Cell: %s", err.Error()))
	}
	delete(sh.data, "id")
	return nil
}
