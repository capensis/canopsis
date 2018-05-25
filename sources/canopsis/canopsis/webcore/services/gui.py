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

from __future__ import unicode_literals
from bottle import static_file, request, redirect, response
import os

from canopsis.common.template import Template

from canopsis.webcore.services import auth as auth_module
from canopsis.webcore.services import session as session_module


def exports(ws):
    skip_login = ws.skip_login

    @ws.application.get('/:lang/static/canopsis/index.html', skip=skip_login)
    @ws.application.get('/static/canopsis/index.html', skip=skip_login)
    def index(lang='en'):
        # Redirect user if not logged in
        if not session_module.get_user():
            redirect('/')

        return static_file('canopsis/index.html', root=ws.root_directory)

    @ws.application.get('/:lang/static/<filename:path>', skip=skip_login)
    @ws.application.get('/static/<filename:path>', skip=skip_login)
    def server_static(filename, lang='en'):
        key = request.params.get('authkey', default=None)

        if key:
            auth_module.autoLogin(key)

        if 'listalarm' in filename or 'timeline' in filename:
            response.set_header("Cache-Control", "public, no-cache")

        return static_file(filename, root=ws.root_directory)

    @ws.application.get('/favicon.ico', skip=skip_login)
    def favicon():
        return

    @ws.application.get('/', skip=skip_login)
    @ws.application.get('/index.html', skip=skip_login)
    @ws.application.get('/:lang/', skip=skip_login)
    @ws.application.get('/:lang/index.html', skip=skip_login)
    @ws.application.get('/:lang/:key/', skip=skip_login)
    @ws.application.get('/:lang/:key/index.html', skip=skip_login)
    def loginpage(lang='en', key=None):
        session = request.environ.get('beaker.session')

        # Try to authenticate user
        key = key or request.params.get('authkey', default=None)
        logerror = request.params.get('logerror', default=None)

        if logerror in [None, '1', '2', '3']:
            logmessage = {
                None: '',
                '1': 'Wrong login or password',
                '2': 'Account disabled',
                '3': 'Plain authentication required'
            }[logerror]
        else:
            logmessage = None

        if key:
            auth_module.autoLogin(key)

        ticket = request.params.get('ticket', default=None)

        footer = ws.db.find_one(
            {'_id': 'cservice.frontend'},
            {'login_footer': 1}
        )
        if footer is not None and 'login_footer' in footer:
            footer = footer['login_footer']
        else:
            footer = None

        if not ticket and not session.get('auth_on', False):
            # Build cservice dict for login page templating
            cservices = {
                'webserver': {provider: 1 for provider in ws.providers},
                'logmessage': logmessage,
                'login_footer': footer
            }

            records = ws.db.find(
                {'crecord_name': {'$in': ['casconfig', 'ldapconfig', 'saml2config']}},
                namespace='object'
            )

            ws.logger.info(u'found {} cservices'.format(len(records)))

            context = {}

            for cservice in records:
                cservice = cservice.dump()
                cname = cservice['crecord_name']

                ws.logger.info(u'found cservices type {}'.format(cname))

                if cname == 'casconfig':
                    cservice['server'] = cservice['server'].rstrip('/')
                    cservice['service'] = cservice['service'].rstrip('/')
                    ws.logger.info(u'cas config : server {}, service {}'.format(
                        cservice['server'],
                        cservice['service'],
                    ))

                context[cname] = cservice

            if len(context.keys()) > 0:
                context["auth_ext"] = True

            # Compile template
            login_page = os.path.join(ws.root_directory, 'login', 'index.html')
            with open(login_page) as src:
                tmplsrc = src.read()

            tmpl = Template(tmplsrc)
            return tmpl(context)

        else:
            redirect('/{0}/static/canopsis/index.html'.format(lang))
