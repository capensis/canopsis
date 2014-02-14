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

import clogging

logger = clogging.getLogger('Reporting Task')

@task
@decorators.log_task
def render_pdf(fileName=None, viewName=None, startTime=None, stopTime=None, interval=None, account=None, mail=None, owner=None, orientation='Portrait', pagesize='A4'):
	if not stopTime:
		stopTime = int(time.time())

	if not startTime:
		if interval:
			startTime = stopTime - interval
		else:
			startTime = -1
	
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

	#set fileName
	if fileName is None:
		toDate = date.fromtimestamp(int(stopTime))
		if startTime and startTime != -1:
			fromDate = date.fromtimestamp(int(startTime))
			fileName = '%s_From_%s_To_%s.pdf' % (view_record.name, fromDate, toDate) 
		else:
			fileName = '%s_%s.pdf' % (view_record.name,toDate) 

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

	# Run rendering
	logger.debug('Run pdf rendering')
	wkhtml_wrapper.run_report()

	logger.info('Put it in grid fs: %s' % file_path)
	doc_id = put_in_grid_fs(file_path, fileName, account,owner)
	logger.debug('Remove tmp report file with docId: %s' % doc_id)
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
