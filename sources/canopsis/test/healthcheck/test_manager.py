#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest

from canopsis.common import root_path
from canopsis.healthcheck.manager import HealthcheckManager, check_checkable
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
        self.assertTrue(Healthcheck.TIME in check)

    def test_check_checkable(self):
        #check = check_checkable("canopsis-engine@")
        #self.assertEqual(check, True)
        # ^ cannot do that on dockerized env

        check = check_checkable("Shinmen Takezō")
        self.assertEqual(check, False)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
