# cli.py
from complier.pipeline import Pipeline
#from pathlib import Path
import argparse
import sys
from complier.errors import LoadError,ConfigValidationError

UNEXPECTED_LOAD_ERROR = 19
UNEXPECTED_CHECK_ERROR = 29


def main():
    parser = argparse.ArgumentParser(prog="pact")
    sub = parser.add_subparsers(dest="command")

    # build <name>
    build = sub.add_parser("build")
    build.add_argument("name", type=str)

    args = parser.parse_args()

    if args.command == "build":

        build_cmd(args.name)
    

    else:
        parser.print_help()




def build_cmd(path: str):

    try:
        pipeline = Pipeline(path)
    except LoadError as e:
        print(e)
        sys.exit(e.exit_code)

    except Exception as e:
        print("an uanexpacted error ecoured")
        # add to error log file 
        print(e)
        sys.exit(UNEXPECTED_LOAD_ERROR)


    try:
        pipeline.run()
    except ConfigValidationError as e:
        print(e)
        sys.exit(e.exit_code)
    except Exception as e:
        print("an uanexpacted error ecoured")
        print(e)
        sys.exit(UNEXPECTED_CHECK_ERROR)










if __name__ == "__main__":
    main()