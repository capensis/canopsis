import unittest

from pyperfstore3.store import TimedStore
from pyperfstore3.timewindow import TimeWindow


class TimedStoreTest(unittest.TestCase):
	"""
	pyperfstore3.store.TimedStore UT on data_name = "test_store"
	"""

	def setUp(self):
		# create a store on test_store collections
		self.store = TimedStore(data_name="test_store")

	def test_connect(self):
		self.assertTrue(self.store.connected())

		self.store.disconnect()

		self.assertFalse(self.store.connected())

		self.store.connect()

		self.assertTrue(self.store.connected())

	def test_CRUD(self):

		data_id = 'test_store_id'

		# start in droping data
		self.store.drop()

		# ensure count is 0
		count = self.store.count(data_id=data_id)
		self.assertEquals(count, 0)

		# let's play with different data_names
		meta = {'min': None, 'max': 0}

		timewindow = TimeWindow()

		before_timewindow = [timewindow.start() - 1000]
		in_timewindow = [
			timewindow.start(),
			timewindow.start() + 5,
			timewindow.stop() - 5,
			timewindow.stop()]
		after_timewindow = [timewindow.stop() + 1000]

		# set timestamps without sort
		timestamps = after_timewindow + before_timewindow + in_timewindow

		for timestamp in timestamps:
			# add a document at starting time window
			self.store.put(data_id=data_id, value=meta, timestamp=timestamp)

		# check for count equals 5
		count = self.store.count(data_id=data_id)
		self.assertEquals(count, 2)

		# check for_data before now
		data = self.store.get(data_id=data_id)
		self.assertEquals(len(data), 1 if len(in_timewindow) > 0 else 0)

		# check for data inside timewindow and just before
		data = self.store.get(data_id=data_id, timewindow=timewindow)
		self.assertEquals(len(data), 1)

		# remove data inside timewindow
		self.store.remove(data_id=data_id, timewindow=timewindow)
		# check for data outside timewindow
		count = self.store.count(data_id=data_id)
		self.assertEquals(count, len(before_timewindow) + len(after_timewindow))

		# remove all data
		self.store.remove(data_id=data_id)
		# check for count equals 0
		count = self.store.count(data_id=data_id)
		self.assertEquals(count, 0)

from pyperfstore3.store import PeriodicStore
from pyperfstore3.timewindow import Period


class PeriodicStoreTest(unittest.TestCase):

	def setUp(self):
		# create a store on test_store collections
		self.store = PeriodicStore(data_name="test_store")

	def test_connect(self):
		self.assertTrue(self.store.connected())

		self.store.disconnect()

		self.assertFalse(self.store.connected())

		self.store.connect()

		self.assertTrue(self.store.connected())

	def test_CRUD(self):
		# start in droping data
		self.store.drop()

		# let's play with different data_names
		data_ids = ['m0', 'm1']
		aggregations = ['mean', 'max', '']
		periods = [Period(**{Period.MINUTE: 60}),
			Period(**{Period.HOUR: 24})]

		timewindow = TimeWindow()

		points = [
			(timewindow.start(), None),  # lower bound
			(timewindow.stop(), 0),  # upper bound
			(timewindow.start() - 1, 1),  # outside timewindow (minus 1)
			(timewindow.start() + 1, 2),  # inside timewindow (plus 1)
			(timewindow.stop() + 1, 3)  # outside timewindow (plus 1)
		]

		sorted_points = sorted(points, key=lambda point: point[0])

		inserted_points = dict()

		# starts to put points for every aggregations and periods
		for data_id in data_ids:
			inserted_points[data_id] = dict()
			for aggregation in aggregations:
				inserted_points[data_id][aggregation] = dict()
				for period in periods:
					inserted_points[data_id][aggregation][period] = points
					# add documents
					self.store.put(data_id=data_id, aggregation=aggregation,
						period=period, points=points)

		points_count_in_timewindow = len(
						[point for point in points if point[0] in timewindow])

		# check for reading methods
		for data_id in data_ids:
			# iterate on data_ids

			for aggregation in aggregations:

				for period in periods:

					count = self.store.count(data_id=data_id, aggregation=aggregation,
						period=period)
					self.assertEquals(count, len(points))

					count = self.store.count(data_id=data_id, aggregation=aggregation,
						period=period, timewindow=timewindow)
					self.assertEquals(count, points_count_in_timewindow)

					data = self.store.get(
						data_id=data_id, aggregation=aggregation, period=period)
					self.assertEquals(len(data), len(points))
					self.assertEquals(data, sorted_points)

					data = self.store.get(
						data_id=data_id, aggregation=aggregation, period=period,
						timewindow=timewindow)
					self.assertEquals(len(data), points_count_in_timewindow)
					self.assertEquals(data,
						[point for point in sorted_points if point[0] in timewindow])

					self.store.remove(
						data_id=data_id, aggregation=aggregation, period=period,
						timewindow=timewindow)

					# check for count equals 1
					count = self.store.count(data_id=data_id, aggregation=aggregation,
						period=period, timewindow=timewindow)
					self.assertEquals(count, 0)
					count = self.store.count(data_id=data_id, aggregation=aggregation,
						period=period)
					self.assertEquals(count, len(points) - points_count_in_timewindow)

					self.store.remove(data_id=data_id, aggregation=aggregation,
						period=period)
					# check for count equals 0
					count = self.store.count(data_id=data_id, aggregation=aggregation,
						period=period)
					self.assertEquals(count, 0)

if __name__ == '__main__':
	unittest.main()
