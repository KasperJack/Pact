from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from core.models import Package, Version




def download_package(pkg: Package):
    version: Version = pkg.versions[pkg.default]



    if "direct" in version.downloads:
        direct_links: list[str] = version.downloads["direct"]
 
 
    if "torrent" in version.downloads:
        torrent_links:list[str] = version.downloads["torrent"]
    
    if "scrape" in version.downloads:
        torrent_links:list[str] = version.downloads["scrape"]