Feature: update alarm by RPC stream
  I need to be able to update alarm on RPC event

  @concurrent
  Scenario: given ack event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-1",
      "connector_name" : "test-connector-name-axe-rpc-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-1",
      "resource" : "test-resource-axe-rpc-1",
      "state" : 2,
      "output" : "test-output-axe-rpc-1"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-1/test-component-axe-rpc-1:
	"""
	{
		"event_type": "ack",
		"parameters": {
          "output": "test-output-axe-rpc-1"
        }
	}
	"""
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "system",
              "m": "test-output-axe-rpc-1",
              "val": 0
            },
            "component": "test-component-axe-rpc-1",
            "connector": "test-connector-axe-rpc-1",
            "connector_name": "test-connector-name-axe-rpc-1",
            "resource": "test-resource-axe-rpc-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "m": "test-output-axe-rpc-1",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given remove ack event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-2",
      "connector_name" : "test-connector-name-axe-rpc-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-2",
      "resource" : "test-resource-axe-rpc-2",
      "state" : 2,
      "output" : "test-output-axe-rpc-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-rpc-2",
      "connector_name" : "test-connector-name-axe-rpc-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-2",
      "resource" : "test-resource-axe-rpc-2",
      "output" : "test-output-axe-rpc-2"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-2/test-component-axe-rpc-2:
	"""
	{
		"event_type": "ackremove",
		"parameters": {
          "output": "test-output-axe-rpc-2"
        }
	}
	"""
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-rpc-2",
            "connector": "test-connector-axe-rpc-2",
            "connector_name": "test-connector-name-axe-rpc-2",
            "resource": "test-resource-axe-rpc-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "a": "root",
                "m": "test-output-axe-rpc-2",
                "val": 0
              },
              {
                "_t": "ackremove",
                "a": "system",
                "m": "test-output-axe-rpc-2",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given cancel event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-3",
      "connector_name" : "test-connector-name-axe-rpc-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-3",
      "resource" : "test-resource-axe-rpc-3",
      "state" : 2,
      "output" : "test-output-axe-rpc-3"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-3/test-component-axe-rpc-3:
	"""
	{
		"event_type": "cancel",
		"parameters": {
          "output": "test-output-axe-rpc-3"
        }
	}
	"""
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "canceled": {
              "_t": "cancel",
              "a": "system",
              "m": "test-output-axe-rpc-3",
              "val": 0
            },
            "component": "test-component-axe-rpc-3",
            "connector": "test-connector-axe-rpc-3",
            "connector_name": "test-connector-name-axe-rpc-3",
            "resource": "test-resource-axe-rpc-3",
            "state": {
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-axe-rpc-3.test-connector-name-axe-rpc-3",
              "m": "test-output-axe-rpc-3",
              "val": 4
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "m": "test-output-axe-rpc-3",
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-rpc-3.test-connector-name-axe-rpc-3",
                "m": "test-output-axe-rpc-3",
                "val": 4
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given assoc ticket event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-4",
      "connector_name" : "test-connector-name-axe-rpc-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-4",
      "resource" : "test-resource-axe-rpc-4",
      "state" : 2,
      "output" : "test-output-axe-rpc-4"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-4/test-component-axe-rpc-4:
    """json
    {
      "event_type": "assocticket",
      "parameters": {
        "ticket": "test-ticket",
        "ticket_url": "test-url",
        "ticket_system_name": "test-system-name",
        "ticket_data": {
          "ticket_param_1": "ticket_value_1"
        },
        "ticket_comment": "test-comment"
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "m": "Ticket ID: test-ticket. Ticket URL: test-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket",
                "ticket_url": "test-url",
                "ticket_system_name": "test-system-name",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "Ticket ID: test-ticket. Ticket URL: test-url. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-ticket",
              "ticket_url": "test-url",
              "ticket_system_name": "test-system-name",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-comment"
            },
            "component": "test-component-axe-rpc-4",
            "connector": "test-connector-axe-rpc-4",
            "connector_name": "test-connector-name-axe-rpc-4",
            "resource": "test-resource-axe-rpc-4",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "val": 0,
                "m": "Ticket ID: test-ticket. Ticket URL: test-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket",
                "ticket_url": "test-url",
                "ticket_system_name": "test-system-name",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given change state event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-5",
      "connector_name" : "test-connector-name-axe-rpc-5",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-5",
      "resource" : "test-resource-axe-rpc-5",
      "state" : 2,
      "output" : "test-output-axe-rpc-5"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-5/test-component-axe-rpc-5:
	"""
	{
		"event_type": "changestate",
		"parameters": {
		  "state": 3,
          "output": "test-output-axe-rpc-5"
        }
	}
	"""
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-rpc-5",
            "connector": "test-connector-axe-rpc-5",
            "connector_name": "test-connector-name-axe-rpc-5",
            "resource": "test-resource-axe-rpc-5",
            "state": {
              "_t": "changestate",
              "a": "system",
              "m": "test-output-axe-rpc-5",
              "val": 3
            },
            "status": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "m": "test-output-axe-rpc-5",
                "val": 3
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given snooze event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-6",
      "connector_name" : "test-connector-name-axe-rpc-6",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-6",
      "resource" : "test-resource-axe-rpc-6",
      "state" : 2,
      "output" : "test-output-axe-rpc-6"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-6/test-component-axe-rpc-6:
	"""
	{
		"event_type": "snooze",
		"parameters": {
		  "duration": {
		    "value": 1,
		    "unit": "h"
		  },
          "output": "test-output-axe-rpc-6"
        }
	}
	"""
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "snooze": {
              "_t": "snooze",
              "a": "system",
              "m": "test-output-axe-rpc-6"
            },
            "component": "test-component-axe-rpc-6",
            "connector": "test-connector-axe-rpc-6",
            "connector_name": "test-connector-name-axe-rpc-6",
            "resource": "test-resource-axe-rpc-6",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "m": "test-output-axe-rpc-6"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given change state event with ok status should update alarm status
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-7",
      "connector_name" : "test-connector-name-axe-rpc-7",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-7",
      "resource" : "test-resource-axe-rpc-7",
      "state" : 2,
      "output" : "test-output-axe-rpc-7"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-7/test-component-axe-rpc-7:
	"""
	{
		"event_type": "changestate",
		"parameters": {
		  "state": 0,
          "output": "test-output-axe-rpc-7"
        }
	}
	"""
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-rpc-7",
            "connector": "test-connector-axe-rpc-7",
            "connector_name": "test-connector-name-axe-rpc-7",
            "resource": "test-resource-axe-rpc-7",
            "state": {
              "_t": "changestate",
              "a": "system",
              "m": "test-output-axe-rpc-7",
              "val": 0
            },
            "status": {
              "_t": "statusdec",
              "a": "system",
              "m": "test-output-axe-rpc-7",
              "val": 0
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "m": "test-output-axe-rpc-7",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "system",
                "m": "test-output-axe-rpc-7",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given change state event with ok status should not update alarm anymore
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-8",
      "connector_name" : "test-connector-name-axe-rpc-8",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-8",
      "resource" : "test-resource-axe-rpc-8",
      "state" : 2,
      "output" : "test-output-axe-rpc-8"
    }
    """
    When I call RPC to engine-axe with alarm test-resource-axe-rpc-8/test-component-axe-rpc-8:
	"""
	{
		"event_type": "changestate",
		"parameters": {
		  "state": 0,
          "output": "test-output-axe-rpc-8"
        }
	}
	"""
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-rpc-8",
      "connector_name" : "test-connector-name-axe-rpc-8",
      "source_type" : "resource",
      "component" :  "test-component-axe-rpc-8",
      "resource" : "test-resource-axe-rpc-8",
      "state" : 3,
      "output" : "test-output-axe-rpc-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-rpc-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-rpc-8",
            "connector": "test-connector-axe-rpc-8",
            "connector_name": "test-connector-name-axe-rpc-8",
            "resource": "test-resource-axe-rpc-8",
            "state": {
              "_t": "changestate",
              "a": "system",
              "m": "test-output-axe-rpc-8",
              "val": 0
            },
            "status": {
              "_t": "statusdec",
              "a": "system",
              "m": "test-output-axe-rpc-8",
              "val": 0
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "m": "test-output-axe-rpc-8",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "system",
                "m": "test-output-axe-rpc-8",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """
