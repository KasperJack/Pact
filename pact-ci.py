# cli.py
from complier.pipeline import Compiler
from complier.errors import PackageNotFoundError
from pathlib import Path
import argparse

def main():
    parser = argparse.ArgumentParser(prog="pact")
    sub = parser.add_subparsers(dest="command")

    # build <name>
    build = sub.add_parser("build")
    build.add_argument("name", type=str)

    args = parser.parse_args()

    if args.command == "build":

        #print(f"building {args.name}...")
        try:
            compiler = Compiler(get_curent_working_bucket(),args.name)
            compiler.run()

        except PackageNotFoundError as e:
            print(e)
        except Exception as e:
            print(f"faild to complile pacakge {args.name}")
            print(e)


    else:
        parser.print_help()











def get_curent_working_bucket() -> Path:

    return Path.cwd() / "test-buckets" / "defult"





if __name__ == "__main__":
    main()