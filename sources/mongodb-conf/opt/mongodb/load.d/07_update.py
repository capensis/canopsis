#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

import logging, math
from cstorage import get_storage
from caccount import caccount
from datetime import timedelta
import time

logger = None
root = caccount(user="root", group="root")
storage = get_storage(account=root, namespace='object')

def init():
	logger.info(' + init')

def update():
	init()
	update_schedule()

def update_schedule():

	records = storage.find(
		{'crecord_type': 'schedule'},
		namespace='object',
		account=root)

	for record in records:
		kwargs = record.data['kwargs']
		cron = record.data['cron']

		#--mail
		if 'mail' in kwargs:
			mail = kwargs['mail']
			if mail:
				if 'sendMail' in mail and mail['sendMail']:
					record.data['exporting_mail'] = True
				if 'recipients' in mail:
					record.data['exporting_recipients'] = mail['recipients']
				if 'subject' in mail:
					record.data['exporting_subject'] = mail['subject']

		#--kwargs
		if 'account' in kwargs:
			record.data['exporting_account'] = kwargs['account']

		if 'task' in kwargs:
			record.data['exporting_task'] = kwargs['task']

		if 'method' in kwargs:
			record.data['exporting_method'] = kwargs['method']

		if 'owner' in kwargs:
			record.data['exporting_owner'] = kwargs['owner']

		if 'viewname' in kwargs:
			kwargs['viewName'] = kwargs['viewname']
			record.data['exporting_viewName'] = kwargs['viewname']
			del kwargs['viewname']

		if 'starttime' in kwargs:
			del kwargs['starttime']

		if 'day_of_week' in cron:
			record.data['crontab_day_of_week'] = cron['day_of_week']

		if 'day' in cron:
			record.data['crontab_day'] = cron['day']

		if 'mouth' in cron:
			record.data['crontab_month'] = cron['month']

		if "every" in record.data:
			record.data['frequency'] = record.data['every']
			del record.data['every']

		exporting = {
			"type": "duration",
			"unit": "day",
			"length": 1
		}

		exporting.update(
			record.data.get('exporting', dict()))

		if 'interval' in kwargs and kwargs['interval'] is not None:
			nbDays = timedelta(seconds=int(kwargs['interval'])).days

			del kwargs['interval']

			if nbDays >= 365:
				exporting['unit'] = 'years'
				exporting['length'] = int(nbDays/365)
			elif nbDays >= 30:
				exporting['unit'] = 'months'
				exporting['length'] = int(nbDays/30)
			elif nbDays >= 7:
				exporting['unit'] = 'weeks'
				exporting['length'] = int(nbDays/7)
			elif nbDays >= 1:
				exporting['unit'] = 'days'
				exporting['length'] = math.floor(nbDays)
			else:
				exporting['unit'] = 'hours'
				exporting['length'] = math.floor(nbDays)

		if 'interval' in kwargs :
			del kwargs['interval']

		if "exporting_intervalUnit" in kwargs and "exporting_intervalLength" in kwargs:
			exporting.update({
				"type": kwargs.get("exporting_advanced", "duration"),
				"length": kwargs['exporting_intervalLength'],
				"unit": kwargs['exporting_intervalUnit']
			})

			del kwargs['exporting_intervalLength']
			del kwargs['exporting_intervalUnit']
			if 'exporting_advanced' in kwargs:
				del kwargs['exporting_advanced']

		kwargs['exporting'] = exporting
		record.data['exporting'] = exporting

		record.data['kwargs'] = kwargs

		storage.put(record)
