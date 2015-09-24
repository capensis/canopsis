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

from canopsis.storage.file import FileStorage, FileStream
from canopsis.common.utils import ensure_iterable
from canopsis.old.mfilter import check

from uuid import uuid4 as uuid
import json
import os


class UnixFileStream(FileStream):
    @property
    def name(self):
        return self.fileobj.name

    def __init__(self, fileobj):
        self.fileobj = fileobj

    def close(self):
        self.fileobj.close()

    def write(self, data):
        self.fileobj.write(data)

    def read(self, size=-1):
        return self.fileobj.read(size)

    def seek(self, pos, from_beginning=False):
        self.fileobj.seek(pos, whence=from_beginning)

    def pos(self):
        return self.fileobj.tell()

    def next(self):
        return UnixFileStream(self.fileobj.next())

    def get_inner_object(self):
        return self.fileobj

    def __eq__(self, other):
        return (
            isinstance(other, self.__class__)
            and self.fileobj.name == other.fileobj.name
        )


class UnixFileStorage(FileStorage):
    META_EXT = '.meta.json'

    def _connect(self, **kwargs):
        return True

    def _path(self, name, with_meta=False):
        fpath = os.path.join(self.uri, name)

        if with_meta:
            mpath = '{0}{1}'.format(fpath, UnixFileStorage.META_EXT)

            return fpath, mpath

        else:
            return fpath

    def put(self, name, data, meta=None):
        fpath = self._path(name)

        with open(fpath, 'w') as f:
            f.write(data)

        if meta is not None:
            self.put_meta(name, meta)

    def put_meta(self, name, meta):
        _, mpath = self._path(name, with_meta=True)

        with open(mpath, 'w') as f:
            json.dump(f, meta)

    def get(self, name, with_meta=False):
        f = open(self._path(name))

        if with_meta:
            return f, self.get_meta(name)

        else:
            return f

    def get_meta(self, name):
        _, mpath = self._path(name, True)

        with open(mpath) as m:
            meta = json.load(m)

        return meta

    def exists(self, name):
        return os.path.exists(self._path(name))

    def find(
        self,
        names=None,
        meta=None,
        sort=None,
        limit=-1,
        skip=0,
        with_meta=False
    ):
        if names is not None:
            names = ensure_iterable(names)

        result = []

        for _, _, files in os.walk(self.uri):
            for filename in files:
                # ignore meta files
                if not filename.endswith(UnixFileStorage.META_EXT):
                    if names is not None and filename not in names:
                        continue

                    if meta is not None:
                        metadata = self.get_meta(filename)

                        if not check(meta, metadata):
                            continue

                    result.append(filename)

        result = result[skip:limit]
        result = [open(filename) for filename in result]

        if with_meta:
            result = [(f, self.get_meta(f.name)) for f in result]

        if sort is not None:
            raise NotImplementedError('sort is not yet supported')

        return result

    def list(self):
        return [
            filename
            for _, _, files in os.walk(self.uri)
            for filename in files
            if not filename.endswith(UnixFileStorage.META_EXT)
        ]

    def new_file(self, name=None, meta=None, data=None):
        if name is None:
            name = str(uuid())

        if meta is None:
            meta = {}

        self.put(name, data, meta)

        return self.get(name)

    def delete(self, names=None):
        if names is None:
            names = self.list()

        names = ensure_iterable(names)

        for name in names:
            os.path.remove(self._path(name))
