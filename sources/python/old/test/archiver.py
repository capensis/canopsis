#!/usr/bin/env python
# --------------------------------
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
from logging import INFO

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

def setFields(_map, **kwargs):
    for key in kwargs:
        _map[key] = kwargs[key]

class KnownValues(unittest.TestCase):
    def setUp(self):
        self.Archiver = Archiver(namespace='unittest',
                                 autolog=True,
                                 logging_level=loggingDEBUG)
        self.Archiver.beat()

    def test_01_check_statuses(self):

        devent = {'rk': 'test_03_check_statuses',
                  'status': 0,
                  'timestamp': 14389,
                  'state': 0}

        event = {'rk': 'test_03_check_statuses',
                 'status': 0,
                 'timestamp': 14400,
                 'state': 0}

        # Check that event stays off even if it appears
        # more than the bagot freq in the stealthy/bagot interval
        for x in xrange(1, 50):
            self.Archiver.check_statuses(event, devent)
            devent = event.copy()
            setFields(event, timestamp=(event['timestamp']+1))
            self.assertEqual(event['status'], OFF)

        # Set state to alarm, event should be On Going
        setFields(event, state=1)
        self.Archiver.check_statuses(event, devent)
        self.assertEqual(event['status'], ONGOING)
        devent = event.copy()

        # Set state back to Ok, event should be Stealthy
        setFields(event, state=0)
        self.Archiver.check_statuses(event, devent)
        self.assertEqual(event['status'], STEALTHY)
        devent = event.copy()

        # Move TS out of stealthy range, event should be On Going
        setFields(event, state=1, timestamp=event['timestamp']+1000)
        self.Archiver.check_statuses(event, devent)
        self.assertEqual(event['status'], ONGOING)
        devent = event.copy()

        # Check that the event is at Bagot when the requirments are met
        for x in xrange(1, 14):
            if x % 2:
                setFields(event, state=0 if event['state'] else 1);
            self.Archiver.check_statuses(event, devent)
            setFields(event, timestamp=(event['timestamp']+1))
            if devent['bagot_freq'] >= self.Archiver.bagot_freq:
                self.assertEqual(event['status'], BAGOT)
            devent = event.copy()


        # Check that the event is On Going if out of the Bagot time interval
        setFields(event, state=1, timestamp=event['timestamp']+4000)
        self.Archiver.check_statuses(event, devent)
        self.assertEqual(event['status'], ONGOING)
        devent = event.copy()



if __name__ == "__main__":
    unittest.main(verbosity=2)


