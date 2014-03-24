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

from cengine import cengine

from cstorage import get_storage
from caccount import caccount

from dateutil.rrule import rrulestr
import datetime
import time


NAME='scheduler'


class engine(cengine):
	def __init__(self, *args, **kwargs):
		super(engine, self).__init__(name=NAME, *args, **kwargs)

		account = caccount(user='root', group='root')
		self.storage = get_storage('jobs', account=account).get_backend()

	def pre_run(self):
		self.beat()

	def beat(self):
		self.logger.info('Reload jobs')

		now = int(time.time())
		prev = now - self.beat_interval

		jobs = self.storage.find({
			'last_execution': {'$lte': prev}
		})

		for job in jobs:
			self.logger.info('Job: {0}'.format(job))

			if job['last_execution'] < 0:
				self.do_job(job)

			else:
				jobStart = datetime.datetime.fromtimestamp(job['start'])

				dtstart = datetime.datetime.fromtimestamp(prev)
				dtend = datetime.datetime.fromtimestamp(now)

				occurences = list(rrulestr(job['rrule'], dtstart=jobStart).between(dtstart, dtend))

				if len(occurences) > 0:
					self.do_job(job)

	def do_job(self, job):
		self.logger.info('Execute job: {0}'.format(job))

		job['params']['id'] = job['_id']

		self.amqp.publish(
			job['params'],
			'task_{0}'.format(job['type']),
			'amq.direct'
		)

		self.storage.update({
			'_id': job['_id']
		},{
			'$set': {
				'last_execution': int(time.time())
			}
		})