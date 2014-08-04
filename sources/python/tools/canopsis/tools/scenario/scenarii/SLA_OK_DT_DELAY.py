import time

component = "test_event_33"
def e1(previousEvent):

	return ('event OK',{
		"connector"		: "test_connector",
		"connector_name": "test_connector_name",
		"event_type"	: "check",
		"source_type"	: "resource",
		"component"		: component,
		"resource"		: "test_custom",
		"state"			: 0,
		"state_type"	: 1,
		'hostgroups'	: ['HG3', 'HG4'],
		"PROC_CRITICAL"	: "H24",
		"PROC_WARNING"	: "H24",
		"output"		: "<h1>MESSAGE</h1>",
		"display_name"	: "DISPLAY_NAME",
		"author"		: "plop",
	})

def e2(previousEvent):

	return ('event ERROR', {
		"connector"		: "test_connector",
		"connector_name": "test_connector_name",
		"event_type"	: "check",
		"source_type"	: "resource",
		"component"		: component,
		"resource"		: "test_custom",
		"state"			: 2,
		"state_type"	: 1,
		'hostgroups'	: ['HG3', 'HG4'],
		"PROC_CRITICAL"	: "H24",
		"PROC_WARNING"	: "H24",
		"output"		: "<h1>MESSAGE</h1>",
		"display_name"	: "DISPLAY_NAME",
		"author"		: "plop",
	})

def e3(previousEvent):
	#This is an ack
	event = previousEvent
	return ('ACK',{
		'ref_rk'		:"%s.%s.%s.%s.%s" % (event['connector'], event['connector_name'], event['event_type'], event['source_type'], event['component']),
		"connector"		: "test_connector",
		"connector_name": "test_connector_name",
		"event_type"	: "ack",
		"source_type"	: "ack",
		"component"		: component,
		#"resource"		: "test_custom",
		"state"			: 0,
		"state_type"	: 1,
		"PROC_CRITICAL"	: "H24",
		"PROC_WARNING"	: "H24",
		"output"		: "aup",
		"display_name"	: "DISPLAY_NAME",
		"author"		: "plop",
	})

def e3b(previousEvent):
	return ('Downtime ...',{
		"duration"		: 60,
		"downtime_id"	: 1,
		"start"			: time.time() + 60,
		"end"			: time.time() + 70,
		"connector"		: "test_connector",
		"fixed"			: True,
		"entry"			: time.time(),
		"author"		: "eric",
		"output"		: "output",
		"connector_name": "test_connector_name",
		"event_type"	: "downtime",
		"source_type"	: "resource",
		"component"		: component,
		"resource"		: "test_custom",
		"state"			: 0,
		"state_type"	: 1,
	})

def e4(previousEvent):
	wait = 71
	print 'wait %s seconds before sending OK event...' % wait
	time.sleep(wait)
	return ('event OK',{
		"connector"		: "test_connector",
		"connector_name": "test_connector_name",
		"event_type"	: "check",
		"source_type"	: "resource",
		"component"		: component,
		"resource"		: "test_custom",
		"state"			: 0,
		"state_type"	: 1,
		'hostgroups'	: ['HG3', 'HG4'],
		"PROC_CRITICAL"	: "H24",
		"PROC_WARNING"	: "H24",
		"output"		: "<h1>MESSAGE</h1>",
		"display_name"	: "DISPLAY_NAME",
		"author"		: "plop",
	})

scenario = [e1,e2,e3,e3b,e4]