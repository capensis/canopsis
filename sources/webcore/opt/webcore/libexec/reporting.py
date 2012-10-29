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

import sys, os, logging, json, subprocess
#from gevent import monkey; monkey.patch_all()
import gevent

import bottle
from bottle import route, get, delete, request, HTTPError, post, static_file, response

from urllib import quote
#gridfs
from pymongo import Connection
import gridfs

import time
from datetime import date

## Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord
from cfile import cfile

import task_mail

#import protection function
from libexec.auth import check_auth, get_account
from libexec.account import check_group_rights

logger = logging.getLogger('Reporting')

#group who have right to access 
group_managing_access = ['group.CPS_reporting_admin']

#########################################################################

@post('/reporting/:startTime/:stopTime/:view_name/:mail',checkAuthPlugin={'authorized_grp':group_managing_access})
@post('/reporting/:startTime/:stopTime/:view_name',checkAuthPlugin={'authorized_grp':group_managing_access})
def generate_report(startTime, stopTime,view_name,mail=None):
	account = get_account()
	storage = cstorage(account=account, namespace='object')
	
	if(isinstance(mail,str)):
		try:
			mail = json.loads(mail)
		except Exception, err:
			logger.error('Error while transform string mail to object' % err)
	
	#get crecord name of the view (id is really harsh)
	try:
		record = storage.get(view_name,account=account)
		
		fromDate = str(date.fromtimestamp(int(startTime) / 1000))
		toDate = str(date.fromtimestamp(int(stopTime) / 1000))
		
		file_name = '%s_From_%s_To_%s.pdf' % (record.name,fromDate,toDate)
	except  Exception, err:
		logger.error(err)
		name_array = view_name.split('.')
		file_name = name_array[len(name_array)-1]
		file_name += '_' + str(date.fromtimestamp(int(startTime) / 1000)) +'.pdf'

	fileName = None
	
	logger.debug('file_name:   %s' % file_name)
	logger.debug('view_name:   %s' % view_name)
	logger.debug('startTime:   %s' % startTime)
	logger.debug('stopTime:    %s' % stopTime)
	
	try:
		import task_reporting
	except Exception, err:
		logger.debug("Check your celeryconfig.py, if you have reporting task imported")
		logger.debug(err)
	try:
		logger.debug('Run celery task')
		result = task_reporting.render_pdf.delay(file_name,
										view_name,
										startTime,
										stopTime,
										account,
										os.path.expanduser("~/etc/wkhtmltopdf_wrapper.json"),
										mail)
		result.wait()
		fileName = result.result
	except Exception, err:
		logger.error(err)

	if fileName:
		if len(fileName['data']) > 0:
			return {'total': 1, 'success': True, 'data': { 'id': str(fileName['data'][0])}}
		else:
			logger.error('Error while generating pdf : %s' % fileName['celery_output'])
			return {'total': 0, 'success': False, 'data': {}}
	else:
		logger.debug('file not found, error while generating pdf')
		return {'total': 0, 'success': False, 'data': {}}
	
@post('/sendreport')
def send_report():
	account = get_account()
	reportStorage = cstorage(account=account, namespace='files')

	recipients = request.params.get('recipients', default=None)
	_id = request.params.get('_id', default=None)
	body = request.params.get('body', default=None)
	subject = request.params.get('subject', default=None)
	
	meta = reportStorage.get(_id)
	meta.__class__ = cfile
	
	mail = {
		'account':account,
		'attachments': meta,
		'recipients':recipients,
		'subject':subject,
		'body': body,
	}
	
	try:
		task = task_mail.send.delay(**mail)
		output = task.get()
		return {'success':True,'total':'1','data':{'output':output}}
	except Exception, err:
		logger.error('Error when run subtask mail : %s' % err)
		return {'success':False,'total':'1','data':{'output':'Mail sending failed'}}


# For highcharts
@post('/export_svg')
def export_svg():
	filename = request.params.get('filename', default=None)
	svg = request.params.get('svg', default=None)
	
	if not filename:
		filename = "chart.svg"
	else:
		filename += ".svg"
		
	
	logger.debug("Export SVG image: %s" % filename)
	
	if svg and filename:
		response.set_header('Content-Disposition', 'attachment; filename="%s"' % filename)
		response.content_type = 'image/svg+xml'
		return svg
	
