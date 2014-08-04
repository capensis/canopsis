component = "test_event_32"
def e1(previousEvent):

	return {
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
	}

def e2(previousEvent):
	#This is an ack
	event = previousEvent
	return {
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
	}

scenario = [e1,e2]