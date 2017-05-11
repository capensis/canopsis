#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from datetime import datetime, timedelta
import time
from unittest import main

from canopsis.alerts.filter import AlarmFilter

from base import BaseTest


class TestFilter(BaseTest):
    def test_check_alarm(self):
        now = datetime.now() - timedelta(minutes=31)
        now_stamp = int(time.mktime(now.timetuple()))
        alarm, value = self.gen_fake_alarm(now_stamp)

        lifter = AlarmFilter({
            'operator': 'eq',
            'key': 'cacao',
            'value': 'maigre'
        })
        self.assertFalse(lifter.check_alarm(value))

        lifter = AlarmFilter({
            'operator': 'eq',
            'key': 'component',
            'value': 'bbb'
        })
        self.assertFalse(lifter.check_alarm(value))

        lifter = AlarmFilter({
            'operator': 'eq',
            'key': 'component',
            'value': 'c'
        })
        self.assertTrue(lifter.check_alarm(value))


if __name__ == '__main__':
    main()
