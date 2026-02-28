import requests
from bs4 import BeautifulSoup
import uuid


class Game:
    def __init__(self,id: uuid.UUID,title: str,download_page: str):
        self.id = id
        self.title = title
        self.download_page = download_page


class SearchResult:
    def __init__(self,parent_site: "site",query: str,games: dict[uuid.UUID,Game]):
        self.parent_site = parent_site
        self.query = query
        self.games = games
        
        self.next = None
        self.prev = None


    def get_next(self):
        pass
    def get_prev(self):
        pass






class Steamumu:
    def __init__(self):
        self.base_url = "https://steamunlocked.org/"
    
    def search_game(self,query: str) -> SearchResult | None:
        query = query.strip()
        query = query.replace(" ","+")

        response = requests.get(f"{self.base_url}/?s={query}", headers={"User-Agent": "Mozilla/5.0"})
        soup = BeautifulSoup(response.text, "html.parser")

        first_h1 = soup.find("h1")
        if first_h1.get_text() == "No results found":
            return None


        games: dict[uuid.UUID,Game] = {}

        for h1 in soup.find_all("h1"):
            #print(h1)
            title = h1.get_text(strip=True)
            parent_link = h1.find_parent("a")
            href = parent_link["href"] if parent_link else "No link"

            id = uuid.uuid4()
            g = Game(id,title,href)
            games[id] = g
            #print(f"{title} -> {href}")


        return SearchResult(self,query,games)
            






def run_test():
    site = Steamumu()
    r: SearchResult | None = site.search_game("far cry")
    if r:
        print(r.games)
        for v in r.games.values():

            print(v.title)

        r.get_next()



run_test()