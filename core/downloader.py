from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from core.models import Package, Version


DOWNLOAD_PRIORITY = ["direct", "torrent", "scrape"]




def download_package(pkg: Package):
    version: Version = pkg.versions[pkg.default]

    #download_type = None
    #download_links: list[str] = []

    for dt in DOWNLOAD_PRIORITY:
        if dt in version.downloads:
         
            download_links: list[str] = version.downloads[dt]
            break

    if not download_links: ## dev check 
        raise RuntimeError(f"No downloads available for version {version.id}")

    print(f"Downloading {download_type} links: {download_links}")
