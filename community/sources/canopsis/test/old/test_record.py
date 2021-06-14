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

from unittest import TestCase, main

from json import dumps

from canopsis.old.record import Record
from canopsis.old.account import Account
import unittest
from canopsis.common import root_path
import xmlrunner


class KnownValues(TestCase):
    def setUp(self):
        self.anonymous_account = Account()
        self.root_account = Account(user="root", group="root")
        self.user_account = Account(user="william", group="capensis")

        self.data = {
            'mydata1': 'data1', 'mydata2': 'data2', 'mydata3': 'data3'}

    def test_01_Init(self):
        record = Record(self.data)
        if record.data != self.data:
            raise Exception('Data corruption ...')

    def test_02_InitFromRaw(self):
        raw = {
            '_id': None,
            'crecord_name': 'titi',
            'crecord_type': 'raw',
            'crecord_write_time': None,
            'enable': True,
            'mydata1': 'data1',
            'mydata3': 'data3',
            'mydata2': 'data2',
            'crecord_type': 'raw'
        }

        record = Record(raw_record=raw)

        dump = record.dump()

        if not isinstance(dump['_id'], type(None)):
            raise Exception('Invalid _id type')

        if record.data != self.data:
            raise Exception('Data corruption ...')

    def test_03_InitFromRecord(self):
        record = Record(self.data)

        record2 = Record(record=record)

        if record2.data != self.data:
            raise Exception('Data corruption ...')

    def test_07_enable(self):
        record = Record(self.data)

        record.set_enable()
        if not record.is_enable():
            raise Exception('Impossible to enable ...')

        record.set_disable()
        if record.is_enable():
            raise Exception('Impossible to disable ...')

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
