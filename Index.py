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

def yahoo_search(query, max_results=5):
    print("[Yahoo Search]")
    url = f"https://search.yahoo.com/search?p={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.find_all("div", class_="dd algo algo-sr Sr"):  # Yahoo specific class names
            title = item.find("h3").text if item.find("h3") else ""
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
        print(f"Erreur Yahoo: {e}")
        return []

def yandex_search(query, max_results=5):
    print("[Yandex Search]")
    url = f"https://yandex.com/search/?text={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    try:
        response = requests.get(url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")

        results = []
        for item in soup.find_all("li", class_="serp-item"):
            title = item.find("h2").text if item.find("h2") else ""
            link = item.find("a")["href"] if item.find("a") else ""
            description = item.find("div", class_="text-container").text if item.find("div", class_="text-container") else ""
            if link:
                results.append({"title": title, "link": link, "description": description})
                if len(results) >= max_results:
                    break
        for idx, result in enumerate(results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
        return results
    except Exception as e:
        print(f"Erreur Yandex: {e}")
        return []

def main():
    print("\n=== Multi-Search Engine Crawler ===")
    query = input("Entrez le mot-clé ou la phrase à rechercher : ")

    print("\nRecherche en cours...\n")

    google_results = google_search(query)
    bing_results = bing_search(query)
    yahoo_results = yahoo_search(query)
    yandex_results = yandex_search(query)

    all_results = google_results + bing_results + yahoo_results + yandex_results

    print("=== Résultats combinés ===")
    if all_results:
        for idx, result in enumerate(all_results):
            print(f"{idx + 1}. {result['title']}\n   {result['link']}\n   {result['description']}\n")
    else:
        print("Aucun résultat pertinent trouvé.")

if __name__ == "__main__":
    main()
