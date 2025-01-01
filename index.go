package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var searchEngines = []string{
	"https://www.google.com/search?q=", // Google
	"https://www.bing.com/search?q=",   // Bing
	"https://duckduckgo.com/?q=",       // DuckDuckGo
}

func searchOnEngine(engine, query string, wg *sync.WaitGroup, results chan<- string, engineName string) {
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

	// Extraire les liens
	links := extractLinks(string(body))
	if len(links) > 0 {
		// Ajouter un titre pour ce moteur
		results <- fmt.Sprintf("\nRésultats pour %s:", engineName)
		for _, link := range links {
			// Nettoyer le lien avant de l'afficher
			cleanedLink := cleanURL(link)
			// Afficher le lien sous forme cliquable (si le terminal le permet)
			results <- clickableLink(getDomainName(cleanedLink), cleanedLink)
		}
	}
}

func extractLinks(html string) []string {
	var links []string
	re := regexp.MustCompile(`https?://[^\s]+`)
	matches := re.FindAllString(html, -1)

	// Utilisation d'un map pour éviter les doublons
	linkMap := make(map[string]struct{})
	for _, link := range matches {
		linkMap[link] = struct{}{}
	}

	// Ajouter les liens sans doublons à la liste finale
	for link := range linkMap {
		links = append(links, link)
	}

	return links
}

// Fonction pour extraire le nom de domaine (ex. https://example.com -> example.com)
func getDomainName(url string) string {
	re := regexp.MustCompile(`https?://([^/]+)`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	}
	return url
}

// Fonction pour nettoyer l'URL (supprimer les paramètres et échapper les caractères)
func cleanURL(url string) string {
	// Décoder les entités HTML comme &amp; en &
	url = strings.ReplaceAll(url, "&amp;", "&")

	// Supprimer les paramètres après ?
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}

	return url
}

// Fonction qui génère un lien cliquable dans certains terminaux (Mac/Linux)
func clickableLink(displayName, url string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, displayName)
}

func main() {
	var query string
	fmt.Print("Entrez le mot-clé à rechercher: ")
	fmt.Scanln(&query)

	var wg sync.WaitGroup
	results := make(chan string, 100)

	// Lancer la recherche pour chaque moteur de recherche
	for idx, engine := range searchEngines {
		wg.Add(1)
		engineName := ""
		switch idx {
		case 0:
			engineName = "Google"
		case 1:
			engineName = "Bing"
		case 2:
			engineName = "DuckDuckGo"
		}
		go searchOnEngine(engine, query, &wg, results, engineName)
	}

	// Attendre que toutes les recherches soient terminées et fermer le canal
	go func() {
		wg.Wait()
		close(results)
	}()

	// Affichage des résultats
	fmt.Println("\nLiens pertinents trouvés :")
	for result := range results {
		fmt.Println(result)
	}
}
