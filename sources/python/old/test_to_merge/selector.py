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
import time

from canopsis.event.selector import Selector
#from Selector import cselector_get, cselector_getall
import logging

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s %(name)s %(levelname)s %(message)s',
                    )

from canopsis.old.record import Record
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

root_account = Account(user="root", group="root")
storage = get_storage(
    account=root_account, namespace='unittest', logging_level=logging.DEBUG)
selector = None


class KnownValues(unittest.TestCase):

    def test_01_InitPutGet(self):
        global selector
        selector = Selector(
            name="myselector", namespace='unittest', storage=storage)
        selector.nocache = True

        _id = selector._id
        print("Selector Id: %s" % _id)

        dump = selector.dump()

        ## Put in db
        storage.put(selector)

        ## Load
        selector = Selector(name="myselector", storage=storage)
        ndump = selector.dump()

        if dump['_id'] != ndump['_id'] or dump['namespace'] != ndump['namespace'] or dump['mfilter'] != ndump['mfilter'] or dump['aaa_owner'] != ndump['aaa_owner']:
            print(dump)
            print(ndump)
            raise Exception('Invalid dump ...')

    def test_02_PutData(self):
        record1 = Record({'_id': 'check1', 'check': 'test1', 'state': 0, 'source_type': 'service'})
        record2 = Record({'_id': 'check2', 'check': 'test2', 'state': 1, 'source_type': 'service'})
        record3 = Record({'_id': 'check3', 'check': 'test3', 'state': 2, 'source_type': 'service', 'state_type': 0, 'previous_state': 1})

        storage.put([record1, record2, record3])

    def test_03_Resolv(self):

        selector.setMfilter({'_id': 'check1'})
        ids = selector.resolv()

        if len(ids) != 1:
            raise Exception('Invalid count (%s)' % ids)

        selector.setMfilter({'state': 0})
        ids = selector.resolv()

        if len(ids) != 1:
            raise Exception('Invalid count (%s)' % ids)

        selector.setMfilter({'$or': [{'check': 'test1'}, {'check': 'test2'}]})
        ids = selector.resolv()

        if len(ids) != 2:
            raise Exception('Invalid count (%s)' % ids)

        ## Check cache
        selector.cache_time = 2

        ids = selector.resolv()
        if len(ids) != 2:
            raise Exception('Invalid count (%s)' % ids)

        selector.mfilter = {'_id': 'check1'}
        ids = selector.resolv()
        if len(ids) != 2:
            raise Exception('Invalid count with cache (%s)' % ids)

        time.sleep(3)
        ids = selector.resolv()
        if len(ids) != 1:
            raise Exception('Invalid count with cache (%s)' % ids)

    def test_04_ResolvIncludeExclude(self):
        selector.setMfilter(
            {'$or': [{'check': 'test1'}, {'check': 'test2'}, {'check': 'test3'}]})

        selector.setExclude_ids(['check1'])
        ids = selector.resolv()
        if len(ids) != 2:
            raise Exception('Invalid count with Exclude (%s)' % ids)

        selector.setExclude_ids([])
        selector.setInclude_ids([])

        selector.setMfilter({'$or': [{'check': 'test1'}, {'check': 'test2'}]})

        selector.setInclude_ids(['check3'])
        ids = selector.resolv()
        if len(ids) != 3:
            raise Exception('Invalid count with Include (%s)' % ids)

        selector.setExclude_ids(['check3'])
        ids = selector.resolv()
        if len(ids) != 2:
            raise Exception('Invalid count with Include + Exclude (%s)' % ids)

        selector.setExclude_ids([])
        selector.setInclude_ids([])

    def test_05_Match(self):
        selector.setMfilter({'$or': [{'check': 'test1'}, {'check': 'test2'}]})

        if not selector.match({'check': 'test1'}):
            raise Exception('Error in match, plain id ...')

        if selector.match({'check': 'toto'}):
            raise Exception('Error in match, wrong id ...')

    def test_06_GetRecords(self):
        selector.setMfilter({'$or': [{'check': 'test1'}, {'check': 'test2'}]})
        records = selector.getRecords()

        if len(records) != 2 and isinstance(records[0], Record):
            raise Exception('Error in get records ("%s")' % records)

    def test_06_GetState(self):
        selector.setMfilter(
            {'$or': [{'check': 'test1'}, {'check': 'test2'}, {'check': 'test3'}]})
        (states, state, state_type) = selector.getState()
        if state != 1:
            raise Exception('Invalid state ("%s")' % state)

        selector.setMfilter({'$or': [{'check': 'test1'}, {'check': 'test2'}]})
        (states, state, state_type) = selector.getState()
        if state != 1:
            raise Exception('Invalid state ("%s")' % state)

        selector.setMfilter(
            {'$or': [{'check': 'test1'}, {'check': 'test2'}, {'check': 'test3'}]})
        (states, state, state_type) = selector.getState()
        if state != 1:
            raise Exception('Invalid state ("%s")' % state)

        selector.setMfilter({'$or': [{'check': 'test1'}]})
        (states, state, state_type) = selector.getState()
        if state != 0:
            raise Exception('Invalid state ("%s")' % state)

    def test_08_Remove(self):
        storage.remove(selector)

    def test_99_DropNamespace(self):
        storage.drop_namespace('unittest')


if __name__ == "__main__":
    unittest.main(verbosity=2)
