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

"""Mongo implementation of the filestorage.
"""

from canopsis.storage.core import Storage


class FileStream(object):
    """File stream object stored in a FileStorage.
    """

    @property
    def name(self):
        """Return file's name.
        """

        raise NotImplementedError()

    def close(self):
        """Close this file stream.
        """

        raise NotImplementedError()

    def write(self, data):
        """Write data in this file stream.

        :param str data: data to write.
        """

        raise NotImplementedError()

    def read(self, size=-1):
        """Read size characters from file content.

        :param int size: number of characters to read. Default all characters.
        """

        raise NotImplementedError()

    def seek(self, pos, from_beginning=False):
        """Set the current position in this file stream.

        :param int pos: seek position.
        :param bool from_beginning: start to seek from the beginning of the
            file.
        """

        raise NotImplementedError()

    def pos(self):
        """Get the current position of this file.
        """

        raise NotImplementedError()

    def next(self):
        """Get next file stream version.
        """

        raise NotImplementedError()

    def get_inner_object(self):
        """Get backend file object.
        """

        raise NotImplementedError()

    def __eq__(self, other):
        """Check that file streams are equal.
        """

        raise NotImplementedError()


class FileStorage(Storage):
    """
    Storage dedicated to manage distributed file stream objects.
    """

    __datatype__ = 'file'  #: registered such as a file storage

    def put(self, name, data):
        """Put data in the storage.

        :param str name: file name.
        :param str data: data to put.
        """

        raise NotImplementedError()

    def get(self, name, version=-1, with_meta=False):
        """Get file stream related to input name.

        :param str name: file stream name.
        :param int version: file stream version. last if -1 (by default).
        :param bool with_meta: return file's metadata.
        :return: corresponding filestream.
        :rtype: FileStream.
        """

        raise NotImplementedError()

    def exists(self, name):
        """True if input file name exists.

        :param str name: file name to check existance.

        :return: True iif a file name equals to input name exists.
        :rtype: bool
        """

        raise NotImplementedError()

    def list(self):
        """Get all file names.

        :return: list of file names.
        :rtype: list
        """

        raise NotImplementedError()

    def find(
            self,
            names=None,
            meta=None,
            sort=None,
            limit=-1,
            skip=0,
            with_meta=False
    ):
        """Try to find file streams where names match with input names or meta
        data match with input meta.

        :param names: file name(s). If str, return one FileStream, list
            otherwise.
        :type names: str or list of str
        :param canopsis.storage.filter.Filter meta: meta information to find.
        :param dict sort: sort criteria.
        :param int limit: limit criteria.
        :param int skip: skip criteria.
        :param bool with_meta: return files with metadata

        :return: file stream(s) depending on input names value.
        :rtype: list or FileStream
        """

        raise NotImplementedError()

    def new_file(self, name=None, meta=None, data=None):
        """Put a file stream with its name and meta data.

        :param str name: file name.
        :param dict meta: file meta data.
        :param str data: file content.

        :rtype: FileStream
        """

        raise NotImplementedError()

    def delete(self, names):
        """Delete file streams related to their names

        :param names: file name(s) to delete.
        :type names: str or list
        """

        raise NotImplementedError()
