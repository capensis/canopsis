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

import bottle
import logging
from bottle import get, request, response, post, HTTPError, redirect

## Canopsis
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

logger = logging.getLogger("auth")
logger.setLevel(logging.INFO)

#session variable
session_accounts = {
    'anonymous': Account()
}

# List of plugins to skip
# CAS will handle the logout method
auth_backends = ['AuthKeyBackend', 'LDAPBackend', 'CASBackend', 'EnsureAuthenticated']
auth_backends_logout = ['AuthKeyBackend', 'LDAPBackend', 'EnsureAuthenticated']


class EnsureAuthenticated(object):
    name = 'EnsureAuthenticated'

    def apply(self, callback, context):
        def decorated(*args, **kwargs):
            s = request.environ.get('beaker.session')

            if not s.get('auth_on', False):
                return HTTPError(401, 'Not authorized')

            return callback(*args, **kwargs)

        return decorated


@post('/auth', skip=auth_backends)
def auth():
    login = request.params.get('username', default=None)
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
        account = Account(storage.get(_id, account=Account(user=login)))
    except Exception as err:
        logger.error(err)

    ## Check
    if not account or account.external:
        if mode == 'plain':
            response.status = 307
            response.set_header('Location', '/auth/external')

            return 'username={0}&password={1}'.format(login, password)

        else:
            return HTTPError(403, 'Plain authentication required')

    if not account.is_enable():
        return HTTPError(403, "This account is not enabled")

    logger.debug(" + Check with local db")
    access = False

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

    redirect('/')


@post('/auth/external')
def doauth():
    # When we arrive here, the Bottle plugins in charge of authentication have
    # initialized the session, we just need to redirect to the index.
    bottle.redirect('/static/canopsis/index.html')


@get('/logged_in')
def logged_in():
    # Route used when came back from CAS or any other external backend
    bottle.redirect('/static/canopsis/index.html')


@get('/autoLogin/:key', skip=auth_backends)
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

    return {'total': 1, 'data': [account.dump()], 'success': True}


#Access for disconnect and clean session
@get('/logout', skip=auth_backends_logout)
@get('/disconnect', skip=auth_backends_logout)
def disconnect():
    s = bottle.request.environ.get('beaker.session')
    user = s.get('account_user', None)

    if not user:
        return HTTPError(403, "Forbidden")

    logger.debug("Disconnect '%s'" % user)
    delete_session()

    bottle.redirect('/')
    # return {'total': 0, 'success': True, 'data': []}


#find the account in memory, or try to find it from database, if not in db log anon
def get_account(_id=None):
    logger.debug("Get Account:")

    s = bottle.request.environ.get('beaker.session')

    if not _id:
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
        storage = get_storage(namespace='object')
        account = Account(storage.get(_id, account=Account(user='root')))
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

    account = Account(record)

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

    record = storage.find_one(mfilter=mfilter, account=Account(user='root'))

    if not record:
        return None

    logger.debug(" + Done")
    key_account = Account(record)

    if account and account._id != key_account._id:
        logger.debug(" + Account missmatch")
        return None

    return key_account


def check_root(account):
    if account._id == 'account.root' or account.group == 'group.CPS_root':
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

    logger.debug("'%s' is not in '%s'" % (account.user, group_id))
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
        ids = [_id]

    logger.debug("Delete session '%s'" % ids)

    for _id in ids:
        try:
            del session_accounts[_id]
        except:
            pass

    if account._id in ids:
        s = bottle.request.environ.get('beaker.session')
        s.delete()
