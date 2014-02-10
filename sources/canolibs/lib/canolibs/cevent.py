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

import socket, time, clogging
import re

logger = clogging.getLogger()

# Change default timeout from 1 to 3 , conflict with gunicorn
socket.setdefaulttimeout(3)

regexp_ip = re.compile("([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})")

dns_cache = {}
cache_time = 1800 #30min

def forger(		connector,
			connector_name,
			event_type,
			source_type='component',
			component=None,
			resource=None,	
			timestamp=None,
			state=0,
			state_type=1,
			output=None,
			long_output=None,
			perf_data=None,
			perf_data_array=None,
			address=None,
			domain=None,
			reverse_lookup=True,
			display_name=None,
			tags=[],
			ticket=None,
			ref_rk=None
		):

	if not timestamp:
		timestamp=int(time.time())

	if not component:
		component = socket.getfqdn()
		if not component:
			component = socket.gethostname()

	if not state:
		state = 0
		
	if not address:
		if bool(regexp_ip.match(component)):
			address = component
			if reverse_lookup:
				dns = None
				
				# get from cache
				try:
					(timestamp, dns) = dns_cache[address.replace('.', '-')]
					logger.info("Use DNS lookup from cache")
					if (timestamp + cache_time) < int(time.time()):
						logger.info(" + Cache is too old")
						del dns_cache[address.replace('.', '-')]
						dns = None
				except:
					logger.info(" + '%s' not in cache" % address)
					
				# reverse lookup
				if not dns:
					try:
						logger.info("DNS reverse lookup for '%s' ..." % address)
						dns = socket.gethostbyaddr(address)
						logger.info(" + Succes: '%s'" % dns[0])
						dns_cache[address.replace('.', '-')] = (int(time.time()), dns)
					except:
						logger.info(" + Failed");
						
				# Dns ok
				if dns:	
					# Split FQDN
					fqdn = dns[0]
					component = fqdn.split('.', 1)[0]
					if not domain:
						try:
							domain = fqdn.split('.', 1)[1]
						except:
							pass
				
				if dns:
					logger.info(" + Component: %s" % component);
					logger.info(" + Address:   %s" % address);
					logger.info(" + Domain:    %s" % domain);


	dump = {
		'connector':		connector,
		'connector_name':	connector_name,
		'event_type':		event_type,
		'source_type':		source_type,
		'component':		component,
		'resource':			resource,	
		'timestamp':		timestamp,
		'state':			state,
		'state_type':		state_type,
		'output':			output,
		'long_output':		long_output,
	}
	
	if perf_data:
		dump["perf_data"] = perf_data
	
	if perf_data_array:
		dump["perf_data_array"] = perf_data_array
	
	if address:
		dump["address"] = address
	
	if domain:
		dump["domain"] = domain
	
	if tags:
		dump["tags"] = tags
	
	if display_name:
		dump["display_name"] = display_name

	if ticket:
		dump["ticket"] = ticket

	if ref_rk:
		dump["ref_rk"] = ref_rk

	return dump

def get_routingkey(event):
	rk = "%s.%s.%s.%s.%s" % (event['connector'], event['connector_name'], event['event_type'], event['source_type'], event['component'])

	if event['resource']:
		rk += ".%s" % event['resource']

	return rk
