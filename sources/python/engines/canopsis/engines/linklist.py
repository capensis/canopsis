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
# from canopsis.old.account import Account
# from canopsis.old.storage import get_storage

from datetime import datetime
from time import time


class engine(Engine):
    etype = 'linklist'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

    def consume_dispatcher(self, event, *args, **kargs):

        self.logger.debug('Enter linklist dispatch')
        linklist = self.get_ready_record(event)

        if linklist:

            self.logger.debug('{}: {}'.format(event_id, linklist))
            event_id = event['_id']

            self.crecord_task_complete(event_id)
