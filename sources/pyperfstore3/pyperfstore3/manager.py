#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------


import time
from md5 import new as md5

from operator import itemgetter

from datetime import datetime

from .timewindow import TimeWindow, Period


class Manager(object):

	__slots__ = ['perfdata3', 'logger']

	def __init__(self, perfdata3, logger):

		super(Manager, self).__init__()

		self.perfdata3 = perfdata3
		self.logger = logger

	@staticmethod
	def get_metric_id(component, resource, metric):

		metric_id_md5 = md5()
		# add co in id
		metric_id_md5.update(component.encode('ascii', 'ignore'))

		if resource is not None:
			# add re in id
			metric_id_md5.update(resource.encode('ascii', 'ignore'))
		# add me in id
		metric_id_md5.update(metric.encode('ascii', 'ignore'))

		result = metric_id_md5.hexdigest()

		return result

	@staticmethod
	def get_document_id(metric_id, timestamp, period):

		md5_result = md5()

		# add id_timestamp in id
		md5_result.update(str(timestamp).encode('ascii', 'ignore'))
		# add period in id
		md5_result.update(period.unit.encode('ascii', 'ignore'))

		result = metric_id.join(md5_result.hexdigest())

		return result

	def put_data(self, metric_id, value, timestamp=time.time(), meta=dict(), period=Period(unit=Period.MINUTE, value=60)):

		sliding_period = period.next_period()

		timestamp = int(timestamp)

		id_timestamp = sliding_period.sliding_timestamp(timestamp)

		self.logger.debug(' + id_timestamp: {0}'.format(id_timestamp))

		_id = self.get_document_id(metric_id, id_timestamp, period)

		field_name = "values.{0}".format(timestamp - id_timestamp)

		result = self.perfdata3.update(
			{
				'_id': _id
			},
			{
				'$set': {
					'period': period.to_dict(),
					'metric_id': metric_id,
					'timestamp': id_timestamp,
					'last_update': timestamp,
					field_name: value,
					'meta': meta
				}
			},
			upsert=1,
			w=1)

		error = result.get("writeConcernError", None)

		if error is not None:
			self.logger.error(' error in updating document: {0}'.format(error))

		error = result.get("writeError")

		if error is not None:
			self.logger.error(' error in updating document: {0}'.format(error))

		self.logger.debug(' + metric updated: {0}'.format(meta))

	def get_data(self, metric_id, interval, period=Period(unit=Period.MINUTE, value=60)):

		timewindow = TimeWindow(**interval)

		sliding_period = period.next_period()

		start = sliding_period.sliding_timestamp(interval.start)

		stop = sliding_period.sliding_timestamp(interval.stop)

		dt = datetime.utcfromtimestamp(start)

		dtstop = datetime.utcfromtimestamp(stop)

		ids = list()

		start_id = Manager.get_document_id(metric_id, start, period)
		ids.append(start_id)

		stop_id = Manager.get_document_id(metric_id, stop, period)
		ids.append(stop_id)

		while dt <= dtstop:

			dt = timewindow.get_next_date(dt, period)

			t = time.mktime(dt)

			t_id = Manager.get_document_id(metric_id, t, period)

			ids.append(t_id)

		request = {"_id": {'$in': ids}}

		documents = self.perfdata3.find(
			request,
			projection={'timestamp': 1, 'values': 1, 'meta': 1})

		documents.hint([('_id', 1)])

		meta = None

		result = list()

		for document in documents:
			if meta is None:
				meta = document.get('meta')

			timestamp = int(document.get('timestamp'))
			values = document.get('values')

			for t, value in values.iteritems():
				result.append( (timestamp + int(t), value) )

		result.sort(key=itemgetter(0))

		return result
