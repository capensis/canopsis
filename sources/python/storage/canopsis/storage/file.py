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

from canopsis.storage import Storage


class FileStream(object):
    """File stream object stored in a FileStorage.
    """

    def close(self):
        """Close this file stream.
        """

        raise NotImplementedError()

    def write(self, data):
        """Write data in this file stream.

        :param str data: data to write.
        """

        raise NotImplementedError()

    def read(self, size=1):
        """Read size characters from file content.

        :param int size: number of characters to read.
        """

        raise NotImplementedError()

    def seek(self, pos, from_beginning=False):
        """Set the current position in this file stream.

        :param int pos: seek position.
        :param bool from_beginning: start to seek from the beginning of the
            file.
        """

        raise NotImplementedError()

    def teel(self):
        """Get the current position of this file.
        """

        raise NotImplementedError()

    def next(self):
        """Get next file stream version.
        """

        raise NotImplementedError()


class FileStorage(Storage):
    """
    Storage dedicated to manage distributed file stream objects.
    """

    __datatype__ = 'file'  #: registered such as a file storage

    def get(self, names, version=-1):
        """Get file stream(s) related to input name(s).

        :param names: file stream names.
        :type names: str or list of str
        :param int version: file stream version. last if -1 (by default).

        :rtype: FileStream or list of FileStream
        """
        raise NotImplementedError()

    def exists(self, name):
        """True if input file name exists.

        :param str name: file name to check existance.

        :return: True iif a file name equals to input name exists.
        :rtype: bool
        """

        raise NotImplementedError()

    def find(self, names=None, meta=None, sort=None, limit=-1, skip=0):
        """Try to find file streams where names match with input names or meta
        data match with input meta.

        :param names: file name(s). If str, return one FileStream, list
            otherwise.
        :type names: str or list of str
        :param canopsis.storage.filter.Filter meta: meta information to find.
        :param dict sort: sort criteria.
        :param int limit: limit criteria.
        :param int skip: skip criteria.

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
        :type names: str or list of str
        """

        raise NotImplementedError()
