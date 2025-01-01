import requests
import urllib.parse
from bs4 import BeautifulSoup

def google_search(query, max_results=5):
    print("[Google Search]")
    url = f"https://www.google.com/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.select("div.tF2Cxc"):
            title = item.select_one("h3").text if item.select_one("h3") else ""
            link = item.select_one("a")["href"] if item.select_one("a") else ""
            description = item.select_one("span.aCOpRe").text if item.select_one("span.aCOpRe") else ""
            if link:
                results.append({"title": title, "link": link, "description": description})
                if len(results) >= max_results:
                    break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur Google: {e}")
        return []

def bing_search(query, max_results=5):
    print("[Bing Search]")
    url = f"https://www.bing.com/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.find_all("li", class_="b_algo"):
            title = item.find("h2").text if item.find("h2") else ""
            link = item.find("a")["href"] if item.find("a") else ""
            description = item.find("p").text if item.find("p") else ""
            if link:
                results.append({"title": title, "link": link, "description": description})
                if len(results) >= max_results:
                    break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur Bing: {e}")
        return []

def duckduckgo_search(query, max_results=5):
    print("[DuckDuckGo Search]")
    url = f"https://duckduckgo.com/html/?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.find_all("a", class_="result__a", href=True):
            title = item.text
            link = item["href"]
            description = ""  # DuckDuckGo descriptions might need extra parsing
            results.append({"title": title, "link": link, "description": description})
            if len(results) >= max_results:
                break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur DuckDuckGo: {e}")
        return []

def torch_search(query, max_results=5):
    print("[Torch Search]")
    url = f"https://www.torchtorsearch.com/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.find_all("a", href=True):
            title = item.text.strip()
            link = item["href"]
            description = ""  # Torch descriptions might need extra parsing
            results.append({"title": title, "link": link, "description": description})
            if len(results) >= max_results:
                break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur Torch: {e}")
        return []

def main():
    print("\n=== Multi-Search Engine Crawler ===")
    query = input("Entrez le mot-clé ou la phrase à rechercher : ")

    print("\nRecherche en cours...\n")

    google_results = google_search(query)
    bing_results = bing_search(query)
    duckduckgo_results = duckduckgo_search(query)
    torch_results = torch_search(query)

    all_results = google_results + bing_results + duckduckgo_results + torch_results

    print("=== Résultats combinés ===")
    if all_results:
        for idx, result in enumerate(all_results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
    else:
        print("Aucun résultat pertinent trouvé.")

if __name__ == "__main__":
    main()
