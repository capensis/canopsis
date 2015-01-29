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

from canopsis.monitoring.parser import CheckParser, PerfDataParser


TEST_CHECK_OUTPUT = """DISK OK - free space: / 3326 MB (56%); | /=2643MB;5948;5958;0;5968
/ 15272 MB (77%);
/boot 68 MB (69%);
/home 69357 MB (27%);
/var/log 819 MB (84%); | /boot=68MB;88;93;0;98
/home=69357MB;253404;253409;0;253414 
/var/log=818MB;970;975;0;980"""

TEST_PERFDATA = "/=2643MB;5948;5958;0;5968,/boot=68MB;88;93;0;98,/home=69357MB;253404;253409;0;253414,/var/log=818MB;970;975;0;980"


class KnownValues(TestCase):
    def test_check_output(self):
        parser = CheckParser(0, TEST_CHECK_OUTPUT)

        expected_status = 0
        expected_text = "DISK OK - free space: / 3326 MB (56%);"
        expected_long_output = """/ 15272 MB (77%);
/boot 68 MB (69%);
/home 69357 MB (27%);
/var/log 819 MB (84%);"""

        expected_perfdata = TEST_PERFDATA

        self.assertEqual(parser.status, expected_status)
        self.assertEqual(parser.text, expected_text)
        self.assertEqual(parser.long_output, expected_long_output)
        self.assertEqual(parser.perfdata, expected_perfdata)

    def test_perfdata(self):
        parser = PerfDataParser(TEST_PERFDATA)

        expected = [
            {
                'crit': 5958.0,
                'max': 5968.0,
                'metric': '/',
                'min': 0.0,
                'unit': 'MB',
                'value': 2643.0,
                'warn': 5948.0
            },{
                'crit': 93.0,
                'max': 98.0,
                'metric': '/boot',
                'min': 0.0,
                'unit': 'MB',
                'value': 68.0,
                'warn': 88.0
            },{
                'crit': 253409.0,
                'max': 253414.0,
                'metric': '/home',
                'min': 0.0,
                'unit': 'MB',
                'value': 69357.0,
                'warn': 253404.0
            },{
                'crit': 975.0,
                'max': 980.0,
                'metric': '/var/log',
                'min': 0.0,
                'unit': 'MB',
                'value': 818.0,
                'warn': 970.0
            }
        ]

        self.assertEqual(expected, parser.perf_data_array)


if __name__ == "__main__":
    main()
