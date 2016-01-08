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

"""Storage synchronization module."""

from canopsis.configuration.model import Parameter, Configuration
from canopsis.middleware.registry import MiddlewareRegistry


class Synchronizer(MiddlewareRegistry):
    """Synchronize data from a source storage to target storages.

    Targets must respect a common storage type with the source storage.
    """

    CONF_FILE = 'middleware/sync.conf'

    SOURCE = 'source'
    TARGETS = 'targets'

    CATEGORY = 'SYNC_MIDDLEWARE'

    def __init__(self, source, targets, *args, **kwargs):
        """

        :param source: source storage to synchronize with targets
        :type source: str

        :param targets: targets to synchronize with source
        :type targets: list of str
        """

        super(Synchronizer, self).__init__(*args, **kwargs)

        self.source = source
        self.targets = targets

    @property
    def source(self):
        """Get source storage.
        """

        return self._source

    @source.setter
    def source(self, value):
        """Change of source storage.
        """

        self._source = value

    @property
    def targets(self):
        """Get target storage types.
        """

        return self._targets

    @targets.setter
    def targets(self, value):
        """Change of target storage types.
        """

        self._targets = value

    def copy(self, source=None, targets=None):
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
                Parameter(Synchronizer.SOURCE),
                Parameter(Synchronizer.TARGETS)
            )
        )

        return result

    def _configure(self, conf, *args, **kwargs):

        super(Synchronizer, self)._configure(conf=conf, *args, **kwargs)

        values = conf[Configuration.VALUES]

        # set shared
        self._update_parameter(values, Synchronizer.SOURCE)
        self._update_parameter(values, Synchronizer.TARGETS)
