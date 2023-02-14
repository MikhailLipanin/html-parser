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

	var sheet storage.Storage

	sheet, err := google_sheets.New()
	if err != nil {
		log.Fatalf("Error with Google Sheets: %s", err.Error())
	}

	dataInStorage := sheet.ReadAllContent()

	log.Println("Data in HTML table:")
	// Put all Error types into Map (to effectively delete items from it)
	var mapStorage = make(map[parsing.ErrorType]bool)
	for _, el := range dataInStorage {
		fmt.Println(el)
		mapStorage[el] = true
	}

	// Modify the Storage data
	for _, errorType := range data {
		delete(mapStorage, errorType)
		if !sheet.IsPresent(errorType.Id) {
			log.Printf("Adding Error type with id %s\n", errorType.Id)
			err := sheet.AddValById(errorType.Id, errorType.Message)
			if err != nil {
				log.Fatalf("Error with Google Sheets: %s", err.Error())
			}
		} else if sheet.GetValById(errorType.Id) != errorType.Message {
			log.Printf("Updating Error type with id %s\n", errorType.Id)
			delete(mapStorage, parsing.ErrorType{
				Id:      errorType.Id,
				Message: sheet.GetValById(errorType.Id),
			})
			err := sheet.UpdateValById(errorType.Id, errorType.Message)
			if err != nil {
				log.Fatalf("Error with Google Sheets: %s", err.Error())
			}
		}
	}

	// delete the Items, that are not occurred in HTML
	for key, _ := range mapStorage {
		log.Printf("Deleting Error type with id %s\n", key)
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
