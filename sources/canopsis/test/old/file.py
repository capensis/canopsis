#!/usr/bin/env python
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

import unittest

from canopsis.old.file import File, get_cfile
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from gridfs.errors import NoFile
import sys
import os

anonymous_account = Account()
root_account = Account(user="root", group="root")

storage = get_storage(
    account=root_account, namespace='unittest')

sample_file_path = os.path.join(
    sys.prefix,
    'var', 'www', 'login', 'img', 'canopsis.png'
)

with open(os.path.expanduser(sample_file_path), 'rb') as f:
    sample_binary = f.read()

sample_binary2 = bin(1234567890123456789)

myfile = None
meta_id = None
bin_id = None


class KnownValues(unittest.TestCase):

    def test_01_Init(self):
        global myfile
        myfile = File(storage=storage)

        self.assertEqual(myfile.data, {})

    def test_02_put_data(self):
        myfile.put_data(sample_binary2)

        self.assertEqual(myfile.binary, sample_binary2)

    def test_03_save_data(self):
        global meta_id, bin_id

        meta_id = myfile.save()
        bin_id = myfile.get_binary_id()

        self.assertTrue(bin_id and meta_id)

    def test_04_put_file(self):
        myfile.put_file(sample_file_path)

        self.assertEqual(myfile.binary, sample_binary)

    def test_05_save_file(self):
        global meta_id, bin_id
        meta_id = myfile.save()

        bin_id = myfile.get_binary_id()

        self.assertTrue(bin_id and meta_id)

    def test_06_Rights(self):

        with self.assertRaises(ValueError):
            storage.put(myfile, account=anonymous_account)

        with self.assertRaises(ValueError):
            storage.remove(myfile, account=anonymous_account)

    def test_07_GetMeta(self):
        global meta_id

        meta = storage.get(meta_id)
        self.assertTrue(meta)

    def test_08_GetBinary(self):
        global bin_id

        binary = storage.get_binary(bin_id)

        self.assertEqual(binary, sample_binary)

    def test_09_RemoveFile(self):
        myfile.remove()

    def test_10_CheckFileRemove(self):
        global meta_id, bin_id

        with self.assertRaises(NoFile):
            storage.get_binary(bin_id)

        with self.assertRaises(KeyError):
            get_cfile(meta_id, storage)

        self.assertFalse(myfile.check())


if __name__ == "__main__":
    unittest.main(verbosity=2)
