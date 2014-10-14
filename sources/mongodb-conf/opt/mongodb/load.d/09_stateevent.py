#!/usr/bin/env python
#--------------------------------
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

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.record import Record

##set root account
root = Account(user="root", group="root")

logger = None


def init():
    storage = get_storage(account=root, namespace='object')

    state_spec = {
        "_id": "state-spec",
        "crecord_type": "state-spec",
        "restore_event": True,
        "bagot": {
            # if event appears >= 10 times in 1hr
            "time": 3600,
            "freq": 10},
        # if event appears again in < 5min
        "stealthy_time": 300,
        "stealthy_show": 300}

    logger.info(" + Creating event state specification")
    record = Record(
        data=state_spec, name="event state specifications", _type='state-spec')
    record.chmod('g+w')
    record.chmod('o+r')
    record.chgrp('group.CPS_root')
    storage.put(record)


def update():
    init()
