from argparse import ArgumentParser
from canopsis.webcore.services.context_graph import ImportKey
from canopsis.context_graph.import_ctx import Manager, ContextGraphImport
import signal
import logging
import os
import sys
import time

#TODO: stack the import −> add a counter or use a mutex
#TODO: change the logging format to match the one used on canopsis

ROOT = "/opt/canopsis/"
UMASK = 0

E_DAEMON_CREATION = "Error while creating the daemon."
E_CHANGE_DIR = "Impossible to change the current working directory to " + \
                   ROOT + "."
E_IMPORT_FAILED = "Error during the import of id {0} : {1}."
I_IMPORT_DONE = "Import {0} done."
I_START_IMPORT = "Start import {0}."
I_DAEMON_RUNNING = "The daemon is running with the pid {0}."


def import_handler(signum, stack):

    importer = ContextGraphImport()
    manager = Manager()

    uuid = manager.get_next_uuid()

    logging.info(I_START_IMPORT.format(uuid))
    manager.update_status(uuid, {ImportKey.F_STATUS: ImportKey.ST_ONGOING})

    start = time.time()
    report = {}
    try:
        importer.import_context(uuid)

    except Exception as e:
        report = {ImportKey.F_STATUS: ImportKey.ST_FAILED,
                  ImportKey.F_INFO: str(e)}
        logging.error(E_IMPORT_FAILED.format(uuid, e))

    else:
        report = {ImportKey.F_STATUS: ImportKey.ST_DONE}
        logging.info(I_IMPORT_DONE.format(uuid))

    end = time.time()
    report[ImportKey.F_EXECTIME] = end - start
    manager.update_status(uuid, report)

    del(importer)
    del(manager)

def daemon_loop():
    signal.signal(signal.SIGUSR1, import_handler)
    while True: # Main loop. Weee
        signal.pause()

def daemonize(function):
    try:
        pid = os.fork()
    except OSError:
        logging.error(E_DAEMON_CREATION)
        exit(1)

    if pid > 0: # parent
        exit(0)

    else: # child
        sid = os.setsid()
        signal.signal(signal.SIGUSR1, signal.SIG_IGN)
        signal.signal(signal.SIGINT, signal.SIG_IGN)
        signal.signal(signal.SIGTERM, signal.SIG_IGN)

        try:
            pid = os.fork()
        except OSError:
            logging.error(E_DAEMON_CREATION)
            exit(1)

        if pid > 0: # parent
            exit(0)
        else:
            try:
                os.chdir(ROOT)
            except OSError:
                logging.error(E_CHANGE_DIR)

                os.umask(UMASK)

        fd_to_close = [sys.stderr, sys.stdout, sys.stdin]

        for fd in fd_to_close:
            try:
                fd.close()
            except OSError:
                pass

        pid = os.getpid()

        with open(ImportKey.PID_FILE, 'w') as fd:
            fd.write("{0}".format(pid))

        logging.info(I_DAEMON_RUNNING.format(pid))

        function()



def main():
    logging.basicConfig(filename='/opt/canopsis/var/log/impord.log',
                        level=logging.DEBUG)
    daemonize(daemon_loop)

if __name__ == "__main__":
    main()
