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

from canopsis.webcore.services.auth import create_session, check_root, check_group_rights


class BaseBackend(object):
    name = 'base'

    def __init__(self, authorized_grp=None, unauthorized_grp=None, *args, **kwargs):
        super(BaseBackend, self).__init__(*args, **kwargs)

        self.logger = logging.getLogger('auth.backend.{0}'.format(self.name))

        self.orig_authorized_grp = authorized_grp or []
        self.orig_unauthorized_grp = unauthorized_grp or []

        self.authorized_grp = self.orig_authorized_grp
        self.unauthorized_grp = self.orig_unauthorized_grp

    def setup_config(self, context):
        conf = context['config'].get('checkAuthPlugin', {})

        self.authorized_grp = conf.get('authorized_grp', self.orig_authorized_grp)
        self.unauthorized_grp = conf.get('unauthorized_grp', self.orig_unauthorized_grp)

    def install_account(self, account):
        if not check_root(account):
            allowed = (len(self.authorized_grp) == 0)

            self.logger.debug('Ensure user {0} is in allowed groups'.format(account.user))
            for group in self.authorized_grp:
                if check_group_rights(account, group):
                    allowed = True
                    break

            if allowed:
                self.logger.debug('Ensure user {0} is not in forbidden groups'.format(account.user))

                for group in self.unauthorized_grp:
                    if check_group_rights(account, group):
                        allowed = False
                        break

            if not allowed:
                return False

        self.logger.debug('Creating session for account {0}'.format(account.user))

        create_session(account)

        return True
