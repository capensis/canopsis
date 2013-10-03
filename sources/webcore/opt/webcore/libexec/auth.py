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

import sys, os, logging, json
import bottle, logging, hashlib, json
from bottle import error, route, get, request, post, HTTPError, redirect

from beaker.middleware import SessionMiddleware

## Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord

logger = logging.getLogger("auth")
logger.setLevel(logging.DEBUG)

#session variable
session_accounts = {
	'anonymous': caccount()
}

#########################################################################
class checkAuthPlugin(object):
	name='checkAuthPlugin'

	def __init__(self,authorized_grp=[], unauthorized_grp=[]):
		self.authorized_grp = authorized_grp
		self.unauthorized_grp = unauthorized_grp
		self.keyword = None
		
	def setup(self, app):
		for other in app.plugins:
			if not isinstance(other, checkAuthPlugin): continue
			if other.keyword == self.keyword:
				raise PluginError("Found another auth plugin with "\
				"conflicting settings (non-unique keyword).")
	
	def apply(self, callback, context):
		authorized_grp = []
		unauthorized_grp = []
		
		if 'config' in context:
			conf = context['config'].get('checkAuthPlugin',{})
		if 'authorized_grp' in conf:
			authorized_grp = conf.get('authorized_grp', [])
		if 'unauthorized_grp' in conf:
			unauthorized_grp = conf.get('unauthorized_grp', [])

		if not isinstance(authorized_grp, list):
			authorized_grp = [authorized_grp]

		if not isinstance(unauthorized_grp, list):
			unauthorized_grp = [unauthorized_grp]

		def do_auth(*args, **kawrgs):
			logger.debug("Auth query")
			access = False

			try:
				path = bottle.request.path
			except:
				path = None
				
			url = bottle.request.url

			account = get_account()

			authkey = request.params.get('authkey', None)
			if authkey:
				logger.debug("Check authkey: '%s'" % authkey)

				account = check_authkey(authkey)

				if account:
					create_session(account)
					access = True
				else:
					logger.debug(" + Failed")
			
			if not account or account.user == 'anonymous':
				return HTTPError(403, "Forbidden")

			logger.debug(" + _id: %s" % account._id)

			if check_root(account):
				access = True
			else:
				if not authorized_grp and not unauthorized_grp:
					access = True
				else:
					logger.debug("Check authorized_grp: %s" % authorized_grp)
					if authorized_grp:
						for group in authorized_grp:
							if check_group_rights(account, group):
								access = True
								break

					logger.debug("Check unauthorized_grp and overwrite access: %s" % unauthorized_grp)
					if unauthorized_grp:
						for group in unauthorized_grp:
							if check_group_rights(account, group):
								access = False
								break


			#logger.debug("Check path: '%s'" % path)
			#if path == "/canopsis/auth.html":
			#	access = True

			logger.debug(" + Access: %s" % access)

			if access:
				logger.debug(" + Auth ok")
				return callback(*args, **kawrgs)
			else:
				logger.error(" + Invalid auth")
				return HTTPError(403, 'Insufficient rights')

		return do_auth
#########################################################################

# Test checkAuthPlugin
@get('/account/checkAuthPlugin1', checkAuthPlugin={'authorized_grp': 'group.Canopsis'})
def account_checkAuthPlugin():
	logger.info("Access granted by checkAuthPlugin")
	return {'total': 0, 'success': True, 'data': []}

@get('/account/checkAuthPlugin2', checkAuthPlugin={'unauthorized_grp': 'group.Canopsis'})
def account_checkAuthPlugin():
	logger.info("Access granted by checkAuthPlugin")
	return {'total': 0, 'success': True, 'data': []}

@get('/auth/:login/:password',skip=['checkAuthPlugin'])
@get('/auth/:login',skip=['checkAuthPlugin'])
@get('/auth',skip=['checkAuthPlugin'])
def auth(login=None, password=None):
	if not login:
		login = request.params.get('login', default=None)

	if not password:
		password = request.params.get('password', default=None)

	if not login or not password:
		return HTTPError(400, "Invalid arguments")

	mode = 'plain'

	if request.params.get('shadow', default=False):
		mode = 'shadow'
	
	if request.params.get('cryptedKey', default=False) or request.params.get('crypted', default=False):
		mode = 'crypted'

	_id = "account.%s" % login

	logger.debug(" + _id:      %s" % _id)
	logger.debug(" + Login:    %s" % login)
	logger.debug(" + Mode:     %s" % mode)

	storage = get_storage(namespace='object')
	account = None

	## Local
	try:
		account = caccount(storage.get(_id, account=caccount(user=login)))
	except Exception, err:
		logger.error(err)

	## External
	# Try to provisionning account
	if not account and mode == 'plain':
		try:
			account = external_prov(login, password)
		except Exception, err:
			logger.error(err)

	## Check
	if not account:
		return HTTPError(403, "Forbidden")

	logger.debug(" + Check password ...")

	if not account.is_enable():
		return HTTPError(403, "This account is not enabled")

	if account.external and  mode != 'plain':
		return HTTPError(403, "Send your password in plain text")

	access = None

	if account.external and mode == 'plain':
		access = external_auth(login, password)

	if access == None:
		logger.debug(" + Check with local db")

		if mode == 'plain':
			access = account.check_passwd(password)

		elif mode == 'shadow':
			access = account.check_shadowpasswd(password)

		elif mode == 'crypted':
			access = account.check_tmp_cryptedKey(password)

	if not access:
		logger.debug(" + Invalid password ...")
		return HTTPError(403, "Forbidden")

	create_session(account)

	output = [ account.dump() ]
	output = {'total': len(output), 'success': True, 'data': output}
	return output

