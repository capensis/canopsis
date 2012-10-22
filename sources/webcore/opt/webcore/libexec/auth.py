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

#session variable
session_accounts = {}

#########################################################################
class checkAuthPlugin(object):
	name='checkAuthPlugin'
	def __init__(self,authorized_grp=[],unauthorized_grp=[]):
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
		authorized_grp = None
		unauthorized_grp = None
		
		if 'config' in context:
			conf = context['config'].get('checkAuthPlugin',{})
		if 'authorized_grp' in conf:
			authorized_grp = conf.get('authorized_grp',None)
		if 'unauthorized_grp' in conf:
			unauthorized_grp = conf.get('unauthorized_grp',None)

		def do_auth(*args, **kawrgs):
			try:
				path = bottle.request.path
			except:
				path = None
				
			url = bottle.request.url
			s = bottle.request.environ.get('beaker.session')

			if s.get('auth_on',False):
				logger.debug(' + Session already open')
				
				account_id = s['account_id']
				account_group = s['account_group']
				account_groups = s['account_groups']
				access = True
			
				if not (account_id == 'account.root' or  account_group == 'group.CPS_root' or 'group.CPS_root' in  account_groups):
					if unauthorized_grp:
						if account_group in unauthorized_grp:
							access = False
						for group in account_groups:
							if group in unauthorized_grp:
								access = False
								
					if authorized_grp:
						access = False
						if account_group in authorized_grp:
							access = True
						for group in account_groups:
							if group in authorized_grp:
								access = True
						
				logger.debug(" + Authentified, Session is Ok.")
			else:
				logger.debug(' + Session not open for this user')
				if path == "/canopsis/auth.html":
					access = True
				else:
					access = False

			if access:
				logger.debug(" + Valid auth")
				return callback(*args, **kawrgs)
			else:
				logger.error(" + Invalid auth")
				return HTTPError(403, 'Insufficient rights')
				#return {'total': 0, 'success': False, 'data': []}
				#return redirect('/static/canopsis/auth.html' + '?url=' + url)
		return do_auth
#########################################################################



@get('/auth/:login/:password',skip=['checkAuthPlugin'])
@get('/auth/:login',skip=['checkAuthPlugin'])
@get('/auth',skip=['checkAuthPlugin'])
def auth(login=None, password=None):
	if not login:
		login = request.params.get('login', default=None)

	shadow = request.params.get('shadow', default=False)
	if shadow:
		shadow = True
	
	cryptedKey = request.params.get('cryptedKey', default=False)

	if not password:
		password = request.params.get('password', default=None)

	if not login or not password:
		return HTTPError(404, "Invalid arguments")

	_id = "account." + login

	logger.debug(" + _id: "+_id)
	logger.debug(" + Login: "+login)
	logger.debug(" + Password: "+password)
	logger.debug("    + is Shadow: "+str(shadow))
	logger.debug("    + is cryptedKey: "+str(cryptedKey))

	storage = get_storage(namespace='object')

	try:
		account = caccount(storage.get(_id, account=caccount(user=login)))
		logger.debug(" + Check password ...")

		if shadow:
			access = account.check_shadowpasswd(password)
			
		elif cryptedKey:
			logger.debug(" + Valid auth key: %s" % (account.make_tmp_cryptedKey()))
			access = account.check_tmp_cryptedKey(password)
			
		else:
			access = account.check_passwd(password)

		if access:
			s = bottle.request.environ.get('beaker.session')
			s['account_id'] = account._id
			s['account_user'] = account.user
			s['account_group'] = account.group
			s['account_groups'] = account.groups
			s['auth_on'] = True
			s.save()

			output = [ account.dump() ]
			output = {'total': len(output), 'success': True, 'data': output}
			return output
		else:
			logger.debug(" + Invalid password ...")
	except Exception, err:
		logger.error(err)
		
	return HTTPError(403, "Forbidden")

@get('/autoLogin/:key',skip=['checkAuthPlugin'])
def autoLogin(key=None):
	if not key:
		return HTTPError(404, "No key provided")
	#---------------------Get storage/account-------------------
	storage = get_storage(namespace='object')
	
	mfilter = {
				'crecord_type':'account',
				'authkey':key,
			}
				
	foundByKey = storage.find(mfilter=mfilter, account=caccount(user='root'))
	#-------------------------if found, create session and redirect------------------------
	if len(foundByKey) == 1:
		account = caccount(foundByKey[0])
		s = bottle.request.environ.get('beaker.session')
		s['account_id'] = account._id
		s['account_user'] = account.user
		s['account_group'] = account.group
		s['account_groups'] = account.groups
		s['auth_on'] = True
		s.save()
		#logger.debug('Autologin success, redirecting browser')
		#redirect('/static/canopsis/index.html')
		account = caccount(foundByKey[0])
		return {'total':1,'data':[account.dump()],'success':True}
	else:
		logger.debug('Autologin failed, no key match the provided one')
		return {'total':0,'data':{},'success':False}

