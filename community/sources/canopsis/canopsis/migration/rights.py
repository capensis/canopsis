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

import json
import os
from uuid import uuid1

from canopsis.common.utils import ensure_iterable
from canopsis.confng import Configuration, Ini
from canopsis.logger import Logger
from canopsis.migration.manager import MigrationModule
from canopsis.organisation.rights import Rights

DEFAULT_ACTIONS_PATH = '~/opt/mongodb/load.d/rights/actions_ids'
DEFAULT_USERS_PATH = '~/opt/mongodb/load.d/rights/default_users'
DEFAULT_ROLES_PATH = '~/opt/mongodb/load.d/rights/default_roles'


class RightsModule(MigrationModule):

    CONF_PATH = 'etc/migration/rights.conf'
    CATEGORY = 'RIGHTS'

    def __init__(
            self,
            actions_path=None,
            users_path=None,
            roles_path=None,
            *args, **kwargs
    ):
        super(RightsModule, self).__init__(*args, **kwargs)

        self.logger = Logger.get('migrationtool', MigrationModule.LOG_PATH)
        self.config = Configuration.load(RightsModule.CONF_PATH, Ini)
        conf = self.config.get(self.CATEGORY, {})

        self.manager = Rights()

        if actions_path is not None:
            actions_path = actions_path
        else:
            actions_path = conf.get('actions_path', DEFAULT_ACTIONS_PATH)
        self.actions_path = os.path.expanduser(actions_path)

        if users_path is not None:
            users_path = users_path
        else:
            users_path = conf.get('users_path', DEFAULT_USERS_PATH)
        self.users_path = os.path.expanduser(users_path)

        if roles_path is not None:
            roles_path = roles_path
        else:
            roles_path = conf.get('roles_path', DEFAULT_ROLES_PATH)
        self.roles_path = os.path.expanduser(roles_path)

    def init(self, clear=True, yes=False):
        self.add_actions(self.load(self.actions_path), clear)
        self.add_users(self.load(self.users_path), clear)
        self.add_roles(self.load(self.roles_path), clear)

    def update(self, yes=False):
        self.init(clear=False)

    def load(self, path):
        try:
            loaded = []

            for fpath in os.listdir(path):
                if fpath.endswith('.json'):
                    fullpath = os.path.join(path, fpath)

                    with open(fullpath) as f:
                        data = ensure_iterable(json.load(f))

                    loaded += data

        except Exception as err:
            self.logger.error(u'Unable to load JSON files "{0}": {1}'.format(
                path,
                err
            ))

            loaded = []

        return loaded

    def add_actions(self, data, clear):
        for action in data:
            for aid in action:
                if self.manager.get_action(aid) is None or clear:
                    self.logger.info(u'Initialize action: {0}'.format(aid))

                    self.manager.add(
                        aid,
                        action[aid].get('desc', 'Empty description')
                    )

    def add_users(self, data, clear):
        for default_user in data:
            user = self.manager.get_user(default_user['_id'])

            if user is None or clear:
                self.logger.info(u'Initialize user: '
                                 '{0}'.format(default_user['_id']))

                self.manager.create_user(
                    default_user['_id'],
                    default_user.get('role', None),
                    rights=default_user.get('rights', None),
                    contact=default_user.get('contact', None),
                    groups=default_user.get('groups', None)
                )

                user = self.manager.get_user(default_user['_id'])

            self.manager.update_fields(
                user['_id'],
                'user',
                {
                    'external': user.get('external',
                                         default_user.get('external', False)),
                    'enable': user.get('enable', default_user.get('enable',
                                                                  True)),
                    'shadowpasswd': user.get('shadowpasswd',
                                             default_user.get('shadowpasswd',
                                                              None)),
                    'mail': user.get('mail', default_user.get('mail', None)),
                    'authkey': user.get('authkey',
                                        default_user.get('authkey',
                                                         str(uuid1())))
                }
            )

    def add_roles(self, data, clear):
        for role in data:
            if self.manager.get_role(role['_id']) is None or clear:
                self.logger.info(u'Initialize role: {0}'.format(role['_id']))

                self.manager.create_role(
                    role['_id'],
                    role.get('profile', None)
                )

            self.logger.info(u'Updating role: {0}'.format(role['_id']))
            record = self.manager.get_role(role['_id'])

            rights = record.get('rights', {})
            groups = record.get('groups', [])

            rights.update(role.get('rights', {}))
            groups += role.get('groups', [])
            groups = list(set(groups))  # make groups unique

            self.manager.update_rights(role['_id'], 'role', rights, record)
            self.manager.update_group(role['_id'], 'role', groups, record)
            self.manager.update_fields(role['_id'], 'role', {
                                       "defaultview": role.get("defaultview", None)})
