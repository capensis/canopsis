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

from caccount import caccount
from cstorage import cstorage
from crecord import crecord

import time
from camqp import camqp
import cevent

init 	= cinit()
logger = init.getLogger()

def launch_celery_task(*args,**kwargs):
	if kwargs.has_key('task') and kwargs.has_key('method'):
		try:
			amqp = camqp()
			amqp.start()

			timer_begin = int(time.time())
			
			#----------Get task informations
			task_name = kwargs['_scheduled']
			celery_task_name = kwargs['task']
			
			module = __import__(kwargs['task'])
			exec "task = module.%s" % kwargs['method']
			
			#-------------Clear arguments
			methodargs = kwargs
			del methodargs['task']
			del methodargs['method']
			del kwargs['_scheduled']
			
			#-------------execute task
			success = True
			
			try:
				result = task.delay(*args,**methodargs)
				result.get()
				result = result.result
				
				if not result['success']:
					raise Exception('Celery task failed')
				
			except Exception, err:
				success = False
				aps_error = str(err)
				logger.error(err)

			#------------Get account and storage
			try:
				if isinstance(kwargs['account'],unicode) or isinstance(kwargs['account'],str):
					account = caccount(user=kwargs['account'])
				else:
					account = kwargs['account']
				#logger.error(account)
				logger.info('Caccount create from passed arguments')
			except Exception, err:
				logger.info('No account specified in the task')
				account = caccount()
			
			try:
				storage = cstorage(account=account, namespace='object')
			except Exception, err:
				logger.info('Error while fecthing storages : %s' % err)
				success = False
				aps_error = str(err)
			
			#-------------time operation-------------------
			timestamp = int(time.time())
			execution_time = (timestamp - timer_begin)
			#-------------Check if function have succeed
			'''
			if success:
				if isinstance(result, list):
					data = result
				else:
					data = [str(result)]

				log = {	'success': True,
						'total': len(data),
						'output':'Task done',
						'timestamp': timestamp,
						'data': data,
						}
				logger.info('Task was a success')
			else:
				log = {	'success': False,
						'total': 0,
						'output': [ str(function_error) ],
						'timestamp':timestamp,
						'data': [],
					  }
				logger.info('Task have failed')
			'''
				
			#-----------------Put log in schedule attribut----------------
			try:
				mfilter = {'crecord_name':task_name}
				search = storage.find_one(mfilter)

				#add execution time
				result['duration'] = execution_time

				if search:
					search.data['log'] = result
					storage.put(search)
					logger.info('Task log updated')
				else:
					logger.error('Task not found in db, can\'t update')
			except Exception, err:
				logger.error('Error when put log in task_log %s' % err)
				success = False
				aps_error = str(err)
			
			#-------------------------Put log in db-------------------------
			'''
			try:
				log['task_name'] = task_name
				log_record = crecord(result,name=task_name)
				storage.put(log_record)
				logger.info('log put in db')
			except Exception,err:
				logger.error('log not added to db, reason : %s' % err)
			'''
			#---------------------Publish amqp event-------------
			# Publish Amqp event
			if success == True:
				status=0
				#result['aps_output'] = 'APS task success'
				#task_output = result
				task_output = ('APS : Task success - Celery : %s - Duration : %is' % (result['celery_output'],execution_time))
			else:
				status=1
				#result['aps_output'] = aps_error
				#task_output = result
				task_output = ('APS : %s - Celery : %s - Duration : %is' % (aps_error,result['celery_output'],execution_time))
				
			event = cevent.forger(
				connector='celery',
				connector_name='task_log',
				event_type='log',
				source_type='resource',
				resource=('task.%s.%s.%s' %  (celery_task_name,account.user,task_name)), 
				output=task_output,
				state=status
				)
			
			#logger.info('Send Event: %s' % event)
			key = cevent.get_routingkey(event)

			amqp.publish(event, key, amqp.exchange_name_events)
						
			logger.info('Amqp event published')

			# Stop AMQP
			amqp.stop()
			amqp.join()

			#--------------------return result-------------------
			return result
			
		except Exception, err:
			logger.error('Error in aps running : %s' % err)
	else:
		logger.error('No task given')
		
	
	
