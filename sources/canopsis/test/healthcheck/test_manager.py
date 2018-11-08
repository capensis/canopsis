#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest

from canopsis.common import root_path
from canopsis.healthcheck.manager import HealthcheckManager
from canopsis.logger.logger import Logger, OutputNull
from canopsis.models.healthcheck import Healthcheck, OK_MSG
import xmlrunner


class HealthcheckTest(unittest.TestCase):

    def setUp(self):
        logger = Logger.get('', None, output_cls=OutputNull)
        self.manager = HealthcheckManager(logger)

    def test_healthcheck(self):
        check = self.manager.check()
        self.assertEqual(check["overall"], True)
        for service in Healthcheck.SERVICES:
            self.assertTrue(service in check)
            self.assertEqual(check[service], OK_MSG)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
