#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest

from canopsis.common import root_path
from canopsis.healthcheck.manager import HealthcheckManager, check_checkable
from canopsis.logger.logger import Logger, OutputNull
from canopsis.models.healthcheck import Healthcheck, OK_MSG
import xmlrunner

# TODO make these unit tests independant of rabbitmq or move them into the
# functional_testing directory
class HealthcheckTest(unittest.TestCase):
    pass
    # def setUp(self):
    #     logger = Logger.get('', None, output_cls=OutputNull)
    #     self.manager = HealthcheckManager(logger)

    # def test_healthcheck(self):
    #     check = self.manager.check()
    #     for service in Healthcheck.SERVICES:
    #         self.assertIn(service, check)
    #         if service != 'engines':
    #             # No engine test on dockerised env
    #             self.assertEqual(check[service], OK_MSG)
    #     self.assertIn(Healthcheck.TIME, check)
    #     self.assertTrue(check['overall'])

    # def test_check_checkable(self):
    #     #check = check_checkable("canopsis-engine@")
    #     #self.assertTrue(check)
    #     # ^ cannot do that on dockerized env

    #     check = check_checkable("Shinmen Takezō")
    #     self.assertFalse(check)

    # def test_check_rabbitmq_state(self):
    #     check = self.manager._check_rabbitmq_state()
    #     self.assertEqual(check.message, OK_MSG)
    #     self.assertTrue(check.state)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
