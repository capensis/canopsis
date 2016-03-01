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

import gevent
from gevent import monkey
monkey.patch_all()

from bottle import default_app as BottleApplication, HTTPError
from beaker.middleware import SessionMiddleware
import mongodb_beaker  # needed by beaker

from canopsis.configuration.model import Parameter, ParamList
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_config
from canopsis.configuration.configurable import Configurable
from canopsis.common.utils import setdefaultattr

# TODO: replace with canopsis.mongo.MongoStorage
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.old.rabbitmq import Amqp

from signal import SIGTERM, SIGINT

import importlib
import sys
import os


config = {
    'server': (
        Parameter('debug', parser=Parameter.bool),
        Parameter('enable_crossdomain_send_events', parser=Parameter.bool),
        Parameter('root_directory', parser=Parameter.path)
    ),
    'auth': (
        Parameter('providers', parser=Parameter.array(), critical=True)
    ),
    'session': (
        Parameter('cookie_expires', parser=int),
        Parameter('secret'),
        Parameter('data_dir', parser=Parameter.path)
    ),
    'webservices': ParamList(parser=Parameter.bool)
}


class EnsureAuthenticated(object):
    name = 'EnsureAuthenticated'
    handle_logout = False

    def __init__(self, ws, *args, **kwargs):
        super(EnsureAuthenticated, self).__init__(*args, **kwargs)

        self.ws = ws
        self.session = ws.require('session')

    def apply(self, callback, context):
        def decorated(*args, **kwargs):
            s = self.session.get()

            if not s.get('auth_on', False):
                return HTTPError(401, 'Not authorized')

            return callback(*args, **kwargs)

        return decorated


@add_config(config)
@conf_paths('webserver.conf')
class WebServer(Configurable):
    @property
    def debug(self):
        return setdefaultattr(self, '_debug', False)

    @debug.setter
    def debug(self, value):
        self._debug = value

    @property
    def enable_crossdomain_send_events(self):
        return setdefaultattr(self, '_crossdomain_evt', False)

    @enable_crossdomain_send_events.setter
    def enable_crossdomain_send_events(self, value):
        self._crossdomain_evt = value

    @property
    def root_directory(self):
        return setdefaultattr(
            self, '_rootdir',
            os.path.expanduser('~/var/www/src/')
        )

    @root_directory.setter
    def root_directory(self, value):
        value = os.path.expanduser(value)

        if os.path.exists(value):
            self._rootdir = value

    @property
    def providers(self):
        return setdefaultattr(self, '_providers', [])

    @providers.setter
    def providers(self, value):
        self._providers = value

    @property
    def cookie_expires(self):
        return setdefaultattr(self, '_cookie', 300)

    @cookie_expires.setter
    def cookie_expires(self, value):
        self._cookie = value

    @property
    def secret(self):
        return setdefaultattr(self, '_secret', 'canopsis')

    @secret.setter
    def secret(self, value):
        self._secret = value

    @property
    def data_dir(self):
        return setdefaultattr(
            self, '_datadir',
            os.path.expanduser('~/var/cache/canopsis/webcore/')
        )

    @data_dir.setter
    def data_dir(self, value):
        value = os.path.expanduser(value)

        if os.path.exists(value):
            self._datadir = value

    # dict properties do not need setters

    @property
    def webservices(self):
        if not hasattr(self, '_webservices'):
            self._webservices = {}

        return self._webservices

    @property
    def beaker_url(self):
        return '{0}.beaker'.format(self.db.uri)

    def __init__(self, *args, **kwargs):
        super(WebServer, self).__init__(*args, **kwargs)

        self.log_name = 'webserver'

        # TODO: Replace with MongoStorage
        self.db = get_storage(account=Account(user='root', group='root'))
        self.amqp = Amqp()
        self.stopping = False

        self.webmodules = {}
        self.auth_backends = {}

    def __call__(self):
        self.logger.info('Initialize gevent signal-handlers')
        gevent.signal(SIGTERM, self.exit)
        gevent.signal(SIGINT, self.exit)

        self.logger.info('Start AMQP thread')
        self.amqp.start()

        self.logger.info('Initialize WSGI Application')
        self.app = BottleApplication()

        self.load_auth_backends()
        self.load_webservices()
        self.load_session()

        return self

    def _load_webservice(self, name):
        modname = 'canopsis.webcore.services.{0}'.format(name)

        if name in self.webmodules:
            return True

        self.logger.info('Loading webservice: {0}'.format(name))

        try:
            mod = importlib.import_module(modname)

        except ImportError as err:
            self.logger.error(
                'Impossible to load webservice {0}: {1}'.format(name, err)
            )

            return False

        else:
            if hasattr(mod, 'exports'):
                self.webmodules[name] = mod
                mod.exports(self)

            else:
                self.logger.error(
                    'Invalid module {0}, no exports()'.format(name)
                )

                return False

        return True

    def load_webservices(self):
        for webservice in self.webservices:
            if self.webservices[webservice]:
                self._load_webservice(webservice)

    def _load_auth_backend(self, name):
        modname = 'canopsis.auth.{0}'.format(name)

        if name in self.auth_backends:
            return True

        self.logger.info('Load authentication backend: {0}'.format(name))

        try:
            mod = importlib.import_module(modname)

        except ImportError as err:
            self.logger.error(
                'Impossible to load authentication backend {}: {}'.format(
                    name, err
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

    def exit(self):
        if not self.stopping:
            self.stopping = True

            self.unload_session()
            self.unload_webservices()
            self.unload_auth_backends()

            self.amqp.stop()
            # TODO: self.amqp.wait() not implemented

            sys.exit(0)

    @property
    def application(self):
        return self.app

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

    def require(self, modname):
        if not self._load_webservice(modname):
            raise ImportError(
                'Impossible to import webservice: {0}'.format(modname)
            )

        return self.webmodules[modname]

    class Error(Exception):
        pass

# Declare WSGI application
ws = WebServer()()
app = ws.application