@get('/keyAuth/:login/:key',skip=['checkAuthPlugin'])
@get('/keyAuth/:login',skip=['checkAuthPlugin'])
@get('/keyAuth',skip=['checkAuthPlugin'])
def keyAuth(login=None, key=None):
	#-----------------------Get info------------------------
	if not login:
		login = request.params.get('login', default=None)
	if not key:
		key = request.params.get('key', default=None)
		
	if not login or not key:
		return HTTPError(404, "Invalid arguments")
	
	#---------------------Get storage/account-------------------
	_id = "account.%s" % login
	
	logger.debug(" + _id: %s" % _id)
	logger.debug(" + Login: %s" % login)
	logger.debug(" + key: %s" % key)
	
	storage = get_storage(namespace='object')

	try:
		account = caccount(storage.get(_id, account=caccount(user=login)))
	except Exception, err:
		logger.error('Error while fetching %s : %s' % (_id,err))
		return HTTPError(403, "There is no account for this login")
	
	#---------------------Check key-------------------------
	if account.check_authkey(key):
		s = bottle.request.environ.get('beaker.session')
		s['account_id'] = account._id
		s['account_user'] = account.user
		s['account_group'] = account.group
		s['account_groups'] = account.groups
		s['auth_on'] = True
		s.save()
		
		logger.debug('Access granted and session open for %s' % _id)
		
		output = [account.dump()]
		return {'total': len(output), 'success': True, 'data': output}
	else:
		logger.error('Wrong key given for %s' % _id)
		return HTTPError(403, "Wrong key, access prohibited")
	
	


#Access for disconnect and clean session
@get('/logout')
@get('/disconnect')
def disconnect():
	logger.error("Disconnect")
	s = bottle.request.environ.get('beaker.session')
	s.delete()
	return {'total': 0, 'success': True, 'data': []}


#decorator in order to protect request
def check_auth(callback):
	def do_auth(*args, **kawrgs):
		try:
			path = kawrgs['path']
		except:
			path = None
	
		url = bottle.request.url
		#get beaker session and test it right after
		s = bottle.request.environ.get('beaker.session')
		#add caccount to parameters
		if s.get('auth_on',False) or path == "canopsis/auth.html":
			logger.debug("Session is Ok.")
			return callback(*args, **kawrgs)

		logger.error("Invalid auth")
		return {'total': 0, 'success': False, 'data': []}
		#return redirect('/static/canopsis/auth.html' + '?url=' + url)

	return do_auth
	
#find the account in memory, or try to find it from database, if not in db log anon
def get_account(_id=None):
	logger.debug("Get Account:")
	if not _id:
		s = bottle.request.environ.get('beaker.session')
		_id = s.get('account_id',0)
		logger.debug(" + Get _id from Beaker Session (%s)" % _id)

	user = s.get('account_user',0)

	logger.debug(" + Try to load account %s ('%s') ..." % (user, _id))

	storage = get_storage(namespace='object')

	try:
		account = session_accounts[_id]
		logger.debug(" + Load account from memory.")
	except:
		if _id:
			record = storage.get(_id, account=caccount(user=user) )
			logger.debug(" + Load account from DB.")
			account = caccount(record)
			session_accounts[_id] = account
		else:
			logger.debug(" + Impossible to load account, return Anonymous account.")
			try:
				return session_accounts['anonymous']
			except:
				session_accounts['anonymous'] = caccount()
				return session_accounts['anonymous']

	return account

#cache is cool, but when you change rights, cache still have old rights, so reload 
def reload_account(_id=None,record=None):
	try:
		logger.debug('Reload Account %s' % _id)
		account = get_account()
		storage = get_storage(namespace='object')
		account_to_update = None
		
		if not record:
			if _id:	
				record = storage.get(_id, account=account)
				account_to_update = caccount(record)
			else:
				record = storage.get(account._id, account=account )
				account_to_update = caccount(record)
				
		if not account_to_update:
			account_to_update = caccount(record)
			
		session_accounts[account_to_update._id] = account_to_update
		s = bottle.request.environ.get('beaker.session')
		s['account_group'] = account_to_update.group
		s['account_groups'] = account_to_update.groups
		logger.debug('Account %s is in following groups : %s' % (_id,str(account_to_update.groups)))
		s.save()
		return True
	except Exception,err:
		logger.error('Account reloading failed : %s' % err)
		return False

def check_group_rights(account,group_id):
	if not (account._id == 'account.root' or  account.group == 'group.CPS_root' or 'group.CPS_root' in account.groups):
		if not group_id in account.groups and group_id != account.group:
			logger.debug('%s is not in %s' % (account.user,group_id))
			return False
	return True
