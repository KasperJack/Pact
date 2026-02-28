from core.loader import load_package



def show_package_info(package_name: str):
    pass


def install_package(package_name: str):

    try:
        pkg = load_package(package_name)
        print(pkg.igdb_id)
    except:
        Invalide


install_package("ion-fury")