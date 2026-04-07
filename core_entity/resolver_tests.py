from core import Resolver

bucket_path = "C:\\Users\\Aya\\Desktop\\game-get\\bucket\\bucket-entity-based"

r = Resolver(bucket_path=bucket_path,package_name="cyberpunk",version="1.6")
r.resolve()