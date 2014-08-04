component = "test_event_32"
import time
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
	wait = 60
	print 'wait %s seconds before sending OK event...' % wait
	time.sleep(wait)
	return ('event OK after 60 seconds',{
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

scenario = [e1,e2,e3]