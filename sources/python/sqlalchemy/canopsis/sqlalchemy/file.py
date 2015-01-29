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

from uuid import uuid4 as uuid

from canopsis.storage.file import FileStorage, FileStream
from canopsis.mongo.storage import MongoStorage
from canopsis.common.init import basestring
from canopsis.common.utils import ensure_iterable

from gridfs import GridFs


class MongoFileStream(FileStream):

    def __init__(self, gridout):

        self.gridout = gridout

    def close(self):
        self.gridout.close()

    def write(self, data):
        self.gridout.write(data=data)

    def read(self, size=1):
        return self.gridout.read(size=size)

    def seek(self, pos, from_beginning=False):
        self.gridout.seek(pos=pos, whence=from_beginning)

    def teel(self):
        return self.gridout.teel()

    def next(self):
        return MongoFileStream(self.gridout.next())


class MongoFileStorage(MongoStorage, FileStorage):

    FILENAME = 'filename'

    def _connect(self, **kwargs):

        result = super(FileStorage, self)._connect(**kwargs)

        if result:

            self.gridfs = GridFs(
                database=self._database, collection=self.get_table())

        return result

    def get(self, names, version=-1):

        names = ensure_iterable(names)
        result = []
        for name in names:
            gridout = self.gridfs.get_version(filename=name, version=version)
            fs = MongoFileStream(gridout)
            result.append(fs)

        return result

    def exists(self, name):

        result = self.gridfs.exists(name)

        return result

    def find(self, names=None, meta=None, sort=None, limit=-1, skip=0):

        request = {}

        if names is not None:
            if isinstance(names, basestring):
                request[MongoFileStorage.FILENAME] = names
            else:
                request[MongoFileStorage.FILENAME] = {'$in': names}

        if meta is not None:
            request.update(meta)

        cursor = self.gridfs.find(request)

        if sort is not None:
            cursor.sort(sort)
        if limit > 0:
            cursor.limit(limit)
        if skip > 0:
            cursor.skip(skip)

        result = (MongoFileStream(gridout) for gridout in cursor)

        return result

    def new_file(self, name=None, meta=None, data=None):

        kwargs = {}

        if name is None:
            name = str(uuid())

        kwargs['_id'] = name
        kwargs['filename'] = name

        if meta is not None:
            kwargs['metadata'] = meta

        gridout = self.gridfs.new_file(**kwargs)

        result = MongoFileStream(gridout)

        if data is not None:
            result.write(data)

        return result

    def delete(self, names):

        names = ensure_iterable(names)
        for name in names:
            self.gridfs.delete(file_id=name)
