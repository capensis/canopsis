#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

import unittest

import logging
logging.basicConfig(format='%(asctime)s %(name)s %(levelname)s %(message)s')

from canopsis.file import File, get_cfile

from canopsis.account import Account
from canopsis.storage import get_storage

from gridfs.errors import NoFile

anonymous_account = Account()
root_account = Account(user="root", group="root")

storage = get_storage(account=root_account , namespace='unittest', logging_level=logging.DEBUG)

sample_file_path = '/opt/canopsis/var/www/canopsis/themes/canopsis/resources/images/logo_small.png'
sample_binary = open(sample_file_path, 'rb').read()

sample_binary2 = bin(1234567890123456789)

myfile = None


class KnownValues(unittest.TestCase):

    def test_01_Init(self):
        global myfile
        myfile = File(storage=storage)

        if myfile.data != {}:
            raise Exception('Data corruption ...')

    def test_02_put_data(self):
        myfile.put_data(sample_binary2)

        if myfile.binary != sample_binary2:
            raise Exception('Data corruption ...')

    def test_03_save_data(self):
        global meta_id, bin_id
        meta_id = myfile.save()
        bin_id = myfile.get_binary_id()

        print("Meta Id: %s, Binary Id: %s" % (meta_id, bin_id))

        if not bin_id or not meta_id:
            raise Exception('Impossible to save File')

    def test_04_put_file(self):
        myfile.put_file(sample_file_path)

        if myfile.binary != sample_binary:
            raise Exception('Data corruption ...')

    def test_05_save_file(self):
        global meta_id, bin_id
        meta_id = myfile.save()

        bin_id = myfile.get_binary_id()
        if not bin_id or not meta_id:
            raise Exception('Impossible to save File')

    def test_06_Rights(self):

        with self.assertRaises(ValueError):
            storage.put(myfile, account=anonymous_account)

        with self.assertRaises(ValueError):
            storage.remove(myfile, account=anonymous_account)

    def test_07_GetMeta(self):
        meta = storage.get(meta_id)
        if not meta:
            raise Exception('Impossible to get meta data')

        print("Meta: %s" % meta)

    def test_08_GetBinary(self):
        binary = storage.get_binary(bin_id)
        if not binary:
            raise Exception('Impossible to get binary data')

        if binary != sample_binary:
            raise Exception('Data corruption ...')

    def test_09_RemoveFile(self):
        myfile.remove()

    def test_10_CheckFileRemove(self):
        with self.assertRaises(NoFile):
            binary = storage.get_binary(bin_id)

        with self.assertRaises(KeyError):
            get_cfile(meta_id, storage)

        if myfile.check():
            raise Exception('File is not deleted ...')

if __name__ == "__main__":
    unittest.main(verbosity=2)
