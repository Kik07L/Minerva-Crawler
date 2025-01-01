package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
func getDomainName(link string) string {
	re := regexp.MustCompile(`https?://([^/]+)`)
	match := re.FindStringSubmatch(link)
	if len(match) > 1 {
		return match[1]
	}
	return link
}

// Fonction pour nettoyer l'URL (supprimer les paramètres et échapper les caractères)
func cleanURL(rawURL string) string {
	// Décoder les entités HTML comme &amp; en &
	rawURL = strings.ReplaceAll(rawURL, "&amp;", "&")

	// Décoder l'URL si nécessaire (en utilisant le package net/url)
	decodedURL, err := url.QueryUnescape(rawURL)
	if err == nil {
		rawURL = decodedURL
	}

	// Analyser l'URL pour séparer la base et les paramètres
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Si l'URL est mal formée, retourner telle quelle
	}

	// Si l'URL commence directement par "&" après le domaine, corriger cela
	if strings.HasPrefix(parsedURL.Path, "&") {
		parsedURL.Path = parsedURL.Path[1:]
	}

	// Supprimer les paramètres inutiles (comme ceux de suivi)
	queryParams := parsedURL.Query()
	paramsToRemove := []string{
		"sa", "ved", "usg", "ref", "tracking", "session", "utm_source", "utm_medium", "utm_campaign",
	}

	// On garde les paramètres nécessaires et on filtre ceux qui sont indésirables
	for key := range queryParams {
		if contains(paramsToRemove, key) {
			queryParams.Del(key) // Supprimer les paramètres indésirables
		}
	}

	// Reconstruire l'URL avec les paramètres filtrés
	parsedURL.RawQuery = queryParams.Encode()

	// Retourner l'URL nettoyée
	return parsedURL.String()
}

// Fonction pour vérifier si un paramètre est indésirable
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// Fonction qui génère un lien cliquable dans certains terminaux (Mac/Linux)
func clickableLink(displayName, link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, displayName)
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
