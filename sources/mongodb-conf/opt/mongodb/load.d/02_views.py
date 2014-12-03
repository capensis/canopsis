#!/usr/bin/env python
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

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.record import Record

import os
import sys
import json

logger = None

views_path = os.path.expanduser('~/opt/mongodb/load.d/views')

# Set root account
root = Account(user="root", group="root")
storage = get_storage(account=root, namespace='object')


def init(clear=True):

    for path, folders, files in os.walk(views_path):
        for filename in files:
            filepath = os.path.join(path, filename)

            with open(filepath) as f:
                try:
                    data = json.loads(f.read())
                except Exception, e:
                    print (
                        '\n + Error while loading JSON file {} :' +
                        ' {}, aborting...\n'.format(filepath, e))
                    sys.exit(1)

                try:
                    _id = data.pop('id')
                except KeyError as err:
                    print >>sys.stderr, "Can't parse view, missing key:", err
                    sys.exit(1)

                create_view(_id, filename, data, clear)


def update():
    # init takes care of existing records
    init(clear=False)


def create_view(_id, name, data, clear, mod='o+r'):

    try:
        record = storage.get(_id)
        if clear:
            storage.remove(record)
            record = None
        logger.info('View already exists {} no insert done'.format(_id))
    except Exception as e:
        record = None

    if record is None:
        logger.info(' + Create view {}'.format(_id))
        record = Record(data, _type='view', name=name, group='group.CPS_view')
        record.chmod(mod)
        storage.put(record)
