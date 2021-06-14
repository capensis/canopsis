Feature: update alarm by RPC stream
  I need to be able to update alarm on RPC event

  Scenario: given ack event should update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test_connector_axe_rpc_1",
      "connector_name" : "test_connector_name_axe_rpc_1",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_1",
      "resource" : "test_resource_axe_rpc_1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-axe with alarm test_resource_axe_rpc_1/test_component_axe_rpc_1:
	"""
	{
		"event_type": "ack",
		"parameters": {
          "output": "noveo alarm"
        }
	}
	"""
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_resource_axe_rpc_1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "system",
              "m": "noveo alarm",
              "val": 0
            },
            "component": "test_component_axe_rpc_1",
            "connector": "test_connector_axe_rpc_1",
            "connector_name": "test_connector_name_axe_rpc_1",
            "resource": "test_resource_axe_rpc_1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "noveo alarm",
                "val": 0
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given remove ack event should update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test_connector_axe_rpc_2",
      "connector_name" : "test_connector_name_axe_rpc_2",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_2",
      "resource" : "test_resource_axe_rpc_2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """
    {
      "event_type" : "ack",
      "connector" : "test_connector_axe_rpc_2",
      "connector_name" : "test_connector_name_axe_rpc_2",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_2",
      "resource" : "test_resource_axe_rpc_2",
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-axe with alarm test_resource_axe_rpc_2/test_component_axe_rpc_2:
	"""
	{
		"event_type": "ackremove",
		"parameters": {
          "output": "noveo alarm"
        }
	}
	"""
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_resource_axe_rpc_2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "test_component_axe_rpc_2",
            "connector": "test_connector_axe_rpc_2",
            "connector_name": "test_connector_name_axe_rpc_2",
            "resource": "test_resource_axe_rpc_2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "noveo alarm",
                "val": 0
              },
              {
                "_t": "ackremove",
                "a": "system",
                "m": "noveo alarm",
                "val": 0
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    Then the response key "data.0.v.ack" should not exist

  Scenario: given cancel event should update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test_connector_axe_rpc_3",
      "connector_name" : "test_connector_name_axe_rpc_3",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_3",
      "resource" : "test_resource_axe_rpc_3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-axe with alarm test_resource_axe_rpc_3/test_component_axe_rpc_3:
	"""
	{
		"event_type": "cancel",
		"parameters": {
          "output": "noveo alarm"
        }
	}
	"""
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_resource_axe_rpc_3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "canceled": {
              "_t": "cancel",
              "a": "system",
              "m": "noveo alarm",
              "val": 0
            },
            "component": "test_component_axe_rpc_3",
            "connector": "test_connector_axe_rpc_3",
            "connector_name": "test_connector_name_axe_rpc_3",
            "resource": "test_resource_axe_rpc_3",
            "state": {
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "test_connector_axe_rpc_3.test_connector_name_axe_rpc_3",
              "m": "noveo alarm",
              "val": 4
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "cancel",
                "a": "system",
                "m": "noveo alarm",
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test_connector_axe_rpc_3.test_connector_name_axe_rpc_3",
                "m": "noveo alarm",
                "val": 4
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given assoc ticket event should update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test_connector_axe_rpc_4",
      "connector_name" : "test_connector_name_axe_rpc_4",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_4",
      "resource" : "test_resource_axe_rpc_4",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-axe with alarm test_resource_axe_rpc_4/test_component_axe_rpc_4:
	"""
	{
		"event_type": "assocticket",
		"parameters": {
		  "ticket": "testticket",
          "output": "noveo alarm"
        }
	}
	"""
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_resource_axe_rpc_4"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "testticket",
              "val": "testticket"
            },
            "component": "test_component_axe_rpc_4",
            "connector": "test_connector_axe_rpc_4",
            "connector_name": "test_connector_name_axe_rpc_4",
            "resource": "test_resource_axe_rpc_4",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "assocticket",
                "a": "system",
                "m": "testticket",
                "val": 0
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given change state event should update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test_connector_axe_rpc_5",
      "connector_name" : "test_connector_name_axe_rpc_5",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_5",
      "resource" : "test_resource_axe_rpc_5",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-axe with alarm test_resource_axe_rpc_5/test_component_axe_rpc_5:
	"""
	{
		"event_type": "changestate",
		"parameters": {
		  "state": 3,
          "output": "noveo alarm"
        }
	}
	"""
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_resource_axe_rpc_5"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "test_component_axe_rpc_5",
            "connector": "test_connector_axe_rpc_5",
            "connector_name": "test_connector_name_axe_rpc_5",
            "resource": "test_resource_axe_rpc_5",
            "state": {
              "_t": "changestate",
              "a": "system",
              "m": "noveo alarm",
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "changestate",
                "a": "system",
                "m": "noveo alarm",
                "val": 3
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given snooze event should update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test_connector_axe_rpc_6",
      "connector_name" : "test_connector_name_axe_rpc_6",
      "source_type" : "resource",
      "component" :  "test_component_axe_rpc_6",
      "resource" : "test_resource_axe_rpc_6",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-axe with alarm test_resource_axe_rpc_6/test_component_axe_rpc_6:
	"""
	{
		"event_type": "snooze",
		"parameters": {
		  "duration": {
		    "seconds": 600,
		    "unit": "s"
		  },
          "output": "noveo alarm"
        }
	}
	"""
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_resource_axe_rpc_6"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "snooze": {
              "_t": "snooze",
              "a": "system",
              "m": "noveo alarm"
            },
            "component": "test_component_axe_rpc_6",
            "connector": "test_connector_axe_rpc_6",
            "connector_name": "test_connector_name_axe_rpc_6",
            "resource": "test_resource_axe_rpc_6",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "snooze",
                "a": "system",
                "m": "noveo alarm"
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
