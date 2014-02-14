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

import sys, os, clogging, json, subprocess
import gevent

import bottle
from bottle import route, get, delete, request, HTTPError, post, static_file, response

from urllib import quote
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
import task_reporting

#import protection function
from libexec.auth import get_account
from libexec.account import check_group_rights

logger = clogging.getLogger()

#group who have right to access 
group_managing_access = ['group.CPS_reporting_admin']

#########################################################################

@post('/reporting/:startTime/:stopTime/:view_name/:mail',checkAuthPlugin={'authorized_grp':group_managing_access})
@post('/reporting/:startTime/:stopTime/:view_name',checkAuthPlugin={'authorized_grp':group_managing_access})
def generate_report(startTime, stopTime,view_name,mail=None):
	stopTime = int(stopTime)
	startTime = int(startTime)

	account = get_account()
	storage = cstorage(account=account, namespace='object')

	if mail:
		try:
			mail = json.loads(mail)
		except Exception, err:
			logger.error('Error while transform string mail to object' % err)
			mail=None
	try:
		record = storage.get(view_name,account=account)
	except Exception, err:
		logger.error(err)
		return {'total': 1, 'success': False, 'data': [str(err)] }
		

	toDate = str(date.fromtimestamp(int(stopTime)))
	if startTime and startTime != -1:
		fromDate = str(date.fromtimestamp(int(startTime)))
		file_name = '%s_From_%s_To_%s.pdf' % (record.name,fromDate,toDate)
	else:
		file_name = '%s_%s.pdf' % (record.name,toDate)

	logger.debug('file_name:   %s' % file_name)
	logger.debug('view_name:   %s' % view_name)
	logger.debug('startTime:   %s' % startTime)
	logger.debug('stopTime:    %s' % stopTime)
	logger.debug('mail:    %s' % mail)

	result = None
	
	try:
		logger.debug('Run celery task')
		result = task_reporting.render_pdf.delay(
										fileName=file_name,
										viewName=view_name,
										startTime=startTime,
										stopTime=stopTime,
										account=account,
										mail=mail
										)
		result.wait()
		result = result.result

	except Exception, err:
		return {'total': 1, 'success': False, 'data': [str(err)] }

	if not result or len(result['data']) == 0:
		logger.error('Error while generating pdf : %s' % result['celery_output'])
		return {'total': 0, 'success': False, 'data': [result['celery_output']] }

	_id = str(result['data'][0])
	
	logger.debug(' + File Id: %s' % _id)
	return {'total': 1, 'success': True, 'data': [{'id': _id}] }
	
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
	
