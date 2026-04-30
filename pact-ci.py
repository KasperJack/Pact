# cli.py
from complier.pipeline import Pipeline
from pathlib import Path
import argparse
import sys
from pact_core.models.package import PackageModel


def main():
    parser = argparse.ArgumentParser(prog="pact")
    sub = parser.add_subparsers(dest="command")

    # build <name>
    build = sub.add_parser("build")
    build.add_argument("name", type=str)

    args = parser.parse_args()

    if args.command == "build":

        full_path = Path(args.name).resolve()
        
        if not full_path.exists():
            print(f"Error: file not found: {full_path}")
            sys.exit(1)


        if not full_path.is_file():
            print(f"Error: expected a file, got a directory: {full_path}")
            sys.exit(1)


        #print(full_path)
        pipeline = Pipeline(full_path)
        pipeline.run()


    else:
        parser.print_help()











def get_curent_working_bucket() -> Path:

    return Path.cwd() / "test-buckets" / "defult"





if __name__ == "__main__":
    main()