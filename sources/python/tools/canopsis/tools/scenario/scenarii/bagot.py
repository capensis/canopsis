import time

component = "bagot"
def e1(previousEvent):

	return ('event OK',{
		"connector"		: "test_connector",
		"connector_name": "test_connector_name",
		"event_type"	: "check",
		"source_type"	: "resource",
		"component"		: component,
		"resource"		: "test_custom",
		"state"			: 2,
		"state_type"	: 1,
		"output"		: "<h1>MESSAGE</h1>",
	})

def e2(previousEvent):

	return ('event KO', {
		"connector"		: "test_connector",
		"connector_name": "test_connector_name",
		"event_type"	: "check",
		"source_type"	: "resource",
		"component"		: component,
		"resource"		: "test_custom",
		"state"			: 0,
		"state_type"	: 1,
		'hostgroups'	: ['HG3', 'HG4'],
		"output"		: "<h1>MESSAGE</h1>",
	})

scenario = []

for x in xrange(10):
	scenario.append(e1)
	scenario.append(e2)