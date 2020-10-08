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
import os.path

from canopsis.common.template import Template

from canopsis.webcore.services import auth as auth_module
from canopsis.webcore.services import session as session_module


def exports(ws):
    skip_login = ws.skip_login

    @ws.application.get('/:lang/static/canopsis-next/index.html', skip=skip_login)
    @ws.application.get('/:lang/static/canopsis-next/dist/index.html', skip=skip_login)
    @ws.application.get('/static/canopsis-next/index.html', skip=skip_login)
    @ws.application.get('/static/canopsis-next/dist/index.html', skip=skip_login)
    def index(lang='en'):
        return static_file('canopsis-next/dist/index.html', root=ws.root_directory)

    @ws.application.get('/:lang/static/canopsis-next/<filename:path>', skip=skip_login)
    @ws.application.get('/static/canopsis-next/<filename:path>', skip=skip_login)
    def server_static(filename, lang='en'):
        response.set_header("Access-Control-Allow-Origin", "*")

        filename = os.path.join('canopsis-next', filename)
        return static_file(filename, root=ws.root_directory)

    @ws.application.get('/:lang/static/canopsis-next/dist/favicon.ico', skip=skip_login)
    @ws.application.get('/static/canopsis-next/dist/favicon.ico', skip=skip_login)
    @ws.application.get('/favicon.ico', skip=skip_login)
    def favicon(**kwargs):
        return

    @ws.application.get('/:lang/static/canopsis/index.html', skip=skip_login)
    @ws.application.get('/static/canopsis/index.html', skip=skip_login)
    @ws.application.get('/:lang/static/<filename:path>', skip=skip_login)
    @ws.application.get('/static/<filename:path>', skip=skip_login)
    @ws.application.get('/', skip=skip_login)
    @ws.application.get('/index.html', skip=skip_login)
    @ws.application.get('/:lang/', skip=skip_login)
    @ws.application.get('/:lang/index.html', skip=skip_login)
    @ws.application.get('/:lang/:key/', skip=skip_login)
    @ws.application.get('/:lang/:key/index.html', skip=skip_login)
    def uiv2(lang='en', key=None, **kwargs):
        session = request.environ.get('beaker.session')

        # Try to authenticate user
        key = key or request.params.get('authkey', default=None)

        if key:
            auth_module.autoLogin(key)

        redirect('/{}{}'.format(lang, ws.config.get('ui', {}).get('url', '')))
