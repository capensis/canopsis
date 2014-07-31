import sys

sys.path.append("../pyperfstore2/")

import unittest

import datetime, time

from utils import MN, HR, D, W, M, Y, roundTime, getTimeSteps

class AggregationTest(unittest.TestCase):

	def setUp(self):		
		self.date = datetime.datetime.utcnow()
		print 'test date : %s' % self.date 

	def testRoundTime(self):		
		print 'MN ', roundTime(self.date, MN)

		print '5MN ', roundTime(self.date, 5*MN)

		print '15MN ', roundTime(self.date, 15*MN)

		print '30MN ', roundTime(self.date, 30*MN)

		print 'HR ', roundTime(self.date, HR)

		print 'D ', roundTime(self.date, D)

		print 'W ', roundTime(self.date, W)

		print 'M ', roundTime(self.date, M)

		print 'Y ', roundTime(self.date, Y)

		pass

	def testTimeSteps(self):
		
		t_date = time.time()		

		stop_date = t_date

		start_date = stop_date - 60 * 15
		timeSteps = getTimeSteps(start_date, stop_date, 5*MN)
		print '5MN', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]

		start_date = stop_date - 60 * 30
		timeSteps = getTimeSteps(start_date, stop_date, 15*MN)
		print '15MN', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]
		
		start_date = stop_date - 60 * 60
		timeSteps = getTimeSteps(start_date, stop_date, 30*MN)
		print '30MN', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]

		start_date = stop_date - 60 * 60 * 5
		timeSteps = getTimeSteps(start_date, stop_date, HR)
		print 'HR', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]

		start_date = stop_date - 60 * 60 * 48
		timeSteps = getTimeSteps(start_date, stop_date, D)
		print 'D', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]

		start_date = stop_date - 60 * 60 * 24 * 20
		timeSteps = getTimeSteps(start_date, stop_date, W)
		print 'W', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]

		start_date = stop_date - 60 * 60 * 24 * 700
		timeSteps = getTimeSteps(start_date, stop_date, Y)
		print 'Y', [ datetime.datetime.fromtimestamp(ts) for ts in timeSteps ]


		pass

	def testAggregate(self):
		pass

if __name__ == "__main__":
	unittest.main()