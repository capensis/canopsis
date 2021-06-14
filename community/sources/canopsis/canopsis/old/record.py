# -*- coding: utf-8 -*-
# !/usr/bin/env python
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

from bson import objectid


class Record(object):
    def __init__(
        self,
        data=None,
        _id=None,
        raw_record=None,
        record=None,
        storage=None,
        account=None,
        _type='raw',
        name='noname'
    ):

        super(Record, self).__init__()

        self.write_time = None
        self._id = _id
        self.enable = True
        self.type = _type
        self.name = name

        if data is None:
            data = {}

        if '_id' in data:
            self._id = data.get('_id')
            del data['_id']

        self.data = data.copy()
        self.storage = storage

        if isinstance(record, Record):
            self.load(record.dump())
        elif raw_record and isinstance(raw_record, dict):
            self.load(raw_record)

    def load(self, dump):

        self.type = str(dump.get('crecord_type', ''))
        self.write_time = dump.get('crecord_write_time', '')
        self.name = dump.get('crecord_name', '')
        self.enable = dump.get('enable', '')

        dump['_id'] = dump.get('_id', None)

        self._id = dump['_id']

        dump.pop('_id', '')
        dump.pop('enable', '')
        dump.pop('crecord_type', '')
        dump.pop('crecord_write_time', '')
        dump.pop('crecord_name', '')

        self.data = dump.copy()

    def save(self, storage=None):
        if not storage:
            if not self.storage:
                raise Exception('For save you must specify storage')
            else:
                storage = self.storage

        return storage.put(self)

    def get(self, key):
        if hasattr(self, key):
            return getattr(self, key)
        else:
            return None

    def dump(self, json=False):
        dump = self.data.copy()

        dump['_id'] = self.get('_id')
        dump['crecord_type'] = self.get('type')
        dump['crecord_write_time'] = self.get('write_time')
        dump['crecord_name'] = self.get('name')

        if 'enable' not in self.data:
            dump['enable'] = self.enable

        if json:
            # Clean objectid
            for key in dump:
                if isinstance(dump[key], objectid.ObjectId):
                    dump[key] = str(dump[key])

        return dump

    def __str__(self):
        return str(self.dump())

    def __repr__(self):
        return self.__str__()

    def is_enable(self):
        return self.enable

    def set_enable(self, autosave=True):
        self.enable = True
        if autosave and self.storage:
            self.save()

    def set_disable(self, autosave=True):
        self.enable = False
        if autosave and self.storage:
            self.save()
