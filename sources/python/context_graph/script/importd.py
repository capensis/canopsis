from argparse import ArgumentParser
import signal
import time
import os

PID_FILE = "/opt/canopsis/var/run/importd.pid"

def signal_handler(signum, stack):
    fd = open("/tmp/importd.log", 'a')
    fd.write("A noze a signal USR1...\n")
    fd.close()

def load_conf():
    pass

def argparse():
    parser = ArgumentParser()
    # if needed, add stuffs here.
    args = parser.parse_args()
    return args

def main():
    #TODO create a daemon.
    fd = open("/tmp/importd.pid", 'w')
    fd.write(str(os.getpid()))
    fd.close()

    signal.signal(signal.SIGUSR1, signal_handler)

    while True: # Main loop. Weee
        signal.pause()


if __name__ == "__main__":
    main()
