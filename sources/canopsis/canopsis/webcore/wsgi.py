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

import importlib
import os
import sys

import uwsgidecorators
import beaker
import beaker.cache
import beaker.session

from bottle import default_app as BottleApplication, HTTPError
from beaker.middleware import SessionMiddleware

from canopsis.vendor import mongodb_beaker
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_array
from canopsis.logger import Logger
from canopsis.old.account import Account
# TODO: replace with canopsis.mongo.MongoStorage
from canopsis.old.storage import get_storage

from canopsis.webcore.services import session as session_module

DEFAULT_DEBUG = False
DEFAULT_ECSE = False
DEFAULT_ROOT_DIR = '~/var/www/src/'
DEFAULT_COOKIES_EXPIRE = 300
DEFAULT_SECRET = 'canopsis'
DEFAULT_DATA_DIR = '~/var/cache/canopsis/webcore/'


class EnsureAuthenticated(object):
    name = 'EnsureAuthenticated'
    handle_logout = False

    def __init__(self, ws, *args, **kwargs):
        super(EnsureAuthenticated, self).__init__(*args, **kwargs)

        self.ws = ws
        self.session = session_module

    def apply(self, callback, context):
        def decorated(*args, **kwargs):
            s = self.session.get()

            if not s.get('auth_on', False):
                return HTTPError(401, 'Not authorized')

            return callback(*args, **kwargs)

        return decorated


class WebServer():

    CONF_PATH = 'etc/webserver.conf'
    LOG_FILE = 'var/log/webserver.log'

    @property
    def application(self):
        return self.app

    @property
    def beaker_url(self):
        return '{}.beaker'.format(self.db.uri)

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

    def __init__(self, config, logger):
        self.config = config
        self.logger = logger

        server = self.config.get('server', {})
        self.debug = server.get('debug', DEFAULT_DEBUG)
        self.enable_crossdomain_send_events = server.get(
            'enable_crossdomain_send_events', DEFAULT_ECSE)
        self.root_directory = os.path.expanduser(
            server.get('root_directory', DEFAULT_ROOT_DIR))

        auth = self.config.get('auth', {})
        self.providers = cfg_to_array(auth.get('providers', ''))
        if len(self.providers) == 0:
            self.logger.critical('Missing providers. Cannot launch webcore module.')
            raise RuntimeError('Missing providers')

        session = self.config.get('session', {})
        self.cookie_expires = int(session.get('cookie_expires',
                                              DEFAULT_COOKIES_EXPIRE))
        self.secret = session.get('secret', DEFAULT_SECRET)
        self.data_dir = session.get('data_dir', DEFAULT_DATA_DIR)

        self.webservices = self.config.get('webservices', {})

        # TODO: Replace with MongoStorage
        self.db = get_storage(account=Account(user='root', group='root'))
        self.stopping = False

        self.webmodules = {}
        self.auth_backends = {}

    def init_app(self):
        self.logger.info('Initialize WSGI Application')
        self.app = BottleApplication()

        self.load_auth_backends()
        self.load_webservices()
        self.load_session()

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

    def _load_auth_backend(self, modname):
        if modname in self.auth_backends:
            return True

        self.logger.info('Load authentication backend: {0}'.format(modname))

        try:
            mod = importlib.import_module(modname)

        except ImportError as err:
            self.logger.error(
                'Impossible to load authentication backend {}: {}'.format(
                    modname, err
                )
            )

            return False

        else:
            backend = mod.get_backend(self)
            self.auth_backends[backend.name] = backend
            self.app.install(backend)

        return True

    def load_auth_backends(self):
        for provider in self.providers:
            self._load_auth_backend(provider)

        # Always add this backend which returns 401 when the login fails
        backend = EnsureAuthenticated(self)
        self.auth_backends[backend.name] = backend
        self.app.install(backend)

    def load_session(self):
        # Since we vendor mongodb_beaker because of broken dep on pypi.python.org
        # we need to setup the beaker class map manually.
        beaker.cache.clsmap['mongodb'] = mongodb_beaker.MongoDBNamespaceManager
        self.app = SessionMiddleware(self.app, {
            'session.type': 'mongodb',
            'session.cookie_expires': self.cookie_expires,
            'session.url': self.beaker_url,
            'session.secret': self.secret,
            'session.lock_dir': self.data_dir
        })

    def unload_session(self):
        pass

    def unload_auth_backends(self):
        pass

    def unload_webservices(self):
        pass

    @uwsgidecorators.signal(9)
    @uwsgidecorators.signal(15)
    def exit(self, num):
        if not self.stopping:
            self.stopping = True

            self.unload_session()
            self.unload_webservices()
            self.unload_auth_backends()

            sys.exit(0)

    class Error(Exception):
        pass


def get_default_app():
    conf = Configuration.load(WebServer.CONF_PATH, Ini)
    logger = Logger.get('webserver', WebServer.LOG_FILE)
    # Declare WSGI application
    ws = WebServer(config=conf, logger=logger).init_app()
    app = ws.application
    return app
