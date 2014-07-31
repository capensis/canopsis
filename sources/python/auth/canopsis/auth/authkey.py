#!/usr/bin/env python2.7
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

from bottle import request, HTTPError

from canopsis.webcore.services.auth import get_account, check_authkey
from canopsis.auth.base import BaseBackend


class AuthKeyBackend(BaseBackend):
    name = 'AuthKeyBackend'

    def apply(self, callback, context):
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = request.environ.get('beaker.session')

            if not s.get('auth_on', False):
                account = self.do_auth()

                if account and not self.install_account(account):
                    self.logger.error('User {0} not allowed'.format(account.user))

                    return HTTPError(403, 'Forbidden')

            return callback(*args, **kwargs)

        return decorated

    def do_auth(self):
        authkey = request.params.get('authkey', None)

        self.logger.info('Trying to authenticate user with authkey: {0}'.format(authkey))

        if authkey:
            account = check_authkey(authkey)

        else:
            account = get_account()

        if not account or account.user == 'anonymous':
            self.logger.error('Authentication failed for authkey: {0}'.format(authkey))
            # Will try with the next backend
            return None

        return account


def get_backend():
    return AuthKeyBackend()
