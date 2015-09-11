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

from canopsis.webcore.services import auth, session, rights


class BaseBackend(object):
    name = 'base'
    handle_logout = False

    def __init__(self, ws, *args, **kwargs):
        super(BaseBackend, self).__init__(*args, **kwargs)

        self.ws = ws
        self.auth = auth
        self.session = session
        self.rights = rights

        self.logger = ws.logger

        self._perms = []

    def setup_config(self, context):
        self.permissions = context['config'].get('permissions', self._perms)

    def install_account(self, uid, user):
        mgr = self.rights.get_manager()

        self.logger.info('Ensure user {0} has sufficient rights'.format(
            uid
        ))

        for p in self.permissions:
            right_id, checksum = p

            if not mgr.check_rights(uid, right_id, checksum):
                return False

        self.logger.info('Creating session for user {0}'.format(uid))

        self.session.create(user)

        return True
