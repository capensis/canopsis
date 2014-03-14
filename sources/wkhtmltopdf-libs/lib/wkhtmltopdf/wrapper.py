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
import os, logging, signal, time, json,random
from tempfile import mkstemp
from subprocess import Popen
from time import sleep
from time import time

logging.basicConfig()

class Wrapper(object):
	def __init__(self,filename, viewName, startTime, stopTime, subset_selection, account, wrapper_conf_file, orientation='Portrait', pagesize='A4'):
		self.logger = logging.getLogger('[Wkhtml wrapper]')
		self.logger.setLevel(logging.DEBUG)

		conf = open(wrapper_conf_file, "r").read()
		self.settings = json.loads(conf)

		self.settings['filename'] = filename
		self.settings['viewName'] = viewName
		self.settings['startTime'] = startTime
		self.settings['stopTime'] = stopTime
		self.settings['subset_selection'] = subset_selection
		self.settings['account'] = account
		self.settings['orientation'] = orientation
		self.settings['pagesize'] = pagesize

		self.xauth = mkstemp()[1]
		self.xDisplay = random.randint(1, 500)
		self.currentX = None

	def run_report(self):
		self.create_xvfb()
		self.export_env(self.xDisplay)
		self.create_report_dir(self.settings['report_dir'])
		self.launch_wkhtml()
		self.clean_x(self.xauth)

	def create_xvfb(self):
		self.logger.debug('Get free Xlock')
		while os.path.isfile("/tmp/.X%s-lock" % self.xDisplay):
			self.xDisplay = random.randint(1, 500)

		cmd = 'Xvfb -screen 0 1024x768x24 -terminate -auth %s -nolisten tcp :%s >/dev/null 2>&1 &' % (self.xauth, self.xDisplay)
		self.logger.debug('Launched cmd for creating Xvfb display')
		self.logger.debug(cmd)
		self.currentX = Popen(cmd, shell=True)

	def export_env(self,environnement):
		self.logger.debug('Set DISPLAY env variable to %s' % environnement)
		os.environ['DISPLAY'] = ':%s' % environnement

	def create_report_dir(self,directory):
		self.logger.debug("Check if report directory exist")

		if not os.path.isdir(directory):
			self.logger.debug("Create directory as %s" % directory)
			os.makedirs(directory)

	def clean_x(self,xauth):
		os.remove(xauth)
		if self.currentX:
			self.currentX.terminate()

	def launch_wkhtml(self):
		output = Popen("wkhtmltopdf -h >> /dev/null",shell=True)
		output.wait()

		# python none to js null
		export_from = self.settings.get('startTime')
		if not export_from:
			export_from = -1

		opts = ' '.join(self.settings['opts'])

		cmd = "wkhtmltopdf -O %s -s %s %s %s %s --window-status %s\
				-T 21mm --header-line --header-spacing 5\
				'http://127.0.0.1:8082/%s/static/canopsis/index.html?exportMode=true&view_id=%s&from=%s&to=%s&subset_selection=%s&authkey=%s'\
				'%s/%s' 2>&1 | grep -v\
				'settings.windowStatus:ready'" % (
													self.settings['orientation'], 
													self.settings['pagesize'], 
													opts, 
													self.settings['header'], 
													self.settings['footer'], 
													self.settings['windowstatus'], 
													self.settings['account'].data.get('locale', 'en'),
													self.settings['viewName'],
													export_from,
													self.settings.get('stopTime', int(time())),
													self.settings.get('subset_selection', "{}"),
													self.settings['account'].get_authkey(),
													self.settings['report_dir'], 
													self.settings['filename']
													)

		self.logger.debug('wkhtmltopdf will be launched with the following command:')
		self.logger.debug(cmd)
		result = Popen(cmd, shell=True)
		waitTime = 0

		#return 0 if everything's ok
		while result.poll() == None:
			sleep(2)
			waitTime = waitTime + 2
			if waitTime >= self.settings.get('timeout',300):
				result.kill()
		self.logger.debug('Pdf rendered, leaving wrapper')