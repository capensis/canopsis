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
from canopsis.mongo.core import MongoStorage
from canopsis.common.init import basestring
from canopsis.common.utils import ensure_iterable

from gridfs import GridFS, NoFile


class MongoFileStream(FileStream):

    @property
    def name(self):
        return self.gridout.filename

    def __init__(self, gridout):

        self.gridout = gridout

    def close(self):
        self.gridout.close()

    def write(self, data):
        self.gridout.write(data=data)

    def read(self, size=-1):

        return self.gridout.read(size=size)

    def seek(self, pos, from_beginning=False):

        self.gridout.seek(pos=pos, whence=from_beginning)

    def pos(self):

        return self.gridout.tell()

    def next(self):
        return MongoFileStream(self.gridout.next())

    def get_inner_object(self):
        return self.gridout

    def __eq__(self, other):

        return (
            isinstance(other, self.__class__)
            and self.gridout.filename == other.gridout.filename
        )


class MongoFileStorage(MongoStorage, FileStorage):

    FILENAME = 'filename'

    def __init__(self, *args, **kwargs):

        super(MongoFileStorage, self).__init__(*args, **kwargs)

        self.gridfs = None

    def _connect(self, **kwargs):

        result = super(MongoFileStorage, self)._connect(**kwargs)

        if result:

            self.gridfs = GridFS(
                database=self._database, collection=self.get_table()
            )

        return result

    def put(self, name, data, meta=None):

        try:
            fs = self.new_file(name=name, meta=meta)
            fs.write(data=data)
        finally:
            fs.close()

    def put_meta(self, name, meta):
        try:
            oldf, _meta = self.get(name, with_meta=True)
            _meta.update(meta)

            fs = self.new_file(name=name, meta=_meta)

            while True:
                data = oldf.read(512)

                if not data:
                    break

                fs.write(data=data)

        finally:
            fs.close()

    def get(self, name, version=-1, with_meta=False):

        result = None

        try:
            gridout = self.gridfs.get_version(filename=name, version=version)
        except NoFile:
            pass
        else:
            if with_meta:
                result = MongoFileStream(gridout), gridout.metadata

            else:
                result = MongoFileStream(gridout)

        return result

    def get_meta(self, name):
        result = self.get(name, with_meta=True)

        if result is not None:
            result = result[1]

        return result

    def exists(self, name):

        result = self.gridfs.exists(filename=name)

        return result

    def find(
        self,
        names=None,
        meta=None,
        sort=None,
        limit=-1,
        skip=0,
        with_meta=False
    ):

        request = {}

        if names is not None:
            if isinstance(names, basestring):
                request[MongoFileStorage.FILENAME] = names
            else:
                request[MongoFileStorage.FILENAME] = {'$in': names}

        if meta is not None:
            for metafield in meta:
                field = 'metadata.{0}'.format(metafield)
                request[field] = meta[metafield]

        cursor = self.gridfs.find(request)

        if sort is not None:
            cursor.sort(sort)
        if limit > 0:
            cursor.limit(limit)
        if skip > 0:
            cursor.skip(skip)

        if with_meta:
            result = (
                (MongoFileStream(gridout), gridout.metadata)
                for gridout in cursor
            )

        else:
            result = (MongoFileStream(gridout) for gridout in cursor)

        return result

    def list(self):

        return self.gridfs.list()

    def new_file(self, name=None, meta=None, data=None):

        kwargs = {}

        if name is None:
            name = str(uuid())

        kwargs['filename'] = name

        if meta is not None:
            kwargs['metadata'] = meta

        gridout = self.gridfs.new_file(**kwargs)

        result = MongoFileStream(gridout)

        if data is not None:
            result.write(data)

        return result

    def delete(self, names=None):

        if names is None:
            names = self.gridfs.list()

        names = ensure_iterable(names)

        for name in names:
            while True:
                fs = self.get(name)

                if fs is None:
                    break

                self.gridfs.delete(file_id=fs.get_inner_object()._id)
