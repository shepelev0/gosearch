package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"gosearch/pkg/crawler"
)

func main() {
	// Флаг поиска
	searchWord := flag.String("s", "", "Search word")
	flag.Parse()

	// Инициализация краулера
	c := crawler.New()

	// Сканируем два сайта
	sites := []string{"https://go.dev", "https://golang.org"}
	var allResults []crawler.Result

	for _, site := range sites {
		results, err := c.Crawl(site, 1) // глубина сканирования (0 — вообще не сканировать, 1 — только главная страница, 2 — главная + ссылки с неё
		if err != nil {
			log.Printf("ошибка сканирования %s: %v", site, err)
			continue
		}
		allResults = append(allResults, results...)
	}

	if *searchWord == "" {
		fmt.Println("Сканирование завершено. Результаты:")
		for _, r := range allResults {
			fmt.Printf("- %s\n", r.URL)
		}
	} else {
		fmt.Printf("Поиск слова \"%s\":\n", *searchWord)
		found := false
		for _, r := range allResults {
			if strings.Contains(strings.ToLower(r.Body), strings.ToLower(*searchWord)) {
				fmt.Println(r.URL)
				found = true
			}
		}
		if !found {
			fmt.Println("Ничего не найдено")
		}
	}
}
