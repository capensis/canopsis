# -*- coding: utf-8 -*-

import unittest

from canopsis.activity.pbehavior import TimeDuration


class TestTimeDuration(unittest.TestCase):

    def test_eq(self):

        td1 = TimeDuration(1, 2, 3, 10)
        td2 = TimeDuration(1, 2, 3, 10)
        td3 = TimeDuration(1, 2, 3, 9)
        td4 = TimeDuration(1, 2, 4, 10)
        td5 = TimeDuration(1, 3, 3, 10)
        td6 = TimeDuration(2, 2, 3, 10)

        self.assertEqual(td1, td2)
        self.assertNotEqual(td1, td3)
        self.assertNotEqual(td1, td4)
        self.assertNotEqual(td1, td5)
        self.assertNotEqual(td1, td6)

    def test_str(self):

        td = TimeDuration(1, 2, 3, 4)
        tds = str(td)

        self.assertEqual('1:2:3/4', tds)

    def test_hash(self):

        td = TimeDuration(1, 2, 3, 4)
        tdS = set([td])

        tdv = tdS.pop()

        self.assertEqual(tdv, td)


if __name__ == '__main__':
    unittest.main()
