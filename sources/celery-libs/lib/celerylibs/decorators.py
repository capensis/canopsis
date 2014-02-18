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

from caccount import caccount
from cstorage import cstorage
from crecord import crecord

import time
from camqp import camqp
import cevent

import clogging

logger 	= clogging.getLogger('Task result to db') 

def simple_decorator(decorator):
    def new_decorator(f):
        g = decorator(f)
        g.__name__ = f.__name__
        g.__module__ = f.__module__
        g.__doc__ = f.__doc__
        g.__dict__.update(f.__dict__)
        return g
    new_decorator.__name__ = decorator.__name__
    new_decorator.__doc__ = decorator.__doc__
    new_decorator.__dict__.update(decorator.__dict__)
    return new_decorator

@simple_decorator
def log_task(func):	
	def wrapper(*args,**kwargs):
		'''
		try:
			task_name = kwargs['_scheduled']
			del kwargs['_scheduled']
		except:
			task_name = None
			logger.info('Not scheduled task')
		'''
		
		try:
			result = func(*args, **kwargs)
			success = True
			logger.info('Task successfully done')
		except Exception, err:
			success = False
			function_error = str(err)
			logger.error(err)
		'''
		try:
			# Get account/storage
			if isinstance(kwargs['account'],unicode):
				account = caccount(user=kwargs['account'])
			else:
				account = kwargs['account']
		except:
			logger.info('No account specified in the task')
			account = caccount()
			
		storage = cstorage(account=account, namespace='task_log')
		taskStorage = cstorage(account=account, namespace='task')
		'''
		
		timestamp = int(time.time())

		# The function have succeed ?
		if success:
			if isinstance(result, list):
				data = result
			else:
				data = [str(result)]

			log = {	'success': True,
					'total': len(data),
					'celery_output': 'Celery task done',
					'timestamp': timestamp,
					'data': data
					}
		else:
			log = {	'success': False,
					'total': 0,
					'celery_output': [ str(function_error) ],
					'timestamp':timestamp,
					'data': []
				  }
		
		'''
		#Put the log
		try:
			# If scheduled
			if task_name:
				logger.info('Task scheduled')
				log_record = crecord(log,name=task_name)
				
				# Replace last log with this one
				try:
					mfilter = {'crecord_name':task_name}
					search = taskStorage.find_one(mfilter)

					if search:
						search.data['log'] = log
						taskStorage.put(search)
						logger.info('Task log updated')
					else:
						logger.error('Task not found in db, can\'t update')
				except Exception, err:
					logger.error('Error when put log in task_log %s' % err)
				
			else:
				logger.info('Not a scheduled task, put log in db')
				log_record = crecord(log)
				
			# Put log in storage
			storage.put(log_record)
		except Exception, err:
			logger.error('Error when put log in task_log %s' % err)

		
		# Publish Amqp event
		if success:
			status=0
		else:
			status=1

		event = cevent.forger(
			connector='celery',
			connector_name='task_log',
			event_type='log',
			source_type='resource',
			output=log['output'],
			state=status
			)
		logger.debug('Send Event: %s' % event)
		key = cevent.get_routingkey(event)
		
		amqp = camqp()
		amqp.start()
		
		amqp.publish(event, key, amqp.exchange_name_events)
		
		amqp.stop()
		amqp.join()
		'''
		return log
	return wrapper
