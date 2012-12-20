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
from datetime import date
from celerylibs import decorators
from random import randint
import os, sys, json 
import time

import hashlib

#TEST
import task_node
import task_mail

init 	= cinit()
logger 	= init.getLogger('Reporting Task') 

@task
@decorators.log_task
def render_pdf(filename=None, viewname=None, starttime=None, stoptime=None, account=None, wrapper_conf_file=None, mail=None, owner=None, orientation='Portrait', pagesize='A4'):
	
	#prepare stoptime and starttime
	if stoptime:
		stoptime = int(stoptime)*1000
	else:
		stoptime = time.time()*1000
	if starttime:
		starttime = int(starttime)*1000
	else:
		startime = 0

	if viewname is None:
		raise ValueError("task_render_pdf : you must at least provide a viewname")

	################setting variable for celery auto launch task##################"

	#if no wrapper conf file set, take the default
	if wrapper_conf_file is None:
		wrapper_conf_file = os.path.expanduser("~/etc/wkhtmltopdf_wrapper.json")
	
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
		
	#set start time if needed
	if starttime != 0:
		starttime = stoptime - starttime 
	
	#get view options
	storage = cstorage(account=account, namespace='object')
	try:
		view_record = storage.get(viewname,account=account)
	except:
		raise Exception("Impossible to find view '%s' with account '%s'" % (viewname, account._id))

	#set filename
	if filename is None:
		toDate = date.fromtimestamp(int(stoptime) / 1000)
		
		if starttime !=0:
			fromDate = date.fromtimestamp(int(starttime) / 1000)
			filename = '%s_From_%s_To_%s.pdf' % (view_record.name, fromDate, toDate) 
		else:
			filename = '%s_%s.pdf' % (view_record.name, toDate) 
		
	ascii_filename = hashlib.md5(filename.encode('ascii', 'ignore')).hexdigest()
	
	logger.info('Filename: %s' % filename)

	#get orientation and pagesize
	view_options = view_record.data.get('view_options', {})
	if isinstance(view_options, dict):
		orientation = view_options.get('orientation', 'Portrait')
		pagesize = view_options.get('pageSize', 'A4')

	logger.info('Orientation: %s' % orientation)
	logger.info('Pagesize: %s' % pagesize)

	##############################################################################
	
	libwkhtml_dir=os.path.expanduser("~/lib")
	sys.path.append(libwkhtml_dir)

	logger.info('Load wkhtmltopdf.wrapper')
	import wkhtmltopdf.wrapper

	# Generate config
	settings = wkhtmltopdf.wrapper.load_conf(	ascii_filename,
							viewname,
							starttime,
							stoptime,
							account,
							wrapper_conf_file,
							orientation=orientation,
							pagesize=pagesize)

	file_path = open(wrapper_conf_file, "r").read()
	file_path = '%s/%s' % (json.loads(file_path)['report_dir'],ascii_filename)

	# Run rendering
	logger.debug('Run pdf rendering')
	result = wkhtmltopdf.wrapper.run(settings)
	result.wait()

	logger.debug('Put it in grid fs')
	id = put_in_grid_fs(file_path, filename, account,owner)
	logger.debug('Remove tmp report file')
	os.remove(file_path)
	
	#Subtask mail (if needed)
	if isinstance(mail, dict):
		#get cfile
		try:
			reportStorage = cstorage(account=account, namespace='files')
			meta = reportStorage.get(id)
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
		
	return id

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
