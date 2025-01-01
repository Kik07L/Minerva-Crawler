import requests
import urllib.parse
from bs4 import BeautifulSoup

def google_search(query, max_results=5):
    print("[Google Search]")
    url = f"https://www.google.com/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    response = requests.get(url, headers=headers)
    soup = BeautifulSoup(response.text, "html.parser")

    results = []
    for item in soup.select("a"):  # Google often changes class names; adjust as needed
        link = item.get("href")
        if link and link.startswith("/url?q="):
            clean_link = link.split("/url?q=")[1].split("&sa=")[0]
            results.append(clean_link)
            if len(results) >= max_results:
                break
    for idx, link in enumerate(results):
        print(f"{idx + 1}. {link}")
    print()
    return results

def bing_search(query, max_results=5):
    print("[Bing Search]")
    url = f"https://www.bing.com/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    response = requests.get(url, headers=headers)
    soup = BeautifulSoup(response.text, "html.parser")

    results = []
    for item in soup.find_all("a", href=True):
        link = item.get("href")
        if link and link.startswith("http"):
            results.append(link)
            if len(results) >= max_results:
                break
    for idx, link in enumerate(results):
        print(f"{idx + 1}. {link}")
    print()
    return results

def duckduckgo_search(query, max_results=5):
    print("[DuckDuckGo Search]")
    url = f"https://duckduckgo.com/html/?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    response = requests.get(url, headers=headers)
    soup = BeautifulSoup(response.text, "html.parser")

    results = []
    for item in soup.find_all("a", href=True):
        link = item.get("href")
        if link.startswith("https://"):
            results.append(link)
            if len(results) >= max_results:
                break
    for idx, link in enumerate(results):
        print(f"{idx + 1}. {link}")
    print()
    return results

def torch_search(query, max_results=5):
    print("[Torch Search]")
    url = f"https://www.torchtorsearch.com/search?q={urllib.parse.quote(query)}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    response = requests.get(url, headers=headers)
    soup = BeautifulSoup(response.text, "html.parser")

    results = []
    for item in soup.find_all("a", href=True):
        link = item.get("href")
        if link.startswith("http"):
            results.append(link)
            if len(results) >= max_results:
                break
    for idx, link in enumerate(results):
        print(f"{idx + 1}. {link}")
    print()
    return results

def main():
    print("\n=== Multi-Search Engine Crawler ===")
    query = input("Entrez le mot-clé ou la phrase à rechercher : ")

    print("\nRecherche en cours...\n")

    google_links = google_search(query)
    bing_links = bing_search(query)
    duckduckgo_links = duckduckgo_search(query)
    torch_links = torch_search(query)

    print("=== Résultats combinés ===")
    unique_links = set(google_links + bing_links + duckduckgo_links + torch_links)
    for idx, link in enumerate(unique_links):
        print(f"{idx + 1}. {link}")

if __name__ == "__main__":
    main()
