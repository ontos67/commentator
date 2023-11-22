// Сервер комментов.
package main

import (
	"log"
	"net/http"
	"os"

	"commentator/pkg/api"
	"commentator/pkg/profanity"
	storage "commentator/pkg/storage/pstg"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Запуск службы...")

	db, err := storage.New()
	defer close(db.CChan)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)
	profanity.ProfanityCheckService(db)

	// запуск сервера
	log.Println("Запуск сервера. Порт: 999...")
	err = http.ListenAndServe(":999", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
