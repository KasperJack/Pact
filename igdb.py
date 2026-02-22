import requests
from dotenv import load_dotenv
import os 
import time

load_dotenv()




def search_games(query: str):

    headers = {
        "Client-ID": os.getenv("CLIENT-ID"),
        "Authorization": "Bearer " + os.getenv("AUTHORIZATION")
    }

    res = requests.post("https://api.igdb.com/v4/games",headers=headers,data=f'search "{query}"; fields id; limit 50;')

    ids = [str(g["id"]) for g in res.json()]
    
    if not ids:
        return []
    id_list = ",".join(ids)

    res2 = requests.post("https://api.igdb.com/v4/games",headers=headers,
                         data=f"""
                            fields name, slug, first_release_date;
                            where id = ({id_list})
                            & platforms = (6)
                            & version_parent = null
                            & game_type = 0
                            & first_release_date != null
                            & first_release_date < {int(time.time())};
                            limit 50;
                         """)

    for g in res2.json():
        print(g["name"])
    print('==================')


    res3 = requests.post("https://api.igdb.com/v4/games",headers=headers,
                         data=f"""
                            search "{query}";
                            fields name, slug, first_release_date;
                            where platforms = (6)
                            & version_parent = null
                            & game_type = 0
                            & first_release_date != null
                            & first_release_date < {int(time.time())};
                            limit 50;

                         """)
    for g in res3.json():
        print(g["name"])




search_games("witcher")