from argparse import ArgumentParser
from canopsis.webcore.services.context_graph import ImportKey
from canopsis.context_graph.import_ctx import Manager, ContextGraphImport
import signal
import logging
import os
import sys
import time
import threading

# TODO: stack the import -> add a counter or use a mutex
# TODO: change the logging format to match the one used on canopsis

ROOT = "/opt/canopsis/"
UMASK = 0

E_DAEMON_CREATION = "Error while creating the daemon."
E_CHANGE_DIR = "Impossible to change the current working directory to " + \
    ROOT + "."
E_IMPORT_FAILED = "Error during the import of id {0} : {1}."
E_BEYOND_REPAIR = "An error beyond repair occured : {0}. Exiting."
I_IMPORT_DONE = "Import {0} done."
I_START_IMPORT = "Start import {0}."
I_DAEMON_RUNNING = "The daemon is running with the pid {0}."

import_mutex = threading.Lock()

counter = 0
manager = Manager()

def execution_time(exec_time):
    """Return from exec_time a human readable string that represent the
    execution time in a human readable format"""

    exec_time = int(exec_time) # we do not care of everything under the second

    hours =  exec_time / 3600
    minutes = (exec_time - 3600 * hours) / 60
    seconds = exec_time - (hours * 3600) - (minutes * 60)

    return "{0}:{1}:{2}".format(str(hours).zfill(2),
                                str(minutes).zfill(2),
                                str(seconds).zfill(2))

def process_import():
    importer = ContextGraphImport()

    uuid = manager.get_next_uuid()
    logging.info("Processing import {0}.".format(uuid))

    logging.info(I_START_IMPORT.format(uuid))
    manager.update_status(uuid, {ImportKey.F_STATUS: ImportKey.ST_ONGOING})

    start = time.time()
    report = {}
    try:
        updated, deleted = importer.import_context(uuid)

    except Exception as e:
        report = {ImportKey.F_STATUS: ImportKey.ST_FAILED,
                  ImportKey.F_INFO: str(e)}
        logging.error(E_IMPORT_FAILED.format(uuid, e))

    else:
        report = {ImportKey.F_STATUS: ImportKey.ST_DONE,
                  ImportKey.F_STATS:
                  {ImportKey.F_DELETED: deleted,
                   ImportKey.F_UPDATED: updated}}

        logging.info(I_IMPORT_DONE.format(uuid))

    end = time.time()
    report[ImportKey.F_EXECTIME] = execution_time(end - start)
    manager.update_status(uuid, report)

    del(importer)

def sig_usr1_handler(signum, stack):
    signal.signal(signal.SIGUSR1, signal.SIG_IGN)
    process_import()

    while manager.pending_in_db():
        process_import()

    signal.signal(signal.SIGUSR1, sig_usr1_handler)

def daemon_loop():
    signal.signal(signal.SIGUSR1, sig_usr1_handler)

    while True:  # Main loop. Weee
        signal.pause()

def daemonize(function):
    try:
        pid = os.fork()
    except OSError:
        logging.error(E_DAEMON_CREATION)
        exit(1)

    if pid > 0:  # parent
        exit(0)

    else:  # child
        sid = os.setsid()
        signal.signal(signal.SIGUSR1, signal.SIG_IGN)
        signal.signal(signal.SIGINT, signal.SIG_IGN)
        signal.signal(signal.SIGTERM, signal.SIG_IGN)

        try:
            pid = os.fork()
        except OSError:
            logging.error(E_DAEMON_CREATION)
            exit(1)

        if pid > 0:  # parent
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

        try:
            function()
        except Exception as e:
            logging.critical(E_BEYOND_REPAIR.format(e))


def main():


    logging.basicConfig(filename='/opt/canopsis/var/log/importd.log',
                        level=logging.DEBUG)
    daemonize(daemon_loop)

if __name__ == "__main__":
    main()
