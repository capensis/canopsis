#!/usr/bin/env python
# --------------------------------
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
from celery.task import task
from celery.task.sets import subtask
from cinit import cinit
from caccount import caccount
from crecord import crecord
from cstorage import cstorage
from cfile import cfile
from ctools import cleanTimestamp
from datetime import date
from celerylibs import decorators
from random import randint
import os, sys, json
import time

import hashlib

import task_mail
from wkhtmltopdf.wrapper import Wrapper

from datetime import timedelta, datetime
import calendar

from dateutil.relativedelta import relativedelta

init 	= cinit()
logger 	= init.getLogger('Reporting Task')
logger.setLevel('DEBUG')

@task
@decorators.log_task
def render_pdf(fileName=None, viewName=None, exporting=None, account=None, mail=None, owner=None, orientation='Portrait', pagesize='A4'):

	logger.info('start render')

	logger.debug("fileName: %s " % fileName)
	logger.debug("viewName: %s " % viewName)
	logger.debug("exporting: %s " % exporting)
	logger.debug("account: %s " % account)
	logger.debug("mail: %s " % mail)

	if exporting is None:
		exporting = {"enable": False}

	now = time.time()

	startTime = None
	stopTime = None

	logger.debug("now: %s " % now)

	timezone = exporting.get(
		'timezone',
		{"type": "local", "value": 0})
	if timezone.get('type') == 'utc':
		timezone['value'] = 0
	elif timezone.get('type') == 'server':
		timezone['value'] = time.timezone
	elif 'value' not in timezone:
		timezone['value'] = 0
		logger.error('timezone value does not exist')

	if not exporting.get('enable', False):
		stopTime = now

	elif exporting.get('type', 'duration') == 'duration':

		_datetime = datetime.fromtimestamp(now)

		kwargs = {
			exporting.get('unit', 'days'): int(exporting.get('length', 1))
		}
		rd = relativedelta(**kwargs)
		_datetime -= rd

		startTime = time.mktime(_datetime.timetuple())
		stopTime = now

	elif exporting['type'] == 'fixed':

		logger.debug('now: %s' % now)

		def get_timestamp(_time):
			"""
			Get a timestamp from an input _time struct.
			"""
			result = 0

			_datetime = datetime.utcfromtimestamp(now)

			before = _time.get('before')
			if before is not None:
				for key in before:
					before[key] = int(before[key])
				rd = relativedelta(**before)
				_datetime -= rd

			kwargs = dict()

			day = _time.get('day')
			if day is not None:
				logger.debug("day: %s" % day)
				kwargs['day'] = int(day)

			month = _time.get('month')
			if month is not None:
				logger.debug("month: %s" % month)
				kwargs['month'] = int(month)

			hour = _time.get('hour')
			if hour is not None:
				logger.debug("hour: %s" % hour)
				kwargs['hour'] = int(hour)

			minute = _time.get('minute')
			if minute is not None:
				logger.debug("minute: %s" % minute)
				kwargs['minute'] = int(minute)

			logger.debug("_datetime: %s" % _datetime)

			_datetime = _datetime.replace(**kwargs)

			logger.debug("_datetime: %s" % _datetime)

			day_of_week = _time.get('day_of_week')
			if day_of_week is not None:
				logger.debug("day_of_week: %s" % day_of_week)
				day_of_week = int(day_of_week)
				weekday = calendar.weekday(_datetime.year, _datetime.month, _datetime.day)
				days = weekday - day_of_week
				logger.debug("days %s" % days)
				td = timedelta(days=days)
				_datetime -= td

			result = time.mktime(_datetime.utctimetuple())
			result -= timezone['value']

			logger.debug('result: %s' % result)

			return result

		startTime = None
		_from = exporting.get('from')
		if _from is not None and _from:
			startTime = int(_from.get('timestamp', get_timestamp(_from)))

		logger.info('from : {0} and startTime : {1} ({2})'.format(_from, startTime, datetime.utcfromtimestamp(startTime)))

		stopTime = now
		_to = exporting.get('to')
		if startTime and _to is not None and _to.get('enable', False):
			stopTime = int(_to.get('timestamp', get_timestamp(_to)))

		logger.info('to : {0} and stopTime : {1} ({2})'.format(_to, stopTime, datetime.utcfromtimestamp(stopTime)))

	else:
		logger.error('Wrong exporting type %s' % exporting['type'])

	if viewName is None:
		raise ValueError("task_render_pdf: you must at least provide a viewName")

	#check if the account is just a name or a real caccount
	if isinstance(account ,str) or isinstance(account ,unicode):
		root_account = caccount(user='root',group='root',mail='root@localhost.local')
		root_storage = cstorage(account=root_account, namespace='object')

		bd_account = root_storage.find_one(mfilter={'_id':'account.%s' % str(account)})

		if isinstance(bd_account, crecord):
			account = caccount(bd_account)
			logger.info('Successfuly retrieve right user from db')
		else:
			account = caccount(mail='anonymous@localhost.local')
			logger.info('Anonymous account created')

	#get view options
	storage = cstorage(account=account, namespace='object')
	try:
		view_record = storage.get(viewName,account=account)
	except:
		raise Exception("Impossible to find view '%s' with account '%s'" % (viewName, account._id))

	logger.info("Account '%s' ask a rendering of view '%s' (%s)" % (account.name, view_record.name, viewName,))

	timezone_td = timedelta(
		seconds=int(timezone.get('value', 0)))

	def get_datetime_with_timezone(timestamp):
		result = datetime.utcfromtimestamp(timestamp)
		result -= timezone_td
		return result

	#set fileName
	if fileName is None:
		toDate = get_datetime_with_timezone(int(stopTime))
		if startTime and startTime != -1:
			fromDate = get_datetime_with_timezone(int(startTime))
			fileName = '%s_From_%s_To_%s' % (view_record.name, fromDate, toDate)
		else:
			fileName = '%s_%s' % (view_record.name,toDate)

		fileName += '_Tz_%s_%s' % (timezone.get('type', 'server'), (timezone.get('value', 0) if timezone.get('value', 0) >= 0 \
			else 'minus_%s' % -timezone['value']))
		fileName += '.pdf'

	logger.info('fileName: %s' % fileName)
	ascii_fileName = hashlib.md5(fileName.encode('ascii', 'ignore')).hexdigest()

	#get orientation and pagesize
	view_options = view_record.data.get('view_options', {})
	if isinstance(view_options, dict):
		orientation = view_options.get('orientation', 'Portrait')
		pagesize = view_options.get('pageSize', 'A4')

	logger.info('Orientation: %s' % orientation)
	logger.info('Pagesize: %s' % pagesize)

	wrapper_conf_file = os.path.expanduser("~/etc/wkhtmltopdf_wrapper.json")
	file_path = open(wrapper_conf_file, "r").read()
	file_path = '%s/%s' % (json.loads(file_path)['report_dir'],ascii_fileName)

	#create wrapper object
	wkhtml_wrapper = Wrapper(	ascii_fileName,
							viewName,
							startTime,
							stopTime,
							account,
							wrapper_conf_file,
							orientation=orientation,
							pagesize=pagesize)

	wkhtml_wrapper.logger = logger

	# Run rendering
	logger.info('Run pdf rendering')
	wkhtml_wrapper.run_report()

	logger.info('Put it in grid fs: %s' % file_path)
	try:
		doc_id = put_in_grid_fs(file_path, fileName, account,owner)
	except Exception as e:
		import inspect
		inspect.trace()
		logger.info('Error in generating file (try to not use slink): %s, %s ' % (e, inspect.trace()))
	logger.info('Remove tmp report file with docId: %s' % doc_id)
	os.remove(file_path)

	#Subtask mail (if needed)
	if isinstance(mail, dict):

		#get cfile
		try:
			reportStorage = cstorage(account=account, namespace='files')
			meta = reportStorage.get(doc_id)
			meta.__class__ = cfile
		except Exception, err:
			logger.error('Error while fetching cfile : %s' % err)

		try:
			mail['account'] = account
			mail['attachments'] = meta
			result = task_mail.send.subtask(kwargs=mail).delay()
			result.get()
			result = result.result

			#if subtask fail, raise exception
			if(result['success'] == False):
				raise Exception('Subtask mail have failed : %s' % result['celery_output'][0])

		except Exception, err:
			logger.error(err)
			raise Exception('Impossible to send mail')

	return doc_id

@task
def put_in_grid_fs(file_path, file_name, account,owner=None):
	storage = cstorage(account, namespace='files')
	report = cfile(storage=storage)
	report.put_file(file_path, file_name, content_type='application/pdf')

	if owner:
		report.chown(owner)

	id = storage.put(report)
	if not report.check(storage):
		logger.error('Report not in grid fs')
		return False
	else:
		return id
