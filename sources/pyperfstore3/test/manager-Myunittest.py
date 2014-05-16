import unittest

from pyperfstore3.manager import Manager
from pyperfstore3.timewindow import Period, TimeWindow
import logging


class ManagerTest(unittest.TestCase):

	def setUp(self):
		logger = logging.getLogger()
		logger.setLevel("DEBUG")
		self.manager = Manager(logger=logger)

	def test_put_get_data(self):

		timewindow = TimeWindow()

		metric_id = 'test_manager'

		self.manager.remove(metric_id=metric_id)

		tv0 = (int(timewindow.start()), None)
		tv1 = (int(timewindow.start() + 1), 0)
		tv2 = (int(timewindow.stop()), 2)
		tv3 = (int(timewindow.stop() + 1000000), 3)

		# set values with timestamp without order
		points = (tv0, tv2, tv1, tv3)

		meta = {'plop': None}

		period = Period(minute=60)

		self.manager.put_data(
			metric_id=metric_id,
			points=points,
			meta=meta,
			period=period)

		data, _meta = self.manager.get_data(
			metric_id=metric_id,
			timewindow=timewindow,
			period=period,
			return_meta=True)

		self.assertEqual(meta, _meta)

		self.assertEqual((tv0, tv1, tv2), data)

		# remove 1 data at stop point
		_interval = (timewindow.stop(), timewindow.stop())
		_timewindow = TimeWindow(_interval)

		self.manager.remove(metric_id, _timewindow, period=period)

		data, _meta = self.manager.get_data(
			metric_id=metric_id,
			timewindow=timewindow,
			period=period,
			return_meta=True)

		self.assertEqual(meta, _meta)

		self.assertEqual(data, (tv0, tv1))

		# get data on timewindow
		data, _meta = self.manager.get_data(
			metric_id=metric_id,
			timewindow=timewindow,
			period=period,
			return_meta=True)

		self.assertEqual(meta, _meta)

		# get all data
		data, _meta = self.manager.get_data(
			metric_id=metric_id,
			period=period,
			return_meta=True)

		self.assertEqual(meta, _meta)

		self.assertEqual(len(data), 3)

		# remove all data
		self.manager.remove(
			metric_id=metric_id)

		data, _meta = self.manager.get_data(
			metric_id=metric_id,
			period=period,
			return_meta=True)

		self.assertEqual(None, _meta)

		self.assertEqual(len(data), 0)

if __name__ == '__main__':
	unittest.main()
