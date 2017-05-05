from argparse import ArgumentParser
from canopsis.webcore.services.context_graph import ImportKey
from canopsis.context_graph.import_ctx import Manager, ContextGraphImport
import atexit
import signal
import time
import os

def signal_handler(signum, stack):
    #TODO: handle stat
    fd = open("/tmp/importd.log", 'a')
    fd.write("A noze a signal USR1... Let's do an import")
    importer = ContextGraphImport()
    manager = Manager()
    uuid = manager.get_next_uuid()

    fd.write("Get uuid done")
    manager.update_status(uuid, {ImportKey.F_STATUS: ImportKey.ST_ONGOING})
    fd.write("Update a status")

    start = time.time()
    report = {}
    try:
        importer.import_context(uuid)

    except Exception as e:
        report = {ImportKey.F_STATUS: ImportKey.ST_FAILED,
                  ImportKey.F_INFO: str(e)}
        fd.write("Failed\n")

    else:
        report = {ImportKey.F_STATUS: ImportKey.ST_DONE}
        fd.write("Done\n")

    finally:
        fd.close()
        end = time.time()
        report[ImportKey.F_EXECTIME] = end - start
        manager.update_status(uuid, report)

def load_conf():
    pass

def argparse():
    parser = ArgumentParser()
    # if needed, add stuffs here.
    args = parser.parse_args()
    return args

def cleanup():
    os.remove(ImportKey.PID_FILE)

def start_up():
    print(os.getpid())
    atexit.register(cleanup)
    try:
        fd = open(ImportKey.PID_FILE, 'w')
        fd.write(str(os.getpid()))
        fd.close()
    except IOError:
        exit(1)

    signal.signal(signal.SIGUSR1, signal_handler)

def main():
    start_up()

    while True: # Main loop. Weee
        signal.pause()


if __name__ == "__main__":
    main()
