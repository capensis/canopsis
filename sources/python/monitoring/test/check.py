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

from canopsis.monitoring.check import CheckRunner
from copy import deepcopy
import socket


CONTEXT = {
    'connector': 'unittest',
    'connector_name': 'monitoring_check',
    'source_type': 'resource',
    'component': 'monitoring',
    'resource': 'check_test'
}


SCHEMA = {
    'meta': {
        'command': {
            'binpath': '/test/check_test',
            'args': {
                'O': 'optional',
                'B': 'boolvar',
                'S': 'strvar',
                'I': 'intvar'
            }
        }
    },
    'type': 'object',
    'properties': {
        'optional': {'type': 'string', 'required': False},
        'boolvar': {'type': 'boolean', 'required': False, "default": False},
        'strvar': {'type': 'string', 'required': True},
        'intvar': {'type': 'integer', 'required': True}
    }
}


class KnownValues(TestCase):
    def setUp(self):
        self.runner = CheckRunner(CONTEXT, 'test')
        self.runner.hostname = 'monitoring'

    def _validate(
        self,
        reverse,
        cmd, cmdargs,
        expected_cmd,
        expected_cmdargs, not_expected_cmdargs
    ):
        cmdargs = ' '.join(cmdargs)

        self.assertEqual(cmd, expected_cmd)

        assertTrue = self.assertTrue if not reverse else self.assertFalse
        assertFalse = self.assertFalse if not reverse else self.assertTrue

        for arg in expected_cmdargs:
            assertTrue(arg in cmdargs)

        for arg in not_expected_cmdargs:
            assertFalse(arg in cmdargs)

    def validate(self, *args, **kwargs):
        self._validate(False, *args, **kwargs)

    def dont_validate(self, *args, **kwargs):
        self._validate(True, *args, **kwargs)

    def test_build_command_all_true(self):
        conf = {
            'optional': 'test',
            'boolvar': True,
            'strvar': 'test',
            'intvar': 4
        }

        cmd, cmdargs = self.runner.build_command(SCHEMA, conf)

        self.validate(
            cmd, cmdargs,
            '/test/check_test',
            ['-O test', '-B', '-S test', '-I 4'],
            []
        )

    def test_build_command_all_false(self):
        conf = {
            'optional': 'test',
            'strvar': 'test',
            'intvar': 4
        }

        cmd, cmdargs = self.runner.build_command(SCHEMA, conf)

        self.validate(
            cmd, cmdargs,
            '/test/check_test',
            ['-O test', '-S test', '-I 4'],
            ['-B']
        )

    def test_build_command_without_optional(self):
        conf = {
            'boolvar': True,
            'strvar': 'test',
            'intvar': 4
        }

        cmd, cmdargs = self.runner.build_command(SCHEMA, conf)

        self.validate(
            cmd, cmdargs,
            '/test/check_test',
            ['-B', '-S test', '-I 4'],
            ['-O test']
        )

    def test_build_command_failure(self):
        conf = {
            'boolvar': True,
            'strvar': 'test',
            'intvar': 4
        }

        cmd, cmdargs = self.runner.build_command(SCHEMA, conf)

        self.dont_validate(
            cmd, cmdargs,
            '/test/check_test',
            ['-O test'],
            ['-B', '-S test', '-I 4']
        )

    def test_gen_event(self):
        expected = deepcopy(CONTEXT)

        expected['event_type'] = 'check'
        expected['output'] = 'test OK: test'
        expected['long_output'] = "there was a test\nand it went good."

        expected['perf_data_array'] = [
            {'metric': 'size', 'min': 0.0, 'unit': 'B', 'value': 2102.0}
        ]

        expected['state'] = 0
        expected['state_type'] = 1

        evt = self.runner.gen_event(0, "test OK: test |size=2102B;;;0\nthere was a test\nand it went good.\n")

        self.assertEqual(evt, expected)


if __name__ == '__main__':
    main()
