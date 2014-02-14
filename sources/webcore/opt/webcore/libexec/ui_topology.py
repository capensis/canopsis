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
import sys, os, clogging, json

import bottle
from bottle import route, get, put, delete, request, HTTPError, post, response

#import protection function
from libexec.auth import get_account

#group_managing_access = ['']
#########################################################################

logger = clogging.getLogger()

import sys, os
operators_path=os.path.expanduser('~/opt/amqp2engines/engines/topology')
sys.path.append(operators_path)

@get('/topology/getOperators')
def get_operators():
	operators = []

	for opfile in os.listdir(operators_path):
		try:
			operator = opfile.split('.')
			if operator[1] == 'py':
				module = __import__(operator[0])
				operators.append(module.options)
				del sys.modules[operator[0]]
		except Exception, err:
			logger.warning("Impossible to parse '%s' (%s)" % (opfile, err))
        
	return { 'total':len(operators), 'success':True ,'data': operators }
