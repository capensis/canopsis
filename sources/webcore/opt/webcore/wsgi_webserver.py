#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

import os, sys, time, logging

import ConfigParser

import bottle
from bottle import route, run, static_file, redirect, request

import libexec.auth

## Hack: Prevent "ExtractionError: Can't extract file(s) to egg cache" when 2 process extract egg at the same time ...
try:
	from beaker.middleware import SessionMiddleware
except:
	time.sleep(2)
	from beaker.middleware import SessionMiddleware

# Reduce beaker verbose
import mongodb_beaker
mongodb_beaker.log.setLevel(logging.INFO)

#from gevent import monkey; monkey.patch_all()

## Configurations
webservices = ['account',  'auth', 'calendar_events', 'event', 'files', 'perfstore', 'reporting', 'rest', 'rights', 'ui_view', 'ui_widgets', 'ui_topology', 'ui_locales']
webservices_mods = {}

config_filename	= os.path.expanduser('~/etc/webserver.conf')
config		= ConfigParser.RawConfigParser()
config.read(config_filename)

mongo_config_file = os.path.expanduser('~/etc/cstorage.conf')
mongo_config = ConfigParser.RawConfigParser()
mongo_config.read(mongo_config_file)

## default config
debug		= True

#mongo config
mongo_host = mongo_config.get("master", "host")
mongo_port = mongo_config.getint("master", "port")
mongo_db = mongo_config.get("master", "db")

session_cookie_expires	= 300
session_secret			= 'canopsis'
session_lock_dir		= os.path.expanduser('~/tmp/webcore_cache')
session_mongo_url		= 'mongodb://%s:%s/%s.beaker' % (mongo_host,mongo_port,mongo_db)
root_directory			= os.path.expanduser("~/var/www/")

try:
	## get config
	debug					= config.getboolean('server', "debug")
	root_directory			= os.path.expanduser(config.get('server', "root_directory"))

	session_cookie_expires	= config.getint('session', "cookie_expires")
	session_secret			= config.get('session', "secret")
	session_data_dir		= os.path.expanduser(config.get('session', "data_dir"))

except Exception, err:
	print "Error when reading '%s' (%s)" % (config_filename, err)

## Logger
logging_level=logging.INFO
if debug:
	logging_level=logging.DEBUG
	
logging.basicConfig(format=r"%(asctime)s [%(process)d] [%(name)s] [%(levelname)s] %(message)s", datefmt=r"%Y-%m-%d %H:%M:%S", level=logging_level)
logger 	= logging.getLogger("webserver")
	
bottle.debug(debug)

## load and unload webservices
def load_webservices():
	logger.info("Load webservices.")
	sys.path.append(os.path.expanduser("~/opt/webcore/libexec/"))
	for webservice in webservices:
		try:
			module = __import__(webservice)
			webservices_mods[webservice] = module
			logger.info(" + '%s' imported." % webservice)
			
			try:
				module.load()
			except AttributeError:
				pass
			except Exception, err:
				logger.error("Impossible to load '%s'. (%s)" % (webservice, err))
				
		except Exception, err:
			logger.error("Impossible to import '%s'. (%s)" % (webservice, err))

def unload_webservices():
	logger.info("Unload webservices.")
	for webservice in webservices_mods:
		module = webservices_mods[webservice]
		try:
			module.unload()
			logger.info(" + '%s' unloaded." % webservice)
		except AttributeError:
			pass
		#except Exception, err:
		#	logger.error("Impossible to unload '%s'. (%s)" % (webservice, err))

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
import gevent, signal, sys
stop_in_progress=False
def signal_handler():
	global stop_in_progress
	if not stop_in_progress:
		stop_in_progress=True
		logger.info("Receive signal to stop worker ...")
		unload_webservices()
		logger.info("Ready to stop.")
		sys.exit(0)

gevent.signal(signal.SIGTERM, signal_handler)
gevent.signal(signal.SIGINT, signal_handler)

## Bottle App
app = bottle.default_app()

load_webservices()

## Install plugins
auth_plugin = libexec.auth.checkAuthPlugin()

bottle.install(auth_plugin)

## Session system with beaker
session_opts = {
    'session.type': 'mongodb',
    'session.cookie_expires': session_cookie_expires,
    'session.url' : session_mongo_url,
    'session.auto': True,
#   'session.timeout': 300,
    'session.secret': session_secret,
    'session.lock_dir' : session_lock_dir,
}

## Basic Handler
@route('/:lang/static/:path#.+#',	skip=['checkAuthPlugin'])
@route('/static/:path#.+#',			skip=['checkAuthPlugin'])
def server_static(path, lang='en'):
	key = request.params.get('authkey', default=None)
	if key:
		autoLogin(key)

	return static_file(path, root=root_directory)

@route('/favicon.ico',skip=[auth_plugin])
def favicon():
	return

@route('/',					skip=['checkAuthPlugin'])
@route('/:key',				skip=['checkAuthPlugin'])
@route('/index.html',		skip=['checkAuthPlugin'])
@route('/:lang/',			skip=['checkAuthPlugin'])
@route('/:lang/:key',		skip=['checkAuthPlugin'])
@route('/:lang/index.html',	skip=['checkAuthPlugin'])
def index(key=None, lang='en'):
	uri_key = request.params.get('authkey', default=None)
	if not key and uri_key:
		key = uri_key

	if key:
		autoLogin(key)

	redirect('/%s/static/canopsis/index.html' % lang)

## Install session Middleware
app = SessionMiddleware(app, session_opts)
