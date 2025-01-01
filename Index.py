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
        for item in soup.select("div.results_links_deep"):
            title = item.select_one("a.result__a").text if item.select_one("a.result__a") else ""
            link = item.select_one("a.result__a")["href"] if item.select_one("a.result__a") else ""
            description = item.select_one("a.result__snippet").text if item.select_one("a.result__snippet") else ""
            if link:
                results.append({"title": title, "link": link, "description": description})
                if len(results) >= max_results:
                    break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur DuckDuckGo: {e}")
        return []

def qwant_search(query, max_results=5):
    print("[Qwant Search]")
    url = f"https://www.qwant.com/?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.select("div.result"):
            title = item.select_one("a.result__url").text if item.select_one("a.result__url") else ""
            link = item.select_one("a.result__url")["href"] if item.select_one("a.result__url") else ""
            description = item.select_one("p.result__desc").text if item.select_one("p.result__desc") else ""
            if link:
                results.append({"title": title, "link": link, "description": description})
                if len(results) >= max_results:
                    break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur Qwant: {e}")
        return []

def ecosia_search(query, max_results=5):
    print("[Ecosia Search]")
    url = f"https://www.ecosia.org/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.find_all("div", class_="result__body"):
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
        print(f"Erreur Ecosia: {e}")
        return []

def main():
    print("\n=== Multi-Search Engine Crawler ===")
    query = input("Entrez le mot-clé ou la phrase à rechercher : ")

    print("\nRecherche en cours...\n")

    google_results = google_search(query)
    duckduckgo_results = duckduckgo_search(query)
    qwant_results = qwant_search(query)
    ecosia_results = ecosia_search(query)

    all_results = google_results + duckduckgo_results + qwant_results + ecosia_results

    print("=== Résultats combinés ===")
    if all_results:
        for idx, result in enumerate(all_results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
    else:
        print("Aucun résultat pertinent trouvé.")

if __name__ == "__main__":
    main()
