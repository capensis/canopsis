import unittest

from pyperfstore3.timeserie import TimeSerie


class TimeSerieTest(unittest.TestCase):
	"""
	UT on timeserie
	"""
	def setUp(self):

		self.timeserie = TimeSerie()

	pass

if __name__ == '__main__':
	unittest.main()
