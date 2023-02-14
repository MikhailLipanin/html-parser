package main

import (
	"fmt"
	"github.com/MikhailLipanin/html-parser/pkg/parsing"
	"github.com/MikhailLipanin/html-parser/pkg/storage"
	"github.com/MikhailLipanin/html-parser/pkg/storage/google_sheets"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error in reading configs: %s", err.Error())
	}
	data := parsing.Parse()
	fmt.Println(data)

	var sheet storage.Storage

	sheet, err := google_sheets.New()
	if err != nil {
		log.Fatalf("Error with Google Sheets: %s", err.Error())
	}

	dataInStorage := sheet.ReadAllContent()

	var mapStorage = make(map[parsing.ErrorType]bool)
	for _, el := range dataInStorage {
		mapStorage[el] = true
	}

	for _, errorType := range data {
		delete(mapStorage, errorType)
		if !sheet.IsPresent(errorType.Id) {
			err := sheet.AddValById(errorType.Id, errorType.Message)
			if err != nil {
				log.Fatalf("Error with Google Sheets: %s", err.Error())
			}
		} else if sheet.GetValById(errorType.Id) != errorType.Message {
			err := sheet.UpdateValById(errorType.Id, errorType.Message)
			if err != nil {
				log.Fatalf("Error with Google Sheets: %s", err.Error())
			}
		}
	}

	for key, _ := range mapStorage {
		err := sheet.DeleteById(key.Id)
		if err != nil {
			log.Fatalf("Error with Google Sheets: %s", err.Error())
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
