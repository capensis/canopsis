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


from bottle import request, static_file, redirect, response, HTTPError

from canopsis.common.ws import route
from canopsis.webcore.services import session as session_module
from canopsis.webcore.services import rights as rights_module
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_FORBIDDEN
from canopsis.auth.check import check
from canopsis.userinterface.manager import UserInterfaceManager
from canopsis.version import CanopsisVersionManager
from canopsis.common.mongo_store import MongoStore

# The USER_FIELDS list contains the names of the keys of the user dictionary
# that are returned in JSON when the authentication is successful.
USER_FIELDS = ['authkey', 'contact', 'crecord_name', 'mail', 'role']


def autoLogin(key):
    user = check(mode='authkey', password=key)

    if not user:
        return HTTPError(HTTP_FORBIDDEN, 'Forbidden')

    if not user.get('enable', False):
        return HTTPError(HTTP_FORBIDDEN, 'Account is disabled')

    session_module.create(user)

    return user


def json_auth_success(user):
    """
    Return a JSON object containing informations about a user.

    This is used as the response to a successful authentication request.

    :param Dict[str, Any] user: The user who was suceesfully authenticated.
    :rtype: str
    """
    response_body = {}

    for field in USER_FIELDS:
        try:
            response_body[field] = user[field]
        except KeyError:
            pass

    return gen_json(response_body)


def exports(ws):
    session = session_module
    rights = rights_module.get_manager()

    @route(
        ws.application.post,
        name='auth',
        wsgi_params={'skip': ws.skip_login},
        payload=[
            'username', 'password',
            'shadow', 'crypted', 'json_response'
        ],
        response=lambda data, adapt: data
    )
    def auth_route(
        username=None, password=None, shadow=False, crypted=False,
        json_response=False
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

                response_body = 'username={0}&password={1}'.format(
                    quote_plus(username),
                    quote_plus(password))

                if json_response:
                    response_body += '&json_response=True'

                return response_body

            else:
                if json_response:
                    return gen_json_error({
                        'description': 'Plain authentication required'
                    }, HTTP_FORBIDDEN)
                else:
                    redirect('/?logerror=3')

        # Local authentication: check if account is activated
        if not user.get('enable', False):
            if json_response:
                return gen_json_error({
                    'description': 'Account disabled'
                }, HTTP_FORBIDDEN)
            else:
                redirect('/?logerror=2')

        user = check(mode=mode, user=user, password=password)

        if not user:
            if json_response:
                return gen_json_error({
                    'description': 'Wrong login or password'
                }, HTTP_FORBIDDEN)
            else:
                redirect('/?logerror=1')

        session.create(user)
        if json_response:
            return json_auth_success(user)
        else:
            redirect('/')

    @route(ws.application.post,
           name='auth/internal',
           wsgi_params={'skip': ws.skip_login},
           response=lambda data, adapt: data)
    def auth_internal(json_response=False, **kwargs):
        # When we arrive here, the Bottle plugins in charge of authentication
        # should have initialized the session.
        if json_response:
            return gen_json_error({
                'description': 'Wrong login or password'
            }, HTTP_FORBIDDEN)
        else:
            # Redirect to /?logerror=1. If the session was initialized, the
            # user will be redirected to the index of canopsis. If not, an
            # error message will be shown.
            redirect('/?logerror=1')

    @route(ws.application.post,
           name='auth/external',
           response=lambda data, adapt: data)
    def auth_external(json_response=False, **kwargs):
        # When we arrive here, the Bottle plugins in charge of authentication
        # should have initialized the session.
        if json_response:
            # Check if a session was indeed initialized by one of the external
            # authentication backends.
            s = session_module.get()

            if s.get('auth_on', False):
                user = session_module.get_user()
                return json_auth_success(user)
            else:
                return gen_json_error({
                    'description': 'Wrong login or password'
                }, HTTP_FORBIDDEN)
        else:
            # Redirect to /?logerror=1. If the session was initialized, the
            # user will be redirected to the index of canopsis. If not, an
            # error message will be shown.
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

    @ws.application.get('/api/internal/login/login_info')
    def get_internal_login_info():
        login_info = {}

        cservices = {
            'webserver': {provider: 1 for provider in ws.providers},
        }

        records = ws.db.find(
            {'crecord_name': {'$in': ['casconfig', 'ldapconfig']}},
            namespace='object'
        )

        for cservice in records:
            cservice = cservice.dump()
            cname = cservice['crecord_name']
            cservices[cname] = cservice

            ws.logger.info(u'found cservices type {}'.format(cname))

            if cname == 'casconfig':
                cservice['server'] = cservice['server'].rstrip('/')
                cservice['service'] = cservice['service'].rstrip('/')
                ws.logger.info(u'cas config : server {}, service {}'.format(
                    cservice['server'],
                    cservice['service'],
                ))

        if "canopsis_cat.webcore.services.saml2" in ws.webmodules:
            result = ws.db.find({'_id': "canopsis"}, namespace='default_saml2')

            cservices["saml2config"] = {
                "url": result[0].data["saml2"]["settings"]["idp"]["singleSignOnService"]["url"]}

        login_info["auth"] = cservices

        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        user_interface = user_interface_manager.get()

        if user_interface is not None:
            login_info["user_interface"] = user_interface.to_dict()

        vStore = MongoStore.get_default()
        vCollection = \
            vStore.get_collection(name=CanopsisVersionManager.COLLECTION)
        vDocument = CanopsisVersionManager(vCollection).\
            find_canopsis_version_document()

        login_info[CanopsisVersionManager.VERSION_FIELD] = vDocument[CanopsisVersionManager.VERSION_FIELD]

        return gen_json(login_info)

    @ws.application.post('/api/internal/login/login_info')
    def update_internal_login_info():

        user_interface_manager = UserInterfaceManager(
            *UserInterfaceManager.provide_default_basics())

        return gen_json(user_interface_manager.get().to_dict())
