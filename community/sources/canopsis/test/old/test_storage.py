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

from canopsis.old.storage import Storage
from canopsis.old.record import Record
from canopsis.old.account import Account

from logging import DEBUG, basicConfig

from json import dumps
import unittest
from canopsis.common import root_path
import xmlrunner

STORAGE = None
MYRECORD = None
ID = None

basicConfig(
    level=DEBUG,
    format='%(asctime)s %(name)s %(levelname)s %(message)s')


class KnownValues(TestCase):
    def setUp(self):
        self.anonymous_account = Account()
        self.root_account = Account(user="root", group="root")
        self.user_account = Account(user="william", group="capensis")
        self.user2_account = Account(user="william2", group="capensis")

        #self.anonymous_account.cat()
        #self.user_account.cat()
        #self.root_account.cat()

        self.data = {'mydata1': 'data1', 'mydata2': 'data2', 'mydata3': 'data3'}

    def test_01_Init(self):
        global STORAGE
        STORAGE = Storage(self.user_account, namespace='unittest', logging_level=DEBUG)

        records = STORAGE.find(account=self.root_account)
        STORAGE.remove(records, account=self.root_account)

    def test_02_CreateRecord(self):
        global MYRECORD
        MYRECORD = Record(self.data, storage=STORAGE)

    def test_03_Put(self):
        global ID
        ID = STORAGE.put(MYRECORD)

    # def test_04_RecordSave(self):
    #     MYRECORD.save(STORAGE)
    #     MYRECORD.save()

    # def test_05_Get(self):
    #     global MYRECORD
    #     MYRECORD = STORAGE.get(ID)
    #     del MYRECORD.data['crecord_creation_time']
    #     if MYRECORD.data != self.data:
    #         raise Exception('Invalid data ...')

    # def test_06_Update_field(self):

    #     data = {'mydata1': 'toto'}

    #     STORAGE.update(ID, data)
    #     record = STORAGE.get(ID)

    #     if record.data['mydata1'] != 'toto':
    #         raise Exception('Put_field failed ...')

    #     STORAGE.update(ID, {'mydata1': self.data['mydata1']})
    #     record = STORAGE.get(ID)

    #     if record.data['mydata1'] != self.data['mydata1']:
    #         raise Exception('Put_field failed ...')

    #     # Check right
    #     with self.assertRaises(KeyError):
    #         STORAGE.update(ID, {'mydata1': 'toto'}, account=self.user2_account)

    # def test_06_Get_field(self):

    #     raw_record = STORAGE.get(ID, mfields=["mydata2"])
    #     print("raw_record: %s" % raw_record)

    #     raw_record["mydata2"]

    #     if len(raw_record) != 2:  # field + _id
    #         raise Exception('Get_field failed ...')

    # def test_07_MultiGet(self):
    #     record1 = Record({'check': 'remove1'})
    #     record2 = Record({'check': 'remove2'})
    #     record3 = Record({'check': 'remove3'})

    #     ids = STORAGE.put([record1, record2, record3])
    #     records = STORAGE.get(ids)

    #     if len(records) != 3:
    #         print(records)
    #         raise Exception("Impossible to get with id's list")

    #     STORAGE.remove(ids)

    # def test_08_UpdateAndPut(self):
    #     MYRECORD.data['mydata4'] = 'data4'
    #     STORAGE.put(MYRECORD)
    #     record = STORAGE.get(ID)
    #     if record.data == self.data:
    #         print("record.data: %s" % record.data)
    #         print("self.data:   %s" % self.data)
    #         raise Exception('Data not updated ...')

    # def test_09_Remove(self):
    #     record1 = Record({'check': 'remove1'})
    #     id1 = STORAGE.put(record1)
    #     record2 = Record({'check': 'remove2'})
    #     id2 = STORAGE.put(record2)
    #     record3 = Record({'check': 'remove3'})
    #     id3 = STORAGE.put(record3)

    #     STORAGE.remove([id1, id2, id3])

    #     STORAGE.remove(ID)

    def test_10_CheckRemove(self):
        self.assertRaises(KeyError, STORAGE.get, ID)

    def test_11_ManyInsert(self):
        record1 = Record({'check': 'test1', 'state': 1})
        record2 = Record({'check': 'test2', 'state': 0})
        record3 = Record({'check': 'test3', 'state': 0})

        STORAGE.put([record1, record2, record3])

    # def test_12_Find(self):
    #     records = STORAGE.find({'check': 'test1'})
    #     #for record in records:
    #     #   record.cat()

    #     if len(records) != 1:
    #         raise Exception('Error in filter ...')

    #     records = STORAGE.find({}, mfields=['state'])
    #     if len(records) != 3:
    #         raise Exception('Error in mfields ...')

    # def test_13_Find_limit(self):
    #     records = STORAGE.find({}, limit=3)

    #     if len(records) != 3:
    #         raise Exception('Error in limit ...')

    # def test_14_FindOne(self):
    #     record = STORAGE.find_one({'check': 'test1'})

    #     if not isinstance(record, Record):
    #         raise Exception('Error in find_one ...')

    # def test_15_count(self):
    #     nb = STORAGE.count({'state': 0})
    #     if nb != 2:
    #         raise Exception('Error in count ...')

    # def test_16_CheckReadRights(self):
    #     # Inserts
    #     STORAGE.put(Record({'check': 'test4'}), account=self.anonymous_account)
    #     STORAGE.put(Record({'check': 'test5'}), account=self.anonymous_account)
    #     _id = STORAGE.put(Record({'check': 'test6'}), account=self.root_account)

    #     ## 3 records for user
    #     ## 2 records for anonymous
    #     ## 6 records for root

    #     records = STORAGE.find(account=self.user_account)
    #     if len(records) != 3:
    #         raise Exception('Invalid rigths for user account ...')

    #     records = STORAGE.find(account=self.anonymous_account)
    #     if len(records) != 2:
    #         raise Exception('Invalid rigths for anonymous account ...')

    #     self.assertRaises(KeyError, STORAGE.get, _id, self.anonymous_account)

    #     record = STORAGE.get(_id, account=self.root_account)

    #     records = STORAGE.find(account=self.root_account)
    #     if len(records) != 6:
    #         raise Exception('Invalid rigths for root account ...')

    # def test_17_CheckWriteRights(self):
    #     # Insert with user account
    #     record = Record({'check': 'test7'})
    #     STORAGE.put(record, account=self.user_account)

    #     ## try to remove with anonymous account
    #     self.assertRaises(
    #         ValueError, STORAGE.remove, record, self.anonymous_account)

    #     ## Change rights
    #     record.chgrp('group.anonymous')
    #     record.chmod('g+w')
    #     STORAGE.put(record)

    #     ## try to remove with anonymous account
    #     STORAGE.remove(record, account=self.anonymous_account)

    # def test_18_MapReduce(self):
    #     from bson.code import Code

    #     mmap = Code("function () {"
    #     "       if (this.state == 0){ emit('ok', 1) }"
    #     "       else if (this.state == 1){ emit('warning', 1) }"
    #     "       else if (this.state == 2){ emit('critical', 1) }"
    #     "       else if (this.state == 3){ emit('unknown', 1) }"
    #     "}")

    #     mreduce = Code("function (key, values) {"
    #     "  var total = 0;"
    #     "  for (var i = 0; i < values.length; i++) {"
    #     "    total += values[i];"
    #     "  }"
    #     "  return total;"
    #     "}")

    #     result = STORAGE.map_reduce({}, mmap, mreduce)

    #     if result['ok'] != 2 and result['warning'] != 1:
    #         raise Exception('Invalid map/reduce result ...')

    # def test_19_tree(self):
    #     record1 = Record({'data': 1}, name="record1")
    #     record2 = Record({'data': 2}, name="record2")
    #     record3 = Record({'data': 3}, name="record3")
    #     record4 = Record({'data': 4}, name="record4")

    #     STORAGE.put([record1, record2, record3, record4])

    #     record2.add_children(record4)

    #     record1.add_children(record2)
    #     record1.add_children(record3)

    #     STORAGE.put([record1, record2])
    #     STORAGE.get_record_childs(record1)
    #     STORAGE.recursive_get(record1)

    #     STORAGE.print_record_tree(record1)

    #     dumps(record1.dump(json=True))

    # def test_20_GetRecursivelyAllChilds(self):
    #     parent_record = STORAGE.find_one({'data': 1})
    #     STORAGE.recursive_get(parent_record)
    #     if len(parent_record.children) != 2:
    #         raise Exception('Wrong parent/children tree builded')

    def test_21_RemoveAll(self):
        records = STORAGE.find(account=self.root_account)
        STORAGE.remove(records, account=self.root_account)

    def test_99_DropNamespace(self):
        STORAGE.drop_namespace('unittest')

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
