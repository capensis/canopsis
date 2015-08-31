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

from canopsis.configuration.parameters import Parameter, Configuration
from canopsis.middleware.manager import Manager


class Synchronizer(Manager):
    """
    Manage synchronization between middlewares of the same type.
    """

    CONF_FILE = 'middleware/sync.conf'

    SOURCE_MIDDLEWARE = 'source_storage'
    TARGET_MIDDLEWARES = 'target_storages'

    CATEGORY = 'SYNC_MIDDLEWARE'

    def __init__(self, source_storage, target_types, *args, **kwargs):
        """

        :param source: source storage to synchronize with targets
        :type source: str

        :param targets: targets to synchronize with source
        :type targets: list of str
        """

        super(Synchronizer, self).__init__(*args, **kwargs)

        self.source = source_storage
        self.targets = target_types

    @property
    def source_storage(self):
        return self._source_storage

    @source_storage.setter
    def source_storage(self, value):
        self._source_storage = value

    @property
    def target_types(self):
        return self._target_types

    @target_types.setter
    def target_types(self, value):
        self._target_types = value

    def copy(self, source_storage=None, target_types=None):
        """
        Copy content of source storage to target storages
        """

        raise NotImplementedError()

    def _get_conf_files(self, *args, **kwargs):

        result = super(Synchronizer, self)._get_conf_files(*args, **kwargs)

        result.append(Synchronizer.CONF_FILE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Synchronizer, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Synchronizer.CATEGORY,
            new_content=(
                Parameter(Synchronizer.SOURCE_MIDDLEWARE),
                Parameter(Synchronizer.TARGET_MIDDLEWARES)))

        return result

    def _configure(self, conf, *args, **kwargs):

        super(Synchronizer, self)._configure(conf=conf, *args, **kwargs)

        values = conf[Configuration.VALUES]

        # set shared
        self._update_parameter(values, Synchronizer.SOURCE_MIDDLEWARE)
        self._update_parameter(values, Synchronizer.TARGET_MIDDLEWARES)
