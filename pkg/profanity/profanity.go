package profanity

import (
	"bufio"
	storage "commentator/pkg/storage/pstg"
	"log"
	"os"
	"regexp"
	"strings"
)

func ProfanityCheck(c []storage.Comment) ([]int, error) {
	var filth = []int{}
	//читаем словарь мата
	file, err := os.Open("./pkg/profanity/thesaurus.ddw")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var bWords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bWords = append(bWords, scanner.Text())
	}
	pattern := "(?i)" + strings.Join(bWords, "|")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	for _, com := range c {
		if re.MatchString(com.Text) {
			filth = append(filth, com.ID)
		}
	}
	return filth, nil
}
func ProfanityCheckService(db *storage.DB) {
	go func(db *storage.DB) {
		var cList []storage.Comment
		for c := range db.CChan {
			cList = append(cList, c)
			if len(cList) >= 2 {
				idForDelete, err := ProfanityCheck(cList)
				if err != nil {
					log.Println("ошибка проверки на мат:", err)
				}
				if len(idForDelete) > 0 {
					for _, id := range idForDelete {
						err := db.DeleteComment(id)
						if err != nil {
							log.Println("ошибка удаления id:", id, err)
						}
					}
				}
				cList = nil
			}
		}
	}(db)
}
