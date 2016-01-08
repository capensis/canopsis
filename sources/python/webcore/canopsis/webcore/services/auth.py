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

from canopsis.common.ws import route
from bottle import redirect, response, HTTPError, urlencode

from hashlib import sha1
from time import time

from . import session
from . import rights


def check(mode='authkey', user=None, password=None):
    def _check_shadow(user, key):
        if user and user['shadowpasswd'].upper() == key.upper():
            return user

        return None

    def _check_plain(user, key):
        shadowpasswd = sha1(key).hexdigest()
        return _check_shadow(user, shadowpasswd)

    def _check_crypted(user, key):
        if user:
            shadowpasswd = user['shadowpasswd'].upper()
            ts = str(int(time() / 10) * 10)
            tmpKey = '{0}{1}'.format(shadowpasswd, ts)

            cryptedKey = sha1(tmpKey).hexdigest().upper()

            if cryptedKey == key.upper():
                return user

        return None

    def _check_authkey(user, key):
        manager = rights.get_manager()

        if not user:
            user = manager.user_storage.get_elements(
                query={
                    'crecord_type': 'user',
                    'authkey': key
                }
            )

        if user and user[0]['authkey'] == key:
            return user[0]

        else:
            return None

    handlers = {
        'plain': _check_plain,
        'shadow': _check_shadow,
        'crypted': _check_crypted,
        'authkey': _check_authkey
    }

    if mode in handlers:
        return handlers[mode](user, password)

    return None


def autoLogin(key):
    user = check(mode='authkey', password=key)

    if not user:
        return HTTPError(403, 'Forbidden')

    if not user.get('enable', False):
        return HTTPError(403, 'Account is disabled')

    session.create(user)

    return user


def exports(ws):
    session = ws.require('session')
    rights = ws.require('rights').get_manager()

    @route(
        ws.application.post,
        name='auth',
        wsgi_params={
            'skip': ws.skip_login
        },
        payload=[
            'username', 'password',
            'shadow', 'crypted'
        ],
        response=lambda data, adapt: data,
        nolog=True
    )
    def auth_route(
        username=None, password=None,
        shadow=False, crypted=False
    ):
        ws.logger.info('/auth')

        if not username or not password:
            redirect('/?logerror=1')

        mode = 'plain'

        if shadow:
            mode = 'shadow'

        if crypted:
            mode = 'crypted'

        # Try to find user in database
        user = rights.get_user(username)

        # No such user, or it's an external one
        if not user or user.get('external', False):
            # Try to redirect authentication to the external backend
            if mode == 'plain':
                response.status = 307
                response.set_header('Location', '/auth/external')

                return 'username={0}&password={1}'.format(
                    urlencode(username),
                    urlencode(password)
                )

            else:
                #return HTTPError(403, 'Plain authentication required')
                redirect('/?logerror=3')

        # Local authentication: check if account is activated
        if not user.get('enable', False):
            #return HTTPError(403, 'This account is not enabled')
            redirect('/?logerror=2')

        user = check(mode=mode, user=user, password=password)

        if not user:
            redirect('/?logerror=1')
            #return HTTPError(403, 'Forbidden')

        session.create(user)
        redirect('/')

    @route(ws.application.post, name='auth/external', nolog=True)
    def auth_external(**kwargs):
        ws.logger.info('/auth/external')

        # When we arrive here, the Bottle plugins in charge of authentication
        # have initialized the session, we just need to redirect to the index.
        redirect('/static/canopsis/index.html')

    @ws.application.get('/logged_in')
    def logged_in():
        # Route used when came back from auth backend
        redirect('/static/canopsis/index.html')

    @route(ws.application.get, wsgi_params={'skip': ws.skip_login}, nolog=True)
    def autologin(key):
        return autoLogin(key)

    @route(ws.application.get, wsgi_params={'skip': ws.skip_logout})
    def logout():
        session.delete()
        redirect('/')
