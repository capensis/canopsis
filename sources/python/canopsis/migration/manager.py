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

from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.model import Parameter

from canopsis.common.utils import lookup

from logging import StreamHandler
import signal
import json
import os


CONF_PATH = 'migration/manager.conf'
CATEGORY = 'MIGRATION'
CONTENT = [
    Parameter('modules', parser=Parameter.array())
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class MigrationTool(Configurable):

    @property
    def modules(self):
        if not hasattr(self, '_modules'):
            self.modules = None

        return self._modules

    @modules.setter
    def modules(self, value):
        if value is None:
            value = []

        self._modules = value

    def __init__(self, modules=None, *args, **kwargs):
        super(MigrationTool, self).__init__(*args, **kwargs)

        if modules is not None:
            self.modules = modules

        self.loghandler = StreamHandler()
        self.logger.addHandler(self.loghandler)

    def fill(self, init=True):
        tools = []

        for module in self.modules:
            try:
                migrationcls = lookup(module)

            except ImportError as err:
                self.logger.error(
                    'Impossible to load module "{0}": {1}'.format(
                        module,
                        err
                    )
                )

                continue

            migrationtool = migrationcls()
            migrationtool.logger.addHandler(self.loghandler)
            tools.append(migrationtool)

        for tool in tools:
            if init:
                tool.init()

            else:
                tool.update()


@conf_paths(CONF_PATH)
@add_category('MODULE', content=[
    Parameter('ask_timeout', parser=int),
    Parameter('version_info')
])
class MigrationModule(Configurable):

    @property
    def ask_timeout(self):
        if not hasattr(self, '_ask_timeout'):
            self.ask_timeout = None

        return self._ask_timeout

    @ask_timeout.setter
    def ask_timeout(self, value):
        if value is None:
            value = 30

        self._ask_timeout = value

    @property
    def version_info(self):
        if not hasattr(self, '_version_info'):
            self.version_info = None

        return self._version_info

    @version_info.setter
    def version_info(self, value):
        if value is None:
            value = '~/var/lib/canopsis/migration.json'

        self._version_info = os.path.expanduser(value)

    def __init__(self, ask_timeout=None, version_info=None, *args, **kwargs):
        super(MigrationModule, self).__init__(*args, **kwargs)

        if ask_timeout is not None:
            self.ask_timeout = ask_timeout

        if version_info is not None:
            self.version_info = version_info

    def get_version(self, item):
        try:
            with open(self.version_info) as f:
                version_info = json.load(f)

        except Exception as err:
            self.logger.error(
                'Impossible to parse version info: {0}'.format(err)
            )

            version_info = {}

        return version_info.get(item, 0)

    def set_version(self, item, version):
        try:
            with open(self.version_info) as f:
                version_info = json.load(f)

        except Exception as err:
            self.logger.error(
                'Impossible to parse version info: {0}'.format(err)
            )

            version_info = {}

        version_info[item] = version

        try:
            with open(self.version_info, 'w') as f:
                json.dump(version_info, f)

        except Exception as err:
            self.logger.error(
                'Impossible to save version info: {0}'.format(err)
            )

    def ask(self, prompt, default=True):
        answered = False
        user_input = 'N'
        default_val = 'Y' if default else 'N'

        def timeout(sig, frame):
            raise Exception('')

        signal.signal(signal.SIGALRM, timeout)

        while not answered:
            signal.alarm(self.ask_timeout)

            try:
                user_input = raw_input(
                    '{0} Y/N (default={1})'.format(
                        prompt,
                        default_val
                    )
                )

                if user_input in ['Y', 'y', 'N', 'n', '']:
                    answered = True

            except Exception:
                user_input = default_val
                answered = True

            signal.alarm(0)

        if user_input == '':
            user_input = default_val

        return (user_input in ['Y', 'y'])

    def init(self):
        raise NotImplementedError()

    def update(self):
        raise NotImplementedError()
