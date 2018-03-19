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

from urllib import quote_plus

from bottle import redirect, response, HTTPError

from canopsis.common.ws import route
from canopsis.webcore.services import session as session_module
from canopsis.webcore.services import rights as rights_module
from canopsis.auth.check import check


def autoLogin(key):
    user = check(mode='authkey', password=key)

    if not user:
        return HTTPError(403, 'Forbidden')

    if not user.get('enable', False):
        return HTTPError(403, 'Account is disabled')

    session_module.create(user)

    return user


def exports(ws):
    session = session_module
    rights = rights_module.get_manager()

    @route(
        ws.application.post,
        name='auth',
        wsgi_params={'skip': ws.skip_login},
        payload=[
            'username', 'password',
            'shadow', 'crypted'
        ],
        response=lambda data, adapt: data
    )
    def auth_route(
        username=None, password=None, shadow=False, crypted=False
    ):
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
                # canopsis only use the default auth backend
                if ws.auth_backends.keys() == ['AuthKeyBackend', u'EnsureAuthenticated']:
                    location = '/auth/internal'
                else:
                    location = '/auth/external'
                response.set_header('Location', location)

                return 'username={0}&password={1}'.format(
                    quote_plus(username),
                    quote_plus(password))

            else:
                # return HTTPError(403, 'Plain authentication required')
                redirect('/?logerror=3')

        # Local authentication: check if account is activated
        if not user.get('enable', False):
            # return HTTPError(403, 'This account is not enabled')
            redirect('/?logerror=2')

        user = check(mode=mode, user=user, password=password)

        if not user:
            redirect('/?logerror=1')
            # return HTTPError(403, 'Forbidden')

        session.create(user)
        redirect('/')

    @route(ws.application.post,
           name='auth/internal',
           wsgi_params={'skip': ws.skip_login})
    def auth_internal(**kwargs):
        # When we arrive here, the Bottle plugins in charge of authentication
        # have initialized the session, we just need to redirect to the index.
        redirect('/?logerror=1')

    @route(ws.application.post, name='auth/external')
    def auth_external(**kwargs):
        # When we arrive here, the Bottle plugins in charge of authentication
        # have initialized the session, we just need to redirect to the index.
        redirect('/?logerror=1')

    @ws.application.get('/logged_in')
    def logged_in():
        # Route used when came back from auth backend
        redirect('/static/canopsis/index.html')

    @route(ws.application.get, wsgi_params={'skip': ws.skip_login})
    def autologin(authkey=''):
        return autoLogin(authkey)

    @route(ws.application.get, wsgi_params={'skip': ws.skip_logout})
    def logout():
        session.delete()
        redirect('/')
