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
# ---------------------------------

import urllib, urllib2, cookielib, json
import clogging


class cwebservices(object):
	def __init__(self, host="127.0.0.1", port=8082):

		self.logger = clogging.getLogger()

		self.logger.debug('Init urlib object ...')
		self.base_url = 'http://' + host + ':' + str(port)

		self.jar = cookielib.CookieJar()
		self.jar.clear_session_cookies()

		self.handler = urllib2.HTTPCookieProcessor(self.jar)
		self.opener = urllib2.build_opener(self.handler)

		urllib2.install_opener(self.opener)

		self.is_login = False

	def get(self, uri, parsing=True):
		url = self.base_url + uri
		self.logger.debug(' + GET '+url)
		data = urllib2.urlopen(url).read()

		if parsing:
			self.logger.debug('   + Try to parse json Data ...')
			try:
				data_json = json.loads(str(data))
				data = data_json['data']
				state = data_json['success']
			except:
				self.logger.debug('     + Failed')
				raise Exception("Failed to parse response ...")

			if not state:
				raise Exception("Request marked failed by server ...")
			
			self.logger.debug('     + Success')
	
		return data

	def login(self, login, password):
		self.logger.debug("Login with '%s'" % login)
		self.get("/auth/" + login + "/" + password, False)

		#print self.get("/online")

		self.is_login = True

	def logout(self):
		self.logger.debug("Logout.")
		self.get("/logout", False)
		
	def put_view(self,view_id=None,view_name=None,parent_id=None,user='root',leaf=True):
		if view_name and parent_id:
			if not view_id:
				view_id = 'view.%s.%s' % (user,view_name)
			url = '%s/ui/view' % self.base_url
			data = {
				'id': view_id,
				'_id': view_id,
				'crecord_name':view_name,
				'parentId':parent_id,
				'items':[],
				'leaf':leaf
				}
			
			req = urllib2.Request(url, json.dumps(data))
			response = self.opener.open(req)
			return response.read()
		else:
			self.logger.debug("Missing parameters (must specified view name/parent)" % login)
			return None
			
	def rename_view(self,view_id,new_name):
		url = '%s/ui/view' % self.base_url
		data = json.dumps({
				'crecord_name':new_name,
				'_id':view_id
			})
		req = urllib2.Request(url, data)
		req.get_method = lambda: 'PUT'
		response = self.opener.open(req)
		return response.read()
	
	def change_view_parent(self,view_id,parent_id):
		url = '%s/ui/view' % self.base_url
		data = json.dumps({
				'_id':view_id,
				'parentId':parent_id
			})
		req = urllib2.Request(url, data)
		req.get_method = lambda: 'PUT'
		response = self.opener.open(req)
		return response.read()
		
	def delete_view_or_dir(self,ids):
		url = '%s/ui/view' % self.base_url
		if not isinstance(ids,list):
			ids = [ids]
		
		data = []
		
		for _id in ids:
			data.append({'_id':_id})
	
		req = urllib2.Request(url, json.dumps(data))
		req.get_method = lambda: 'DELETE'
		response = self.opener.open(req)
		return response.read()
		
	def valid_server_response(self,response):
		response = response['data']
		for record in response:
			if not response[record].get('success',False):
				raise Exception(response[record].get('output','error'))
		return True

	def __del__(self):
		if self.is_login:
			self.logout()

