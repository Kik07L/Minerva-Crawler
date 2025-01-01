# Minerva Crawler

```
    :::   :::   ::::::::::: ::::    ::: :::::::::: :::::::::  :::     :::     :::  
      :+:+: :+:+:      :+:     :+:+:   :+: :+:        :+:    :+: :+:     :+:   :+: :+: 
    +:+ +:+:+ +:+     +:+     :+:+:+  +:+ +:+        +:+    +:+ +:+     +:+  +:+   +:+ 
   +#+  +:+  +#+     +#+     +#+ +:+ +#+ +#++:++#   +#++:++#:  +#+     +:+ +#++:++#++: 
  +#+       +#+     +#+     +#+  +#+#+# +#+        +#+    +#+  +#+   +#+  +#+     +#+  
 #+#       #+#     #+#     #+#   #+#+# #+#        #+#    #+#   #+#+#+#   #+#     #+#   
###       ### ########### ###    #### ########## ###    ###     ###     ###     ###    
```

Minerva Crawler is a command-line tool that searches for a given keyword across multiple search engines and displays the results. It supports Google, Bing, DuckDuckGo, Ask, and Ecosia. The tool also provides a global keyword suggestion based on the most frequent keyword found in the search results.

## Features

- Searches across multiple search engines simultaneously
- Displays clickable links in supported terminals

## Prerequisites

- Go 1.16 or higher
- Internet connection

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/Minerva-Crawler.git
    cd Minerva-Crawler
    ```

2. Build the project:
    ```sh
    go build -o minerva-crawler
    ```

## Usage

1. Run the executable:
    ```sh
    ./minerva-crawler
    ```

2. Enter the keyword you want to search for when prompted:
    ```
    Enter the keyword to search: example
    ```

3. View the search results.

## Code Overview

The main functionality is implemented in the `index.go` file. The key functions include:

- `searchOnEngine`: Searches for the keyword on a specific search engine and extracts links from the results.
- `extractLinks`: Extracts links from the HTML content of the search results.
- `cleanURL`: Cleans the extracted URLs by removing unnecessary parameters.
- `clickableLink`: Generates clickable links for supported terminals.
- `extractKeywords`: Extracts keywords from the URLs.
- `displayAnimation`: Displays an ASCII art animation at startup.
- `main`: The entry point of the program that orchestrates the search and displays the results.
