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
from canopsis.middleware.core import Middleware


CONF_PATH = 'migration/serie.conf'
CATEGORY = 'SERIE'
CONTENT = [
    Parameter('storage')
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class SerieModule(MigrationModule):
    @property
    def storage(self):
        if not hasattr(self, '_storage'):
            self.storage = None

        return self._storage

    @storage.setter
    def storage(self, value):
        if value is None:
            value = 'storage-default-serie2://'

        if isinstance(value, basestring):
            value = Middleware.get_middleware_by_uri(value)

        self._storage = value

    def __init__(self, storage=None, *args, **kwargs):
        super(SerieModule, self).__init__(*args, **kwargs)

        if storage is not None:
            self.storage = storage

    def init(self):
        pass

    def update(self):
        items = self.storage.find_elements(query={
            'computations_per_interval': {'$exists': False}
        })

        for item in items:
            item['computations_per_interval'] = 1
            self.storage.put_element(element=item, _id=item['_id'])
