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
import os, logging, signal, time, json
from subprocess import Popen

logger = logging.getLogger('WRAPPER')

def load_conf(filename, viewname, starttime, stoptime, account, wrapper_conf_file):
	conf = open(wrapper_conf_file, "r").read()
	settings = json.loads(conf)
	settings['filename'] = filename
	settings['viewname'] = viewname
	settings['starttime'] = starttime
	settings['stoptime'] = stoptime
	settings['account'] = account
	return settings

def check_xorg(lock, xvfb_cmd):
	import random
	from tempfile import mkstemp
	global DISP, XAUTH

	XAUTH = mkstemp()[1]
	DISP = random.randint(1, 500)

	while os.path.isfile("/tmp/.X%s-lock" % DISP):
		DISP = random.randint(1, 500)

	cmd = 'Xvfb -screen 0 1024x768x24 -terminate -auth %s -nolisten tcp :%s >/dev/null 2>&1 &' % (XAUTH, DISP)
	logger.debug(cmd)
	output = Popen(cmd, shell=True)

def export_env(interface):
	logger.debug(" [WK_WRAPPER] :: Set env DISPLAY to %s" % DISP)
	os.environ['DISPLAY'] = ':%s' % DISP

def check_report_dir(report_dir):
	logger.debug(" [WK_WRAPPER] :: Check if report directorie exist")
	if not os.path.isdir(report_dir):
		logger.debug(" [WK_WRAPPER] :: Create it at %s" % report_dir)
		os.makedirs(report_dir)

def	get_cookie(cookiejar, account):
	output = Popen("wkhtmltopdf -h >> /dev/null",shell=True)
	output.wait()
	authkey = account.get_authkey()
	logger.debug(" [WK_WRAPPER] :: Recreate cookie (%s)" % cookiejar)
	#logger.debug("wkhtmltopdf --load-error-handling ignore --cookie-jar %s \"http://127.0.0.1:8082/auth/%s/%s?cryptedKey=True\" /dev/null" % (cookiejar, account.user, authkey))
	#output = Popen("wkhtmltopdf --load-error-handling ignore --cookie-jar %s \"http://127.0.0.1:8082/auth/%s/%s?cryptedKey=True\" /dev/null" % (cookiejar, account.user, authkey), shell=True)

	logger.error("wkhtmltopdf --load-error-handling ignore --cookie-jar %s \"http://127.0.0.1:8082/keyAuth/%s/%s\" /dev/null" % (cookiejar, account.user, authkey))
	output = Popen("wkhtmltopdf --load-error-handling ignore --cookie-jar %s \"http://127.0.0.1:8082/keyAuth/%s/%s\" /dev/null" % (cookiejar, account.user, authkey), shell=True)

	output.wait()
	
	if os.path.isfile(cookiejar):
		if os.stat(cookiejar).st_size==0:
			logger.error(" [WK_WRAPPER] :: Error while cookie forging")
			exit()
		else:
			logger.debug(" [WK_WRAPPER] :: Cookie created")
	else:
		logger.error(" [WK_WRAPPER] :: Error while cookie forging")
		exit()

def clean_x():
	os.remove(XAUTH)
	if os.path.exists('/tmp/.X%s-lock' % DISP):
		pid = int(open('/tmp/.X%s-lock' % DISP).read().strip())
		os.kill(pid, signal.SIGINT)

def run(settings):
	filename 		= settings['filename']
	viewname 		= settings['viewname']
	starttime 		= settings['starttime']
	stoptime		= settings['stoptime']
	cookiejar		= settings['cookiejar']
	windowstatus	= settings['windowstatus']
	opts			= settings['opts']
	xlock			= settings['xlock']
	xvfb_cmd		= settings['xvfb_cmd']
	display_int		= settings['display_int']
	report_dir		= settings['report_dir']
	header			= settings['header']
	footer			= settings['footer']
	account			= settings['account']

	check_xorg(xlock, xvfb_cmd)
	export_env(display_int)
	check_report_dir(report_dir)
	get_cookie(cookiejar, account)

	runscript = "var export_view_id='%s';var export_from=%s;var export_to=%s" % (viewname, starttime, stoptime)
	opts = ' '.join(opts)

	logger.debug("wkhtmltopdf %s %s %s --window-status %s -T 21mm --header-line --header-spacing 5 --cookie-jar %s --run-script \"%s\" 'http://127.0.0.1:8082/static/canopsis/reporting.html' '%s/%s'" % (opts, header, footer, windowstatus, cookiejar, runscript, report_dir, filename))

	result = Popen("wkhtmltopdf %s %s %s --window-status %s -T 21mm --header-line --header-spacing 5 --cookie-jar %s --run-script \"%s\" 'http://127.0.0.1:8082/static/canopsis/reporting.html' '%s/%s'" % (opts, header, footer, windowstatus, cookiejar, runscript, report_dir, filename), shell=True)

	result.wait()
	clean_x()

	return result
