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

import logging


class BaseBackend(object):
    name = 'base'
    handle_logout = False

    def __init__(self, ws, *args, **kwargs):
        super(BaseBackend, self).__init__(*args, **kwargs)

        self.ws = ws
        self.session = ws.require('session')
        self.rights = ws.require('rights').get_manager()
        self.auth = ws.require('auth')

        self.logger = logging.getLogger('auth.backend.{0}'.format(self.name))

        self._perms = []

    def setup_config(self, context):
        self.permissions = context['config'].get('permissions', self._perms)

    def install_account(self, user):
        self.logger.debug('Ensure user {0} has sufficient rights'.format(
            user['_id']
        ))

        for p in self.permissions:
            right_id, checksum = p

            if not self.rights.check_rights(user['_id'], right_id, checksum):
                return False

        self.logger.debug('Creating session for user {0}'.format(user['_id']))

        self.session.create(user)

        return True
