// Сервер комментов.
package main

import (
	"log"
	"net/http"

	"commentator/pkg/api"
	"commentator/pkg/profanity"
	storage "commentator/pkg/storage/pstg"
)

func main() {
	db, err := storage.New()
	defer close(db.CChan)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)
	profanity.ProfanityCheckService(db)

	// запуск сервера
	err = http.ListenAndServe(":9999", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
