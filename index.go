package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

var searchEngines = []string{
	"https://www.google.com/search?q=",   // Google
	"https://www.bing.com/search?q=",     // Bing
	"https://duckduckgo.com/?q=",         // DuckDuckGo
}

func searchOnEngine(engine, query string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	url := engine + query
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while fetching %s: %v", engine, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response from %s: %v", engine, err)
		return
	}

	// Extraction de liens à partir du HTML
	links := extractLinks(string(body))
	for _, link := range links {
		results <- link
	}
}

func extractLinks(html string) []string {
	// Utilisation d'une expression régulière pour extraire les liens des résultats
	var links []string
	re := regexp.MustCompile(`https?://[^\s]+`)
	matches := re.FindAllString(html, -1)

	// Retirer les doublons (liens)
	linkMap := make(map[string]struct{})
	for _, link := range matches {
		linkMap[link] = struct{}{}
	}

	for link := range linkMap {
		links = append(links, link)
	}

	return links
}

func main() {
	var query string
	fmt.Print("Entrez le mot-clé à rechercher: ")
	fmt.Scanln(&query)

	var wg sync.WaitGroup
	results := make(chan string, 100)

	// Effectuer des recherches en parallèle sur plusieurs moteurs
	for _, engine := range searchEngines {
		wg.Add(1)
		go searchOnEngine(engine, query, &wg, results)
	}

	// Attendre que toutes les goroutines aient fini
	go func() {
		wg.Wait()
		close(results)
	}()

	// Afficher les résultats
	fmt.Println("\nLiens pertinents trouvés :")
	for link := range results {
		fmt.Println(link)
	}
}
