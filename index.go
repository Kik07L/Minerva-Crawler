package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

var searchEngines = []string{
	"https://www.google.com/search?q=", // Google
	"https://www.bing.com/search?q=",   // Bing
	"https://duckduckgo.com/?q=",       // DuckDuckGo
	"https://www.ask.com/web?q=",       // Ask
	"https://www.ecosia.org/search?q=", // Ecosia
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func searchOnEngine(engine, query string, wg *sync.WaitGroup, results chan<- string, engineName string, summary chan<- string) {
	defer wg.Done()
	url := engine + query
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while fetching %s: %v", engine, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response from %s: %v", engine, err)
		return
	}

	// Extraire les liens
	links := extractLinks(string(body))
	if len(links) > 0 {
		// Ajouter un titre pour ce moteur
		results <- fmt.Sprintf("\n%s%s:%s", colorBlue, engineName, colorReset)
		for _, link := range links {
			// Nettoyer le lien avant de l'afficher
			cleanedLink := cleanURL(link)
			// Afficher le lien sous forme cliquable (si le terminal le permet)
			results <- fmt.Sprintf("%s- %s%s", colorGreen, clickableLink(getDomainName(cleanedLink), cleanedLink), colorReset)
			// Envoyer le lien au canal de résumé
			summary <- cleanedLink
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

	// Supprimer les paramètres inutiles après "&sa=" ou "?"
	if idx := strings.Index(rawURL, "&sa="); idx != -1 {
		rawURL = rawURL[:idx]
	}
	if idx := strings.Index(rawURL, "?"); idx != -1 {
		rawURL = rawURL[:idx]
	}

	// S'assurer que les URL sont propres en supprimant les caractères indésirables
	return rawURL
}

// Fonction qui génère un lien cliquable dans certains terminaux (Mac/Linux)
func clickableLink(displayName, link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, displayName)
}

func extractKeywords(link string) []string {
	// Extraire les mots-clés de l'URL
	re := regexp.MustCompile(`[a-zA-Z]+`)
	return re.FindAllString(link, -1)
}

func displayAnimation() {
	asciiArt := []string{
		"        :::   :::   ::::::::::: ::::    ::: :::::::::: :::::::::  :::     :::     :::  ",
		"      :+:+: :+:+:      :+:     :+:+:   :+: :+:        :+:    :+: :+:     :+:   :+: :+: ",
		"    +:+ +:+:+ +:+     +:+     :+:+:+  +:+ +:+        +:+    +:+ +:+     +:+  +:+   +:+ ",
		"   +#+  +:+  +#+     +#+     +#+ +:+ +#+ +#++:++#   +#++:++#:  +#+     +:+ +#++:++#++: ",
		"  +#+       +#+     +#+     +#+  +#+#+# +#+        +#+    +#+  +#+   +#+  +#+     +#+  ",
		" #+#       #+#     #+#     #+#   #+#+# #+#        #+#    #+#   #+#+#+#   #+#     #+#   ",
		"###       ### ########### ###    #### ########## ###    ###     ###     ###     ###    ",
	}

	colors := []string{colorRed, colorGreen, colorYellow, colorBlue, colorPurple, colorCyan, colorWhite}

	for i := 0; i < 3; i++ { // Repeat the animation 3 times
		for _, line := range asciiArt {
			color := colors[i%len(colors)]
			fmt.Println(color + line + colorReset)
			time.Sleep(50 * time.Millisecond) // Adjust the delay to make it faster
		}
		clearScreen()
	}

	// Print the logo one last time to leave it displayed
	for _, line := range asciiArt {
		fmt.Println(line)
	}
}

func main() {
	clearScreen()
	displayAnimation()

	var query string
	fmt.Print("Entrez le mot-clé à rechercher: ")
	fmt.Scanln(&query)

	var wg sync.WaitGroup
	results := make(chan string, 100)
	summary := make(chan string, 100)

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
		case 3:
			engineName = "Ask"
		case 4:
			engineName = "Ecosia"
		}
		go searchOnEngine(engine, query, &wg, results, engineName, summary)
	}

	// Attendre que toutes les recherches soient terminées et fermer les canaux
	go func() {
		wg.Wait()
		close(results)
		close(summary)
	}()

	// Affichage des résultats
	fmt.Println("\nLiens pertinents trouvés :")
	for result := range results {
		fmt.Println(result)
	}

	// Générer une supposition de réponse globale
	keywordCount := make(map[string]int)
	for link := range summary {
		keywords := extractKeywords(link)
		for _, keyword := range keywords {
			keywordCount[keyword]++
		}
	}

	// Trouver le mot-clé le plus fréquent
	var mostFrequentKeyword string
	var maxCount int
	for keyword, count := range keywordCount {
		if count > maxCount {
			mostFrequentKeyword = keyword
			maxCount = count
		}
	}

	// Afficher la supposition de réponse globale
	if mostFrequentKeyword != "" {
		fmt.Printf("\n%sSupposition de réponse globale :%s %s%s%s\n", colorYellow, colorReset, colorGreen, mostFrequentKeyword, colorReset)
	} else {
		fmt.Printf("\n%sAucune supposition de réponse globale trouvée.%s\n", colorRed, colorReset)
	}
}
