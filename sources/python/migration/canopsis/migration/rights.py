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

from canopsis.migration.manager import MigrationModule
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.model import Parameter

from canopsis.common.utils import ensure_iterable
from canopsis.organisation.rights import Rights

from uuid import uuid1
import json
import os


CONF_PATH = 'migration/rights.conf'
CATEGORY = 'RIGHTS'
CONTENT = [
    Parameter('actions_path'),
    Parameter('users_path'),
    Parameter('roles_path')
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class RightsModule(MigrationModule):

    @property
    def actions_path(self):
        if not hasattr(self, '_actions_path'):
            self.actions_path = None

        return self._actions_path

    @actions_path.setter
    def actions_path(self, value):
        if value is None:
            value = '~/opt/mongodb/load.d/rights/actions_ids'

        self._actions_path = os.path.expanduser(value)

    @property
    def users_path(self):
        if not hasattr(self, '_users_path'):
            self.users_path = None

        return self._users_path

    @users_path.setter
    def users_path(self, value):
        if value is None:
            value = '~/opt/mongodb/load.d/rights/default_users'

        self._users_path = os.path.expanduser(value)

    @property
    def roles_path(self):
        if not hasattr(self, '_roles_path'):
            self.roles_path = None

        return self._roles_path

    @roles_path.setter
    def roles_path(self, value):
        if value is None:
            value = '~/opt/mongodb/load.d/rights/default_roles'

        self._roles_path = os.path.expanduser(value)

    def __init__(
        self,
        actions_path=None,
        users_path=None,
        roles_path=None,
        *args, **kwargs
    ):
        super(RightsModule, self).__init__(*args, **kwargs)

        self.manager = Rights()

        if actions_path is not None:
            self.actions_path = actions_path

        if users_path is not None:
            self.users_path = users_path

        if roles_path is not None:
            self.roles_path = roles_path

    def init(self, clear=True):
        self.add_actions(self.load(self.actions_path), clear)
        self.add_users(self.load(self.users_path), clear)
        self.add_roles(self.load(self.roles_path), clear)

    def update(self):
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
            self.logger.error('Unable to load JSON files "{0}": {1}'.format(
                path,
                err
            ))

            loaded = []

        return loaded

    def add_actions(self, data, clear):
        for action in data:
            for aid in action:
                if self.manager.get_action(aid) is None or clear:
                    self.logger.info('Initialize action: {0}'.format(aid))

                    self.manager.add(
                        aid,
                        action[aid].get('desc', 'Empty description')
                    )

    def add_users(self, data, clear):
        for user in data:
            if self.manager.get_user(user['_id']) is None or clear:
                self.logger.info('Initialize user: {0}'.format(user['_id']))

                self.manager.create_user(
                    user['_id'],
                    user.get('role', None),
                    rights=user.get('rights', None),
                    contact=user.get('contact', None),
                    groups=user.get('groups', None)
                )

            self.manager.update_fields(
                user['_id'],
                'user',
                {
                    'external': user.get('external', False),
                    'enable': user.get('enable', True),
                    'shadowpasswd': user.get('shadowpass', None),
                    'mail': user.get('mail', None),
                    'authkey': user.get('authkey', str(uuid1()))
                }
            )

    def add_roles(self, data, clear):
        for role in data:
            if self.manager.get_role(role['_id']) is None or clear:
                self.logger.info('Initialize role: {0}'.format(role['_id']))

                self.manager.create_role(
                    role['_id'],
                    role.get('profile', None)
                )

            self.logger.info('Updating role: {0}'.format(role['_id']))
            record = self.manager.get_role(role['_id'])

            rights = record.get('rights', {})
            groups = record.get('groups', [])

            rights.update(role.get('rights', {}))
            groups += role.get('groups', [])
            groups = list(set(groups))  # make groups unique

            self.manager.update_rights(role['_id'], 'role', rights, record)
            self.manager.update_group(role['_id'], 'role', groups, record)
