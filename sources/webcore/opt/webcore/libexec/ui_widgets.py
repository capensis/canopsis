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
import re

import bottle
from bottle import route, get, put, delete, request, HTTPError, response

## Canopsis
from cstorage import get_storage
from libexec.auth import get_account

logger = clogging.getLogger()

#########################################################################

www_path = os.path.expanduser("~/var/www/")
base_path = "%s/canopsis/" % www_path

def get_internal_widgets():
	wlist = []
	if os.path.exists(base_path + "widgets/"):
		wlist += os.listdir(base_path + "widgets/")
	return wlist

def get_external_widgets():
	wlist = []
	if os.path.exists(www_path + "widgets/"):
		wlist += os.listdir(www_path + "widgets/")
	return wlist

def get_widget_json(json_path):
	try:
		FH = open (json_path, 'r' )
		widget_info = FH.read()
		widget_info = json.loads(widget_info)[0]
		FH.close()
		logger.debug("     + Success")
		return widget_info
	except Exception, err:
		logger.debug("     + Failed (%s)" % err)

	return None

def get_thumb_url(widget_path, url_path):
	default_thumb = "themes/canopsis/resources/images/thumb_widget.png"

	thumb_path = "%s/thumb.png" % widget_path
	if os.path.exists(thumb_path):
		return "%s/thumb.png" % url_path

	return default_thumb

#### GET
@get('/ui/widgets')
def get_all_widgets():
	#account = get_account()
	#storage = get_storage(namespace='object')

	output = []

	logger.debug(" + Search all widgets ...")
	widgets =  get_internal_widgets()
	widgets += get_external_widgets()

	for widget in widgets:
		# Externals
		widget_path = "%s/widgets/%s/" % (www_path, widget)
		url_path = "canopsis/widgets/%s/" % widget
		if not os.path.exists("%s/widget.json" % widget_path):
			# Internals
			widget_path = "%s/widgets/%s/" % (base_path, widget)
			url_path = "widgets/%s/" % widget
			if not os.path.exists("%s/widget.json" % widget_path):
				continue

		logger.info("   + Load '%s' (%s)" % (widget, widget_path))

		widget_info = get_widget_json("%s/widget.json" % widget_path)

		widget_info["thumb"] = get_thumb_url(widget_path, url_path)

		if widget_info:
			output.append(widget_info)
		
	output={'total': len(output), 'success': True, 'data': output}
	return output

#### Widgets CSS
@get('/ui/widgets.css', skip=['checkAuthPlugin'])
def get_widgets_css():
	widgets =  get_internal_widgets()
	widgets += get_external_widgets()

	output = ""

	logger.debug(" + Search all widgets CSS...")
	for widget in widgets:
		css_uri = "/static/canopsis/widgets"

		# Externals
		widget_path = "%s/widgets/%s/" % (www_path, widget)
		if not os.path.exists("%s/widget.json" % widget_path):
			# Internals
			widget_path = "%s/widgets/%s/" % (base_path, widget)
			if not os.path.exists("%s/widget.json" % widget_path):
				return
		else:
			css_uri = "/static/widgets"

		css =		"%s.css" 	% widget
		css_ie =	"%s-ie.css"	% widget
		
		if os.path.exists("%s/%s" % (widget_path, css)):
			logger.debug(" - %s" % css)
			output += "@import '%s/%s/%s';\n" % (css_uri, widget, css)

		#MSIE 8.0 MSIE 7.0
		user_agent = request.environ.get('HTTP_USER_AGENT')
		if re.search("MSIE [78]", user_agent) and os.path.exists("%s/%s" % (widget_path, css_ie)):
			logger.debug(" - %s" % css_ie)
			output += "@import '%s/%s/%s';\n" % (css_uri, widget, css_ie)

	response.content_type = 'text/css'
	return output

#### external widgets libs
@get('/ui/thirdpartylibs.js', skip=['checkAuthPlugin'])
def get_external_widgets_libs():
	widgets =  get_internal_widgets()
	widgets += get_external_widgets()

	output = ""

	logger.debug(" + Search all widgets thirdparty libs...")
	for widget in widgets:
		widget_path = "%s/widgets/%s/libs" % (www_path, widget)
		if not os.path.exists(widget_path):
			continue

		list_of_files = os.listdir(widget_path)
		for filename in [_file for _file in list_of_files if '.js' in _file]:
			with open('%s/%s' % (widget_path,filename), 'r') as f:
				output += f.read()

	response.content_type = 'application/javascript'
	return output