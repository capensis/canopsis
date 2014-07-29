#!/usr/bin/env python
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
from logging import DEBUG as loggingDEBUG, basicConfig as loggingBasicConfig

from canopsis.old.archiver import Archiver

ARCHIVER = None

loggingBasicConfig(
    level=loggingDEBUG,
    format='%(asctime)s %(name)s %(levelname)s %(message)s',
    )

#Statuses
OFF = 0
ONGOING = 1
STEALTHY = 2
BAGOT = 3
CANCELED = 4

class KnownValues(unittest.TestCase):
    def setUp(self):
        self.Archiver = Archiver(namespace='unittest',  autolog=True, logging_level=loggingDEBUG)
        self.Archiver.bagot_freq = 10
        self.Archiver.bagot_time = 3600
        self.Archiver.furtif_time = 300
        self.Archiver.restore_event = True

    def test_01_cancel(self):
        event = {'cancel': True, 'output': 'output', 'author': 'eric', 'state': 1}
        devent = {'output': 'test', 'status': STEALTHY, 'state' : 0}

        self.Archiver.process_cancel(devent, event)
        #self.assertTrue('cancel' not in event)


        event = {'cancel': True, 'output': 'output', 'author': 'eric', 'state': 0}
        devent = {'output': 'test', 'status': STEALTHY, 'state': 0}

        self.Archiver.process_cancel(devent, event)

        #Cancellation
        self.assertTrue(event['cancel']['cancel'])
        self.assertEqual(event['cancel']['comment'], 'output')
        self.assertEqual(event['cancel']['previous_status'], STEALTHY)
        self.assertEqual(event['cancel']['author'], 'eric')

    def test_02_uncancel(self):
        #No previous cancellation done
        event = {'cancel': False, 'output': 'output', 'author': 'eric', 'state': 0}
        devent = {'output': 'test', 'status': STEALTHY, 'state': 0}

        self.Archiver.process_cancel(devent, event)
        self.assertTrue('cancel' not in event)

        #Previous cancellation existed
        event = {'cancel': False, 'output': 'output', 'author': 'eric', 'state': 0}
        devent = {
            'output': 'test',
            'status': STEALTHY,
            'cancel': {
                'cancel': True,
                'previous_status': ONGOING,
                'previous_state': 1
            },
            'state': 0}

        self.Archiver.process_cancel(devent, event)
        #Cancellation
        self.assertEqual(event['cancel']['comment'], 'output')
        self.assertEqual(event['cancel']['previous_status'], STEALTHY)
        self.assertEqual(event['cancel']['author'], 'eric')

        #Previous status set to Ok by default
        event = {'cancel': False, 'output': 'output', 'author': 'eric', 'state': 0}
        devent = {
            'output': 'test',
            'status': ONGOING,
            'cancel': {
                'cancel': True,
                'previous_status': 1,
                'previous_state': 1
            },
            'state': 0}


        self.Archiver.process_cancel(devent, event)
        self.assertEqual(event['cancel']['previous_status'], 1)

    def test_03_check_bagot(self):

        event = {'rk':'testrk' , 'timestamp' : 1}
        devent = {'timestamp' : 1}
        self.Archiver.check_bagot(event, devent)

        self.assertEqual(event['furtif_freq'], 1)
        self.assertEqual(event['last_furtif'], 1)
        self.assertEqual(event['status'], STEALTHY)

        for x in xrange(1,10):
            event['timestamp'] += 1
            self.Archiver.check_bagot(event, devent)
            self.assertEqual(event['furtif_freq'], x+1)
            self.assertEqual(event['last_furtif'], x+1)
            self.assertEqual(event['timestamp'], x+1)
            if x == 9:
                self.assertEqual(event['status'], BAGOT)
            else:
                self.assertEqual(event['status'], STEALTHY)


    def test_04_check_statuses(self):
        event = {'rk': 'testrk',
                 'timestamp': 1,
                 'state': 0}

        devent = {'rk': 'testrk',
                  'timestamp': 1,
                  'state': 0}

        # Check that event stays off even if it appears
        # more than the bagot freq in the stealthy/bagot interval
        for x in xrange(1, 20):
            self.Archiver.check_statuses(event, devent)
            event = devent

        self.assertEqual(event['status'], OFF)

        # Check that the event becomes Stealthy and then Bagot
        for x in xrange(1, 20):
            if (x % 2):
                devent['state'] = 1
            self.Archiver.check_statuses(event, devent)
            event = devent
            if x >= 10:
                self.assertEqual(event['status'], BAGOT)
            elif x > 3:
                self.assertEqual(event['status'], STEALTHY)

        # Check that the event is On Going if out of the Bagot time interval
        event = {'rk': 'testrk',
                 'timestamp': 4000,
                 'state': 2}
        self.Archiver.check_statuses(event, devent)
        self.assertEqual(event['status'], ONGOING)
        devent = event

        # Check that the event is now Off if out of the Stealthy time interval
        event = {'rk': 'testrk',
                 'timestamp': 4500,
                 'state': 0}
        self.Archiver.check_statuses(event, devent)
        self.assertEqual(event['status'], OFF)



if __name__ == "__main__":
    unittest.main(verbosity=2)