@get('/autoLogin/:key',skip=['checkAuthPlugin'])
def autoLogin(key=None):
	if not key:
		return HTTPError(400, "No key provided")

	account = check_authkey(key)

	if not account:
		return HTTPError(403, "Forbidden")

	if not account.is_enable():
		logger.debug(" + Account is disabled")
		return HTTPError(403, "Account is disabled")

	create_session(account)

	return {'total':1,'data':[ account.dump() ], 'success':True}
	
#Access for disconnect and clean session
@get('/logout',		skip=['checkAuthPlugin'])
@get('/disconnect', skip=['checkAuthPlugin'])
def disconnect():
	s = bottle.request.environ.get('beaker.session')
	user = s.get('account_user', None)

	if not user:
		return HTTPError(403, "Forbidden")

	logger.debug("Disconnect '%s'" % user)
	delete_session()
	return {'total': 0, 'success': True, 'data': []}

#find the account in memory, or try to find it from database, if not in db log anon
def get_account(_id=None):
	logger.debug("Get Account:")

	if not _id:
		s = bottle.request.environ.get('beaker.session')
		_id = s.get('account_id', None)

	logger.debug(" + _id: %s" % _id)

	if not _id:
		logger.debug("  + Failed")
		return session_accounts['anonymous']

	logger.debug("Load account from cache")
	account = session_accounts.get(_id, None)

	if account:
		return account

	logger.debug("  + Failed")

	logger.debug("Load account from DB")
	
	try:
		user = s.get('account_user', None)	
		storage = get_storage(namespace='object')
		account = caccount(storage.get(_id, account=caccount(user=user)))
	except Exception as err:
		logger.debug("  + Failed: %s" % err)
		return session_accounts['anonymous']

	session_accounts[_id] = account
	return account

#cache is cool, but when you change rights, cache still have old rights, so reload 
def reload_account(_id):
	logger.info('Reload cache and session')

	try:
		account = get_account()
		storage = get_storage(namespace='object')
		record = storage.get(_id, account=account)
	except Exception as err:
		logger.error('Impossible to load record: %s' % err)
		return False

	account = caccount(record)

	if not account:
		logger.error('Impossible to load account')
		return False

	session_accounts[account._id] = account

	logger.debug(' + Session and cache are reloaded')
	return True

def check_authkey(key, account=None):
	storage = get_storage(namespace='object')
	
	mfilter = {
		'crecord_type': 'account',
		'authkey': key,
	}
	
	logger.debug("Search authkey: '%s'" % key)

	record = storage.find_one(mfilter=mfilter, account=caccount(user='root'))

	if not record:
		return None

	logger.debug(" + Done")
	key_account = caccount(record)

	if account and account._id != key_account._id:
		logger.debug(" + Account missmatch")
		return None

	return key_account

def check_root(account):
	if account._id == 'account.root' or  account.group == 'group.CPS_root':
		return True

	if 'group.CPS_root' in account.groups:
		return True

	return False

def check_group_rights(account, group_id):
	if check_root(account):
		return True

	if group_id == account.group:
		return True

	if group_id in account.groups:
		return True

	logger.debug("'%s' is not in '%s'" % (account.user,group_id))
	return False

def create_session(account):
	session_accounts[account._id] = account

	s = bottle.request.environ.get('beaker.session')
	s['account_id'] = account._id
	s['account_user'] = account.user
	s['account_group'] = account.group
	s['account_groups'] = account.groups
	s['auth_on'] = True
	s.save()

	return s

def delete_session(_id=None):
	account = get_account()

	if not _id:	
		_id = account._id

	if isinstance(_id, list):
		ids = _id
	else:
		ids = [ _id ]

	logger.debug("Delete session '%s'" % ids)

	for _id in ids:
		try:
			del session_accounts[_id]
		except:
			pass

	if account._id in ids:
		s = bottle.request.environ.get('beaker.session')
		s.delete()

def external_prov(login, password):
	import auth_ldap
	
	logger.debug(" + Check External provisionning")
	return auth_ldap.prov(login, password)

def external_auth(login=None, password=None):
	import auth_ldap

	logger.debug(" + Check External Auth")
	return auth_ldap.auth(login, password)	
