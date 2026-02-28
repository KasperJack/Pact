import requests
import os


import sys
sys.stdout.reconfigure(encoding='utf-8')

from dotenv import load_dotenv

load_dotenv()

API_KEY = os.getenv("RAWG_API_KEY")

def search_rawg(query: str):
    url = "https://api.rawg.io/api/games"
    
    params = {
        "key": API_KEY,
        "search": query,
        "platforms": 4,
        "page_size": 50,
        'game_type': 'game'
    }

    response = requests.get(url, params=params)

    if response.status_code != 200:
        print("Error:", response.text)
        return []

    data = response.json()
    #print(data)
    results = data.get("results", [])

    for game in results:
        print(f"{game['name']} ({game.get('released', 'N/A')})")
        
    #print('=====================')
    return results


search_rawg("far cry")