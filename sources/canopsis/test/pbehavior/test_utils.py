from unittest import TestCase, main

from canopsis.pbehavior.utils import check_valid_rrule
import unittest
from canopsis.common import root_path
import xmlrunner

class UtilsTest(TestCase):

    VALID_RRULE = 'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU'
    INVALID_RRULE = 'InVaLiDrRuLe'

    def test_check_valid_rrule(self):
        res = check_valid_rrule(self.VALID_RRULE)

        self.assertTrue(res)

    def test_check_invalid_rrule(self):
        with self.assertRaises(ValueError):
            check_valid_rrule(self.INVALID_RRULE)

if __name__ == '__main__':
    output = root_path + "/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
