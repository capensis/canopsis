

import time


class Timer:
    def __init__(self, logger):
        self.logger = logger
        self.start_time = 0
        self.action = ""

    def start(self, action):
        self.action = action
        self.logger.info("Starting timing for action: {}".format(action))
        self.start_time =  int(round(time.time() * 1000))

    def stop(self):
        stop_time = int(round(time.time() * 1000))
        duration = stop_time - self.start_time
        self.logger.info("Action : {0} took {1} ms to complete ".format(self.action, duration))
        self.reset()

    def reset(self):
        self.start_time = 0
        self.action = ""
