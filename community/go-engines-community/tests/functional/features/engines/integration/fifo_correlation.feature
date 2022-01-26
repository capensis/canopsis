Feature: create and update meta alarm and child alarm in order
  I need to be able to create and update meta alarm

  Scenario: given child event and meta alarm event should process events in order
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-fifo-correlation-1",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-fifo-correlation-1"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-fifo-correlation-1",
      "connector_name": "test-connector-name-fifo-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-fifo-correlation-1",
      "resource": "test-resource-fifo-correlation-1",
      "state": 2,
      "output": "test-output-fifo-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-fifo-correlation-1",
        "connector_name": "test-connector-name-fifo-correlation-1",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "test-component-fifo-correlation-1",
        "resource": "test-resource-fifo-correlation-1",
        "output": "test-output-fifo-correlation-1"
      },
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-fifo-correlation-1-meta"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-fifo-correlation-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-fifo-correlation-1",
            "connector": "test-connector-fifo-correlation-1",
            "connector_name": "test-connector-name-fifo-correlation-1",
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
              },
              {
                "_t": "comment",
                "m": "test-output-fifo-correlation-1"
              },
              {
                "_t": "ack",
                "m": "test-output-fifo-correlation-1-meta"
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

  Scenario: given meta alarm event and child event should process events in order
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-fifo-correlation-2",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-fifo-correlation-2"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-fifo-correlation-2",
      "connector_name": "test-connector-name-fifo-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-fifo-correlation-2",
      "resource": "test-resource-fifo-correlation-2",
      "state": 2,
      "output": "test-output-fifo-correlation-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&correlation=true
    Then the response code should be 200
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-fifo-correlation-2-meta"
      },
      {
        "connector": "test-connector-fifo-correlation-2",
        "connector_name": "test-connector-name-fifo-correlation-2",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "test-component-fifo-correlation-2",
        "resource": "test-resource-fifo-correlation-2",
        "output": "test-output-fifo-correlation-2"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-fifo-correlation-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-fifo-correlation-2",
            "connector": "test-connector-fifo-correlation-2",
            "connector_name": "test-connector-name-fifo-correlation-2",
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
              },
              {
                "_t": "ack",
                "m": "test-output-fifo-correlation-2-meta"
              },
              {
                "_t": "comment",
                "m": "test-output-fifo-correlation-2"
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

  Scenario: given multiple child events and meta alarm event should process child events first
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-fifo-correlation-3",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-fifo-correlation-3"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-fifo-correlation-3",
      "connector_name": "test-connector-name-fifo-correlation-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-fifo-correlation-3",
      "resource": "test-resource-fifo-correlation-3",
      "state": 2,
      "output": "test-output-fifo-correlation-3"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&correlation=true
    Then the response code should be 200
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-fifo-correlation-3",
        "connector_name": "test-connector-name-fifo-correlation-3",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "test-component-fifo-correlation-3",
        "resource": "test-resource-fifo-correlation-3",
        "output": "test-output-fifo-correlation-3-1"
      },
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-fifo-correlation-3-meta"
      },
      {
        "connector": "test-connector-fifo-correlation-3",
        "connector_name": "test-connector-name-fifo-correlation-3",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "test-component-fifo-correlation-3",
        "resource": "test-resource-fifo-correlation-3",
        "output": "test-output-fifo-correlation-3-2"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-fifo-correlation-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-fifo-correlation-3",
            "connector": "test-connector-fifo-correlation-3",
            "connector_name": "test-connector-name-fifo-correlation-3",
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
              },
              {
                "_t": "comment",
                "m": "test-output-fifo-correlation-3-1"
              },
              {
                "_t": "comment",
                "m": "test-output-fifo-correlation-3-2"
              },
              {
                "_t": "ack",
                "m": "test-output-fifo-correlation-3-meta"
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

  Scenario: given child event and multiple meta alarm events should process child event first
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-fifo-correlation-4",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-fifo-correlation-4"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-fifo-correlation-4",
      "connector_name": "test-connector-name-fifo-correlation-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-fifo-correlation-4",
      "resource": "test-resource-fifo-correlation-4",
      "state": 2,
      "output": "test-output-fifo-correlation-4"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&correlation=true
    Then the response code should be 200
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-fifo-correlation-4-meta"
      },
      {
        "connector": "test-connector-fifo-correlation-4",
        "connector_name": "test-connector-name-fifo-correlation-4",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "test-component-fifo-correlation-4",
        "resource": "test-resource-fifo-correlation-4",
        "output": "test-output-fifo-correlation-4"
      },
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-fifo-correlation-4-meta"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-fifo-correlation-4"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-fifo-correlation-4",
            "connector": "test-connector-fifo-correlation-4",
            "connector_name": "test-connector-name-fifo-correlation-4",
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
              },
              {
                "_t": "ack",
                "m": "test-output-fifo-correlation-4-meta"
              },
              {
                "_t": "comment",
                "m": "test-output-fifo-correlation-4"
              },
              {
                "_t": "comment",
                "m": "test-output-fifo-correlation-4-meta"
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

#  todo race condition
#  Scenario: given multiple child events and events for multiple meta alarms should process child event first
#    Given I am admin
#    When I do POST /api/v4/cat/metaalarmrules:
#    """json
#    {
#      "name": "test-metaalarmrule-fifo-correlation-5-1",
#      "type": "attribute",
#      "config": {
#        "alarm_patterns": [
#          {
#            "v": {
#              "component": "test-component-fifo-correlation-5"
#            }
#          }
#        ]
#      }
#    }
#    """
#    Then the response code should be 201
#    Then I save response metaAlarmRuleID1={{ .lastResponse._id }}
#    When I do POST /api/v4/cat/metaalarmrules:
#    """json
#    {
#      "name": "test-metaalarmrule-fifo-correlation-5-2",
#      "type": "attribute",
#      "config": {
#        "alarm_patterns": [
#          {
#            "v": {
#              "component": "test-component-fifo-correlation-5"
#            }
#          }
#        ]
#      }
#    }
#    """
#    Then the response code should be 201
#    Then I save response metaAlarmRuleID2={{ .lastResponse._id }}
#    When I wait the next periodical process
#    When I send an event:
#    """json
#    {
#      "connector": "test-connector-fifo-correlation-5",
#      "connector_name": "test-connector-name-fifo-correlation-5",
#      "source_type": "resource",
#      "event_type": "check",
#      "component":  "test-component-fifo-correlation-5",
#      "resource": "test-resource-fifo-correlation-5",
#      "state": 2,
#      "output": "test-output-fifo-correlation-5"
#    }
#    """
#    When I wait the end of 3 events processing
#    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID1 }}"}]}&correlation=true
#    Then the response code should be 200
#    When I save response metaAlarmConnector1={{ (index .lastResponse.data 0).v.connector }}
#    When I save response metaAlarmConnectorName1={{ (index .lastResponse.data 0).v.connector_name }}
#    When I save response metaAlarmComponent1={{ (index .lastResponse.data 0).v.component }}
#    When I save response metaAlarmResource1={{ (index .lastResponse.data 0).v.resource }}
#    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID2 }}"}]}&correlation=true
#    Then the response code should be 200
#    When I save response metaAlarmConnector2={{ (index .lastResponse.data 0).v.connector }}
#    When I save response metaAlarmConnectorName2={{ (index .lastResponse.data 0).v.connector_name }}
#    When I save response metaAlarmComponent2={{ (index .lastResponse.data 0).v.component }}
#    When I save response metaAlarmResource2={{ (index .lastResponse.data 0).v.resource }}
#    When I send an event:
#    """json
#    [
#      {
#        "connector": "{{ .metaAlarmConnector1 }}",
#        "connector_name": "{{ .metaAlarmConnectorName1 }}",
#        "source_type": "resource",
#        "event_type": "comment",
#        "component":  "{{ .metaAlarmComponent1 }}",
#        "resource": "{{ .metaAlarmResource1 }}",
#        "output": "test-output-fifo-correlation-5-1-meta"
#      },
#      {
#        "connector": "test-connector-fifo-correlation-5",
#        "connector_name": "test-connector-name-fifo-correlation-5",
#        "source_type": "resource",
#        "event_type": "comment",
#        "component":  "test-component-fifo-correlation-5",
#        "resource": "test-resource-fifo-correlation-5",
#        "output": "test-output-fifo-correlation-5-1"
#      },
#      {
#        "connector": "{{ .metaAlarmConnector2 }}",
#        "connector_name": "{{ .metaAlarmConnectorName2 }}",
#        "source_type": "resource",
#        "event_type": "comment",
#        "component":  "{{ .metaAlarmComponent2 }}",
#        "resource": "{{ .metaAlarmResource2 }}",
#        "output": "test-output-fifo-correlation-5-2-meta"
#      },
#      {
#        "connector": "test-connector-fifo-correlation-5",
#        "connector_name": "test-connector-name-fifo-correlation-5",
#        "source_type": "resource",
#        "event_type": "comment",
#        "component":  "test-component-fifo-correlation-5",
#        "resource": "test-resource-fifo-correlation-5",
#        "output": "test-output-fifo-correlation-5-2"
#      }
#    ]
#    """
#    When I wait the end of 4 events processing
#    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-fifo-correlation-5"}]}&with_steps=true
#    Then the response code should be 200
#    Then the response body should contain:
#    """json
#    {
#      "data": [
#        {
#          "v": {
#            "component": "test-component-fifo-correlation-5",
#            "connector": "test-connector-fifo-correlation-5",
#            "connector_name": "test-connector-name-fifo-correlation-5",
#            "state": {
#              "val": 2
#            },
#            "status": {
#              "val": 1
#            },
#            "steps": [
#              {
#                "_t": "stateinc",
#                "val": 2
#              },
#              {
#                "_t": "statusinc",
#                "val": 1
#              },
#              {
#                "_t": "metaalarmattach",
#                "a": "engine.correlation"
#              },
#              {
#                "_t": "metaalarmattach",
#                "a": "engine.correlation"
#              },
#              {
#                "_t": "comment",
#                "m": "test-output-fifo-correlation-5-1-meta"
#              },
#              {
#                "_t": "comment",
#                "m": "test-output-fifo-correlation-5-1"
#              },
#              {
#                "_t": "comment",
#                "m": "test-output-fifo-correlation-5-2"
#              },
#              {
#                "_t": "comment",
#                "m": "test-output-fifo-correlation-5-2-meta"
#              }
#            ]
#          }
#        }
#      ],
#      "meta": {
#        "page": 1,
#        "page_count": 1,
#        "per_page": 10,
#        "total_count": 1
#      }
#    }
#    """
