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

import ConfigParser
import logging
import time
import sys
import os

import importlib
import imp

from bottle import route, static_file, redirect, request, \
    debug as bottle_debug, default_app, install

import libexec.auth

## Hack: Prevent "ExtractionError: Can't extract file(s) to egg cache" when 2 process extract egg at the same time ...
try:
    from beaker.middleware import SessionMiddleware
except:
    time.sleep(2)
    from beaker.middleware import SessionMiddleware

from signal import SIGTERM, SIGINT
from gevent import signal

# Reduce beaker verbose
import mongodb_beaker
mongodb_beaker.log.setLevel(logging.INFO)

#from gevent import monkey; monkey.patch_all()

## Configurations
config_filename = os.path.expanduser('~/etc/webserver.conf')
config = ConfigParser.RawConfigParser()
config.read(config_filename)

#put config in builtins to make it readable from modules in libexec
__builtins__["config"] = config

webservices = []
webservices_mods = {}

for webservice, enabled in config.items('webservices'):
    enabled = int(enabled)

    if enabled and webservice not in webservices:
        webservices.append(webservice)

mongo_config_file = os.path.expanduser('~/etc/cstorage.conf')
mongo_config = ConfigParser.RawConfigParser()
mongo_config.read(mongo_config_file)

## default config
debug = False

#mongo config
mongo_host = mongo_config.get("master", "host")
mongo_port = mongo_config.getint("master", "port")
mongo_userid = mongo_config.get("master", "userid")
mongo_password = mongo_config.get("master", "password")
mongo_db = mongo_config.get("master", "db")

session_cookie_expires = 300
session_secret = 'canopsis'
session_lock_dir = os.path.expanduser('~/tmp/webcore_cache')
root_directory = os.path.expanduser("~/var/www/")

if mongo_userid and mongo_password:
    session_mongo_url = 'mongodb://{0}:{1}@{2}:{3}/{4}.beaker'.format(
        mongo_userid,
        mongo_password,
        mongo_host,
        mongo_port,
        mongo_db
    )

else:
    session_mongo_url = 'mongodb://{0}:{1}/{2}.beaker'.format(
        mongo_host,
        mongo_port,
        mongo_db
    )

try:
    ## get config
    debug = config.getboolean('server', "debug")
    root_directory = os.path.expanduser(config.get('server', "root_directory"))

    session_cookie_expires = config.getint('session', "cookie_expires")
    session_secret = config.get('session', "secret")
    session_data_dir = os.path.expanduser(config.get('session', "data_dir"))

except Exception as err:
    print("Error when reading '%s' (%s)" % (config_filename, err))

## Logger
logging_level = logging.INFO
if debug:
    logging_level = logging.DEBUG

logging.basicConfig(
    format=r"%(asctime)s [%(process)d] [%(name)s] [%(levelname)s] %(message)s",
    datefmt=r"%Y-%m-%d %H:%M:%S", level=logging_level)
logger = logging.getLogger("webserver")

bottle_debug(debug)


## load and unload webservices
def load_webservices():
    for webservice in webservices:
        logger.info('Loading webservice: {0}'.format(webservice))

        modpath = os.path.expanduser('~/opt/webcore/libexec/{0}.py'.format(webservice))
        modname = webservice

        try:
            mod = imp.load_source(modname, modpath)

            if hasattr(mod, 'logger'):
                mod.logger.setLevel(logging_level)

            if hasattr(mod, 'load'):
                mod.load()

            webservices_mods[modname] = mod

        except Exception as err:
            logger.error('Impossible to load webservice {0}: {1}'.format(modname, err))


def unload_webservices():
    logger.info("Unload webservices.")
    for webservice in webservices_mods:
        module = webservices_mods[webservice]

        if hasattr(module, 'unload'):
            logger.info('Unloading module {0}'.format(webservice))
            module.unload()


def autoLogin(key=None):
    if key and len(key) == 56:
        logger.debug('Autologin:')
        output = libexec.auth.autoLogin(key)
        if output['success']:
            logger.info(' + Success')
            return True
        else:
            logger.info(' + Failed')
            return False

## Bind signals
stop_in_progress = False


def signal_handler():
    global stop_in_progress
    if not stop_in_progress:
        stop_in_progress = True
        logger.info("Receive signal to stop worker ...")
        unload_webservices()
        logger.info("Ready to stop.")
        sys.exit(0)

signal(SIGTERM, signal_handler)
signal(SIGINT, signal_handler)

## Bottle App
app = default_app()

load_webservices()

## Install plugins

if config.has_option('auth', 'providers'):
    providers = config.get('auth', 'providers').split(',')

    for provider in providers:
        logger.info('Loading authentication provider: {0}'.format(provider))

        modname = 'canopsis.auth.{0}'.format(provider)

        try:
            mod = importlib.import_module(modname)

        except ImportError as err:
            logger.error('Impossible to load authentication backend {0}: {1}'.format(modname, err))

        else:
            install(mod.get_backend())

install(libexec.auth.EnsureAuthenticated())

## Session system with beaker
session_opts = {
    'session.type': 'mongodb',
    'session.cookie_expires': session_cookie_expires,
    'session.url': session_mongo_url,
    'session.auto': True,
    #'session.timeout': 300,
    'session.secret': session_secret,
    'session.lock_dir': session_lock_dir,
}


## Basic Handler
@route('/:lang/static/canopsis/index.html')
@route('/static/canopsis/index.html')
def index(lang='en'):
    return static_file('canopsis/index.html', root=root_directory)


@route('/:lang/static/canopsis/index.debug.html')
@route('/static/canopsis/index.debug.html')
def index_debug(lang='en'):
    return static_file('canopsis/index.debug.html', root=root_directory)


@route('/:lang/static/:path#.+#', skip=libexec.auth.auth_backends)
@route('/static/:path#.+#', skip=libexec.auth.auth_backends)
def server_static(path, lang='en'):
    key = request.params.get('authkey', default=None)
    if key:
        autoLogin(key)

    return static_file(path, root=root_directory)


@route('/favicon.ico', skip=[libexec.auth.auth_backends])
def favicon():
    return


@route('/', skip=libexec.auth.auth_backends)
@route('/:key', skip=libexec.auth.auth_backends)
@route('/:lang/', skip=libexec.auth.auth_backends)
@route('/:lang/:key', skip=libexec.auth.auth_backends)
@route('/index.html', skip=libexec.auth.auth_backends)
@route('/:lang/index.html', skip=libexec.auth.auth_backends)
def loginpage(key=None, lang='en'):
    s = request.environ.get('beaker.session')

    key = key or request.params.get('authkey', default=None)

    if key:
        autoLogin(key)

    ticket = request.params.get('ticket', default=None)

    if not ticket and not s.get('auth_on', False):
        return static_file('login/index.html', root=root_directory)

    else:
        redirect('/{0}/static/canopsis/index.html'.format(lang))

## Install session Middleware
app = SessionMiddleware(app, session_opts)
