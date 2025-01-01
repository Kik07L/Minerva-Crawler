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

	links := extractLinks(string(body))
	if len(links) > 0 {
		results <- fmt.Sprintf("\n%s%s:%s", colorBlue, engineName, colorReset)
		for _, link := range links {
			cleanedLink := cleanURL(link)
			results <- fmt.Sprintf("%s- %s%s", colorGreen, clickableLink(getDomainName(cleanedLink), cleanedLink), colorReset)
			summary <- cleanedLink
		}
	}
}

func extractLinks(html string) []string {
	var links []string
	re := regexp.MustCompile(`https?://[^\s]+`)
	matches := re.FindAllString(html, -1)

	linkMap := make(map[string]struct{})
	for _, link := range matches {
		linkMap[link] = struct{}{}
	}

	for link := range linkMap {
		links = append(links, link)
	}

	return links
}

func getDomainName(link string) string {
	re := regexp.MustCompile(`https?://([^/]+)`)
	match := re.FindStringSubmatch(link)
	if len(match) > 1 {
		return match[1]
	}
	return link
}

func cleanURL(rawURL string) string {
	rawURL = strings.ReplaceAll(rawURL, "&amp;", "&")

	if idx := strings.Index(rawURL, "&sa="); idx != -1 {
		rawURL = rawURL[:idx]
	}
	if idx := strings.Index(rawURL, "?"); idx != -1 {
		rawURL = rawURL[:idx]
	}

	return rawURL
}

func clickableLink(displayName, link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, displayName)
}

func extractKeywords(link string) []string {
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

	for i := 0; i < 3; i++ {
		for _, line := range asciiArt {
			color := colors[i%len(colors)]
			fmt.Println(color + line + colorReset)
			time.Sleep(50 * time.Millisecond)
		}
		clearScreen()
	}

	for _, line := range asciiArt {
		fmt.Println(line)
	}
}

func main() {
	clearScreen()
	displayAnimation()

	var query string
	fmt.Print("Enter the keyword to search: ")
	fmt.Scanln(&query)

	var wg sync.WaitGroup
	results := make(chan string, 100)
	summary := make(chan string, 100)

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

	go func() {
		wg.Wait()
		close(results)
		close(summary)
	}()

	fmt.Println("\nRelevant links found:")
	for result := range results {
		fmt.Println(result)
	}

	keywordCount := make(map[string]int)
	for link := range summary {
		keywords := extractKeywords(link)
		for _, keyword := range keywords {
			keywordCount[keyword]++
		}
	}

	var mostFrequentKeyword string
	var maxCount int
	for keyword, count := range keywordCount {
		if count > maxCount {
			mostFrequentKeyword = keyword
			maxCount = count
		}
	}

	if mostFrequentKeyword != "" {
		fmt.Printf("\n%sGlobal response guess:%s %s%s%s\n", colorYellow, colorReset, colorGreen, mostFrequentKeyword, colorReset)
	} else {
		fmt.Printf("\n%sNo global response guess found.%s\n", colorRed, colorReset)
	}
}
