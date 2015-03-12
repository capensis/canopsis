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

from canopsis.engines.core import Engine
from canopsis.old.account import Account
from canopsis.old.storage import get_storage


class engine(Engine):
    """
    In charge of doing write operation on DB.
    """

    etype = "writer"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.clean_collection = {}
        collections_to_clean = ['events', 'events_log']

        for collection in collections_to_clean:
            self.clean_collection[collection] = get_storage(
                collection,
                account=Account(user='root')
            ).get_backend()

        self.object = get_storage(
            'object',
            account=Account(user='root')
        ).get_backend()
