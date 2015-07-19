#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
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

import sys
sys.path.append("/opt/canopsis/opt/webcore/")

import logging
import traceback

logging.basicConfig(
	format=r"%(asctime)s [%(process)d] [%(name)s] [%(levelname)s] %(message)s",
	datefmt=r"%Y-%m-%d %H:%M:%S",
	level=logging.DEBUG
)

import wsgi_webserver
from webtest import TestApp

from canopsis.old.account import Account
from canopsis.old.storage import get_storage

storage = get_storage(namespace='object')

user = 'root'
pwd = 'root'
shadow = Account().make_shadow('root')
crypted = Account().make_tmp_cryptedKey(shadow=shadow)
authkey = storage.get('account.%s' % user, account=Account(user='root')).data['authkey']

app = TestApp(wsgi_webserver.app)


def quit(code=0):
	wsgi_webserver.unload_webservices()
	sys.exit(code)


def get(uri, args={}, status=None, params=None):
	print("Get %s" % uri)
	resp = app.get(uri, status=status, params=params)
	print(" + code: %s" % resp.status_int)

	return resp


def test(func):
	def wrapper(*args, **kwargs):
		try:
			print("### %s" % func.__name__)
			func(*args, **kwargs)
			return True
		except Exception:
			traceback.print_exc(file=sys.stdout)
			quit(1)
			return False

	return wrapper


@test
def logout():
	get('/logout')


@test
def login_failed():
	get('/auth/toto/dummy', status=[403])
	get('/auth/%s/dummy' % user, status=[403])


@test
def login_plain():
	get('/auth/%s/%s' % (user, pwd), status=[200])
	logout()


@test
def login_plain_ldap():
	storage.remove('account.toto', account=Account(user='root'))
	get('/auth/toto/aqzsedrftg123;', status=[200])
	logout()

	get('/auth/toto/tata', status=[403])
	get('/auth/toto/aqzsedrftg123;', status=[200])
	logout()


@test
def login_shadow():
	params = {'shadow': 1}
	get('/auth/%s/%s' % (user, shadow), params=params, status=[200])
	logout()


@test
def login_crypted():
	params = {'crypted': 1}
	get('/auth/%s/%s' % (user, crypted), params=params, status=[200])
	logout()


@test
def login_authkey():
	get('/autoLogin/dummy', status=[403])
	get('/autoLogin/%s' % authkey, status=[200])
	logout()

	params = {'authkey': authkey}
	get('/autoLogin/dummy', status=[403])
	get('/account/me', params=params, status=[200])
	logout()


@test
def login_checkAuthPlugin():
	get('/account/checkAuthPlugin1', status=[403])
	#resp = get('/canopsis/auth.html', status=[200])

	# Login
	get('/autoLogin/%s' % authkey, status=[200])

	get('/account/checkAuthPlugin1', status=[200])
	logout()

	get('/auth/canopsis/canopsis', status=[200])
	get('/ui/view', status=[200])
	get('/account/checkAuthPlugin1', status=[200])
	get('/account/checkAuthPlugin2', status=[403])

	logout()

## Execute test
#login_failed()
#login_plain()
login_plain_ldap()
#login_shadow()
#login_crypted()
#login_authkey()
#login_checkAuthPlugin()


quit()
