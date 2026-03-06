#from core.loader import load_package
from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from core.models import Package,Version


from core import load_package, download_package, InvalidManifestError, PackageNotFoundError
import sys








def show_package_info(pkg: Package):
    print(pkg.name)
    print(pkg.igdb_id)
    print(pkg.release_year)


def show_version_info(pkg: Package):
    v: Version = pkg.versions[pkg.default]
    print(v.version)
    print(v.source)
    print(v.size_mb)
    
    




def install_package(package_name: str):

    try:
        # try to check of defult version is defined porpely 
        # check missing keys errors 
        pkg = load_package(package_name)
    except PackageNotFoundError as e:
        print(e)
        # suggest another packages with close name 
        sys.exit(1)
    except InvalidManifestError as e:
        print(e)
        sys.exit(2)

    show_package_info(pkg)
    show_version_info(pkg)
    download_package(pkg)
    
install_package("ion-fury")