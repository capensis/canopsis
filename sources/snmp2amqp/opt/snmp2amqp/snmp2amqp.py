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

import sys, os, signal, json

sys.path.append(os.path.expanduser('~/lib/canolibs'))

from pwd import getpwnam
from cinit import cinit

DAEMON_NAME='snmp2amqp'

init 		= cinit()
logger 	= init.getLogger()

RUN = False

myamqp = None
transportDispatcher = None
mibs = {}

from camqp import camqp

import cevent

from pysnmp.carrier.asynsock.dispatch import AsynsockDispatcher
from pysnmp.carrier.asynsock.dgram import udp
from pyasn1.codec.ber import decoder
from pysnmp.proto import api

from mibtools import mib
from mibtools import severity_to_state

sys.path.append(os.path.expanduser('~/etc/'))
import snmp2amqp_conf

########################################################
#
#   Functions
#
########################################################

RUN = True
#### Connect signals
def signal_handler(signum, frame):
	logger.info("Receive signal to stop daemon...")
	global RUN
	RUN = False
	transportDispatcher.closeDispatcher()

def parse_trap(mib, trap_oid, agent, varBinds):

	notification = mib.get_notification(trap_oid)

	## Parse trap
	if notification:
		try:
			logger.info("[%s][%s] %s-%s: %s (%s)" % (agent, mib.name, notification['SEVERITY'], notification['STATE'], notification['TYPE'], trap_oid))
		except Exception, err:
			logger.error("Impossible to parse notification, check mib conversion ...")
			return None
		
		arguments = notification['ARGUMENTS']
		summary	  = notification['SUMMARY']
		
		nb_string_arg = summary.count('%s')
		
		if varBinds and nb_string_arg:
			for i in range(nb_string_arg):
				logger.debug(" + Get value %s" % i)
				value = None
				oid, components = varBinds[i]
				component = components[0]
				if component != None:
					#value = component._componentValues[0]
					for info in component._componentValues:
						if info:
							value = str(info)

					logger.debug("   + %s" % value)

				if value:
					summary = summary.replace('%s', value, 1)
							

		logger.info(" + Summary: %s" % summary)

		component = agent
		resource = mib.name
		source_type = 'resource'
		state = severity_to_state[notification['SEVERITY']]
		output = notification['TYPE']
		long_output = summary

		## convert trap to event
		event = cevent.forger(
				connector='snmp',
				connector_name=DAEMON_NAME,
				component=component,
				resource=resource,
				timestamp=None,
				source_type=source_type,
				event_type='trap',
				state=state,
				output=output,
				long_output=long_output)

		#own fields
		event['snmp_severity'] = notification['SEVERITY']
		event['snmp_state'] = notification['STATE']
		event['snmp_oid'] = trap_oid

		logger.debug("Event: %s" % event)
		## send event on amqp
		key = cevent.get_routingkey(event)						
		myamqp.publish(event, key, myamqp.exchange_name_events)
		

########################################################
#
#   Callback
#
########################################################

def cbFun(transportDispatcher, transportDomain, transportAddress, wholeMsg):
	while wholeMsg:
		msgVer = int(api.decodeMessageVersion(wholeMsg))
		if api.protoModules.has_key(msgVer):
			pMod = api.protoModules[msgVer]
		else:
			print 'Unsupported SNMP version %s' % msgVer
			return

		reqMsg, wholeMsg = decoder.decode( wholeMsg, asn1Spec=pMod.Message())

		#print 'Trap from %s[%s]:' % transportAddress

		reqPDU = pMod.apiMessage.getPDU(reqMsg)
		if reqPDU.isSameTypeWith(pMod.TrapPDU()):
			if msgVer == api.protoVersion1:
				agent		= pMod.apiTrapPDU.getAgentAddr(reqPDU).prettyPrint()
				enterprise	= str(pMod.apiTrapPDU.getEnterprise(reqPDU))
				gtrap		= str(pMod.apiTrapPDU.getGenericTrap(reqPDU))
				strap		= str(pMod.apiTrapPDU.getSpecificTrap(reqPDU))
				varBinds 	= pMod.apiTrapPDU.getVarBindList(reqPDU)
				timestamp 	= pMod.apiTrapPDU.getTimeStamp(reqPDU)

				trap_oid = enterprise + '.0.' + strap

				if enterprise in snmp2amqp_conf.blacklist_enterprise:
					logger.debug("Blacklist enterprise: '%s'." % enterprise)
					return wholeMsg

				if trap_oid in snmp2amqp_conf.blacklist_trap_oid:
					logger.debug("Blacklist trap: '%s'." % trap_oid)
					return wholeMsg

				mib = None
				try:
					mib = mibs[enterprise]
				except:	
					logger.warning("Unknown trap from '%s': %s" % (agent, trap_oid))
					logger.warning(" + Unknown enterprise '%s'" % enterprise)
					#if varBinds:
					#	for oid, components in varBinds:
					#		print "  + ", oid

				if mib:
					try:
						parse_trap(mib, trap_oid, agent, varBinds)
					except Exception, err:
						logger.error("Impossible to parse trap: %s" % err)
				

	
			#else:
			#	varBinds = pMod.apiPDU.getVarBindList(reqPDU)

				
	return wholeMsg


########################################################
#
#   Main
#
########################################################

def main():
	signal.signal(signal.SIGINT, signal_handler)
	signal.signal(signal.SIGTERM, signal_handler)
	
	# global
	global myamqp, transportDispatcher

	# Connect to amqp bus
	logger.debug("Start AMQP ...")
	myamqp = camqp()

	logger.info("Load all MIBs ...")
	for oid in snmp2amqp_conf.mibs.keys():
		mibs[oid] = mib(snmp2amqp_conf.mibs[oid])
	

	logger.info("Init SNMP listenner ...")
	transportDispatcher = AsynsockDispatcher()

	transportDispatcher.registerTransport(
		udp.domainName, udp.UdpSocketTransport().openServerMode((snmp2amqp_conf.interface, snmp2amqp_conf.port))
		)

	transportDispatcher.registerRecvCbFun(cbFun)
	transportDispatcher.jobStarted(1) # this job would never finish

	## set euid of process
	os.setuid(getpwnam('canopsis')[2])

	myamqp.start()

	logger.info("Wait SNMP traps ...")
	try:
		transportDispatcher.runDispatcher()
	except Exception, err:
		## Impossible to stop transportDispatcher properly ...
		logger.error(err)
		pass

	logger.info("Stop SNMP daemon ...") 

	logger.debug("Stop AMQP ...")
	myamqp.stop()
	myamqp.join()

if __name__ == "__main__":
	main()
