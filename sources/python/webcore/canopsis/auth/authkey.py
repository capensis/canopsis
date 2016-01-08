#!/usr/bin/env python2.7
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

from canopsis.auth.base import BaseBackend

from bottle import request, HTTPError


class AuthKeyBackend(BaseBackend):
    name = 'AuthKeyBackend'

    def __init__(self, *args, **kwargs):
        super(AuthKeyBackend, self).__init__(*args, **kwargs)

    def apply(self, callback, context):
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = self.session.get()

            if not s.get('auth_on', False):
                user = self.do_auth()

                if user and not self.install_account(user['_id'], user):
                    self.logger.error('User {0} not allowed'.format(
                        user['_id']
                    ))

                    return HTTPError(403, 'Forbidden')

            return callback(*args, **kwargs)

        return decorated

    def do_auth(self):
        authkey = request.params.get('authkey', None)

        self.logger.info(
            'Trying to authenticate user with authkey: {0}'.format(
                authkey
            )
        )

        if authkey:
            user = self.auth.check(mode='authkey', password=authkey)

        else:
            user = self.session.get_user()

        if not user:
            self.logger.error('Authentication failed for authkey: {0}'.format(
                authkey
            ))

        return user


def get_backend(ws):
    return AuthKeyBackend(ws)
