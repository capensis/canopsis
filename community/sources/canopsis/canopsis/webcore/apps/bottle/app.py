# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals
from canopsis.common import root_path
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.logger import Logger
from canopsis.confng import Configuration, Ini
from canopsis.common.amqp import AmqpPublisher
from canopsis.common.amqp import get_default_connection as \
    get_default_amqp_connection
from bottle import default_app as BottleApplication
from signal import SIGTERM, SIGINT
import bottle
import sys
import os
import importlib
import gevent

# DO NOT EVER MODIFY THE 2 LINES BELOW OR UNDESIRED BEHAVIOR ***WILL** HAPPEN.
from gevent import monkey
monkey.patch_all()


# TODO: replace with canopsis.mongo.MongoStorage

DEFAULT_DEBUG = False
DEFAULT_ECSE = False
DEFAULT_ROOT_DIR = '~/var/www/src/'
DEFAULT_DATA_DIR = '~/var/cache/canopsis/webcore/'

bottle.BaseRequest.MEMFILE_MAX = 1024*1024*1024


class OldApi():

    CONF_PATH = 'etc/oldapi.conf'
    LOG_FILE = root_path + '/var/log/oldapi.log'

    @property
    def application(self):
        return self.app

    @property
    def beaker_url(self):
        return self.db.beaker_uri

    @property
    def skip_login(self):
        return [bname for bname in self.auth_backends.keys()]

    @property
    def skip_logout(self):
        return [
            bname
            for bname in self.auth_backends.keys()
            if not self.auth_backends[bname].handle_logout
        ]

    def __init__(self, config, logger, amqp_pub):
        self.config = config
        self.logger = logger
        self.amqp_pub = amqp_pub

        server = self.config.get('server', {})
        self.debug = server.get('debug', DEFAULT_DEBUG)
        self.enable_crossdomain_send_events = server.get(
            'enable_crossdomain_send_events', DEFAULT_ECSE)
        self.root_directory = os.path.expanduser(
            server.get('root_directory', DEFAULT_ROOT_DIR))

        self.webservices = self.config.get('webservices', {})

        # TODO: Replace with MongoStorage
        self.db = get_storage(account=Account(user='root', group='root'))
        self.stopping = False

        self.webmodules = {}
        self.auth_backends = {}

    def init_app(self):
        self.logger.info('Initialize gevent signal-handlers')
        gevent.signal(SIGTERM, self.exit)
        gevent.signal(SIGINT, self.exit)
        self.logger.info('Initialize WSGI Application')
        self.app = BottleApplication()

        self.load_webservices()

        self.logger.info('WSGI fully loaded.')
        return self

    def _load_webservice(self, modname):
        if modname in self.webmodules:
            return True

        if modname is None:
            return False

        self.logger.info('Loading webservice: {0}'.format(modname))

        try:
            mod = importlib.import_module(modname)

        except ImportError as err:
            self.logger.error(
                'Impossible to load webservice {0}: {1}'.format(modname, err)
            )

            return False

        else:
            if hasattr(mod, 'exports'):
                self.webmodules[modname] = mod
                mod.exports(self)

            else:
                self.logger.error(
                    'Invalid module {0}, no exports()'.format(modname)
                )

                return False

        return True

    def load_webservices(self):
        for module in sorted(self.webservices.keys()):
            enable = int(self.webservices[module])
            if enable == 1:
                self._load_webservice(module)
            else:
                self.logger.info(
                    u'Webservice {} skipped by configuration.'.format(module))

        self.logger.info(u'Service loading completed.')

    def unload_webservices(self):
        pass

    def exit(self):
        if not self.stopping:
            self.stopping = True

            self.unload_webservices()
            self.amqp_pub.connection.disconnect()

            sys.exit(0)

    class Error(Exception):
        pass


def get_default_app(logger=None, oldapiconf=None, amqp_conn=None, amqp_pub=None):
    if oldapiconf is None:
        oldapiconf = Configuration.load(OldApi.CONF_PATH, Ini)

    if logger is None:
        logger = Logger.get('oldapi', OldApi.LOG_FILE)

    if amqp_conn is None:
        amqp_conn = get_default_amqp_connection()

    if amqp_pub is None:
        amqp_pub = AmqpPublisher(amqp_conn, logger)

    # Declare WSGI application
    ws = OldApi(config=oldapiconf, logger=logger, amqp_pub=amqp_pub).init_app()
    app = ws.application
    return app
