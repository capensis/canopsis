from unittest import TestCase

from canopsis.pbehavior.utils import check_valid_rrule

class UtilsTest(TestCase):

    VALID_RRULE = 'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU'
    INVALID_RRULE = 'InVaLiDrRuLe'

    def test_check_valid_rrule(self):
        res = check_valid_rrule(self.VALID_RRULE)

        self.assertTrue(res)

    def test_check_invalid_rrule(self):
        with self.assertRaises(ValueError):
            check_valid_rrule(self.INVALID_RRULE)