// Сервер комментов.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"commentator/pkg/api"
	"commentator/pkg/profanity"
	storage "commentator/pkg/storage/pstg"
)

type Config struct {
	Port  string `json:"port"`
	DBadr string `json:"dbadr"`
}

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Запуск службы...")
	var conf Config
	b, err := os.ReadFile("./cmd/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.New(conf.DBadr)
	defer close(db.CChan)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	profanity.ProfanityCheckService(db)

	// запуск сервера
	log.Println("Запуск сервера. Порт", conf.Port)
	err = http.ListenAndServe(conf.Port, api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
