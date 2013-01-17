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

import sys, os, logging, json
import re

import bottle
from bottle import route, get, put, delete, request, HTTPError, response

## Canopsis
from cstorage import get_storage
from libexec.auth import check_auth, get_account

logger = logging.getLogger("ui-widgets")

#########################################################################

base_path = os.path.expanduser("~/var/www/canopsis/")

def get_widgets():
	return os.listdir(base_path + "widgets/")

#### GET
@get('/ui/widgets')
def get_all_widgets():
	#account = get_account()
	#storage = get_storage(namespace='object')

	output = []

	widgets = get_widgets()

	logger.debug(" + Search all widgets ...")
	for widget in widgets:
		widget_path = "%s/widgets/%s/" % (base_path, widget)

		logger.debug("   + Load '%s' (%s)" % (widget, widget_path))
		try:
			FH = open (widget_path + "/widget.json", 'r' )
			widget_info = FH.read()
			widget_info = json.loads(widget_info)
			
			output.append(widget_info[0])

			FH.close()
			logger.debug("     + Success")
		except Exception, err:
			logger.debug("     + Failed (%s)" % err)
		
	output={'total': len(output), 'success': True, 'data': output}
	return output

#### Widgets CSS
@get('/ui/widgets.css', skip=['checkAuthPlugin'])
def get_widgets_css():
	widgets = get_widgets()
	output = ""

	logger.debug(" + Search all widgets CSS...")
	for widget in widgets:
		css = "widgets/%s/%s.css" % (widget, widget)
		iecss = "widgets/%s/%s-ie.css" % (widget, widget)
		
		if os.path.exists(base_path + css):
			logger.debug(" - %s" % css)
			output += "@import '/static/canopsis/%s';\n" % css

		#MSIE 8.0 MSIE 7.0
		user_agent = request.environ.get('HTTP_USER_AGENT')
		if re.search("MSIE [78]", user_agent) and os.path.exists(base_path + iecss):
			logger.debug(" - %s" % iecss)
			output += "@import '/static/canopsis/%s';\n" % iecss

	response.content_type = 'text/css'
	return output

