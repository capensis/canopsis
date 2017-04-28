#!/usr/bin/env python

from canopsis.context_graph.import_ctx import ContextGraphImport
from argparse import ArgumentParser

def parse_args():
    parser = ArgumentParser()
    parser.add_argument(type=str, dest="file", default="/opt/canopsis/tmp/import.json",
                        help="The file to parse.")
    args = parser.parse_args()
    return args

def main():
    args = parse_args()
    file_ = args.file

    import_manager = ContextGraphImport()
    try:
        import_manager.import_ctx(file_, 0)
    except Exception as e:
        fd = open("/opt/canopsis/tmp/status.log", 'w')
        fd.write(str(e))
        fd.close()
    finally:
        fd = open("/opt/canopsis/tmp/status.log", 'w')
        fd.write("I am ok")
        fd.close()


if __name__ == "__main__":
    main()
