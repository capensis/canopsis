#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

import logging
import signal
import time
import os


class Init(object):

    class getHandler(object):
        def __init__(self, logger):
            self.logger = logger
            self.RUN = True

        def status(self):
            return self.RUN

        def signal_handler(self, signum, frame):
            self.logger.warning("Receive signal to stop daemon...")
            if self.callback:
                self.callback()
            self.stop()

        def run(self, callback=None):
            self.callback = callback
            signal.signal(signal.SIGINT, self.signal_handler)
            signal.signal(signal.SIGTERM, self.signal_handler)

        def stop(self):
            self.RUN = False

        def set(self, statut):
            self.RUN = statut

        def wait(self):
            while self.RUN:
                try:
                    time.sleep(1)
                except:
                    break
            self.stop()

    def getLogger(self, name, level="INFO", logging_level=None):
        if logging_level is None:
            self.level = level
        else:
            self.level = logging_level

        logging.basicConfig(format='%(asctime)s %(levelname)s %(name)s [%(module)s %(lineno)s] %(message)s')
        self.logger = logging.getLogger(name)
        self.logger.setLevel(self.level)

        return self.logger

    def get_confpath(self, conftype):
        """
            Get path to config file.

            :param conftype: Type of configuration (webserver, websocket, amqp, storage, ...)
            :type conftype: basestring

            :returns: Absolute path to config file
        """

        envvar = 'CPS_CONFPATH_{0}'.format(conftype.upper())

        if conftype == 'webcore':
            default = '~/etc/webserver.conf'

        elif conftype == 'websocket':
            default = '~/etc/websocket.conf'

        elif conftype == 'amqp':
            default = '~/etc/amqp.conf'

        elif conftype == 'storage':
            default = '~/etc/cstorage.conf'

        elif conftype == 'engines':
            default = '~/etc/amqp2engines'

        elif conftype == 'logging':
            default = '~/etc/logging.conf'

        else:
            default = None

        return os.path.expanduser(os.getenv(envvar, default))
