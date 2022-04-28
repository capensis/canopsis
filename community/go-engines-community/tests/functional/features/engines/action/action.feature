Feature: execute action on trigger
  I need to be able to trigger action on event

  Scenario: given scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-1-name",
      "enabled": true,
      "priority": 20,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "resource": "test-resource-action-1"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-1/test-component-action-1"
            }
          ],
          "type": "assocticket",
          "parameters": {
            "output": "test-scenario-action-1-action-1-output {{ `{{ .Alarm.Value.Connector }}` }}",
            "ticket": "test-scenario-action-1-action-1-ticket"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-1"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-1/test-component-action-1"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-scenario-action-1-action-2-output {{ `{{ .Alarm.Value.Connector }}` }}",
            "author": "test-scenario-action-1-action-2-author"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-action-1",
      "connector_name" : "test-connector-name-action-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-1",
      "resource" : "test-resource-action-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-action-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "assocticket",
              "a": "root",
              "user_id": "root",
              "m": "test-scenario-action-1-action-1-ticket",
              "val": "test-scenario-action-1-action-1-ticket"
            },
            "ack": {
              "_t": "ack",
              "a": "test-scenario-action-1-action-2-author",
              "user_id": "root",
              "m": "test-scenario-action-1-action-2-output test-connector-action-1"
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "root",
                "user_id": "root",
                "m": "test-scenario-action-1-action-1-ticket"
              },
              {
                "_t": "ack",
                "a": "test-scenario-action-1-action-2-author",
                "user_id": "root",
                "m": "test-scenario-action-1-action-2-output test-connector-action-1"
              }
            ],
            "connector": "test-connector-action-1",
            "connector_name": "test-connector-name-action-1",
            "component": "test-component-action-1",
            "resource": "test-resource-action-1"
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

  Scenario: given delayed scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-2-name",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "delay": {
        "value": 5,
        "unit": "s"
      },
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "resource": "test-resource-action-2"
              }
            }
          ],
          "type": "assocticket",
          "parameters": {
            "output": "test-output-action-2-{{ `{{ .Alarm.Value.Connector }}` }}",
            "ticket": "test-ticket-action-2"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-2"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-2/test-component-action-2"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-output-action-2-{{ `{{ .Alarm.Value.Connector }}` }}"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-action-2",
      "connector_name" : "test-connector-name-action-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-2",
      "resource" : "test-resource-action-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-action-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "assocticket",
              "a": "root",
              "user_id": "root",
              "m": "test-ticket-action-2",
              "val": "test-ticket-action-2"
            },
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-output-action-2-test-connector-action-2"
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "root",
                "user_id": "root",
                "m": "test-ticket-action-2"
              },
              {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-output-action-2-test-connector-action-2"
              }
            ],
            "connector": "test-connector-action-2",
            "connector_name": "test-connector-name-action-2",
            "component": "test-component-action-2",
            "resource": "test-resource-action-2"
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

  Scenario: given scenario with emit trigger and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-3-name-1",
      "enabled": true,
      "priority": 22,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "resource": "test-resource-action-3"
              }
            }
          ],
          "type": "assocticket",
          "parameters": {
            "output": "test-output-action-3-{{ `{{ .Alarm.Value.Connector }}` }}",
            "ticket": "test-ticket-action-3"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-3-name-2",
      "enabled": true,
      "priority": 23,
      "triggers": ["assocticket"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-3"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-3/test-component-action-3"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-output-action-3-{{ `{{ .Alarm.Value.Connector }}` }}"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-action-3",
      "connector_name" : "test-connector-name-action-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-3",
      "resource" : "test-resource-action-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-action-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "assocticket",
              "a": "root",
              "user_id": "root",
              "m": "test-ticket-action-3",
              "val": "test-ticket-action-3"
            },
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-output-action-3-test-connector-action-3"
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "root",
                "user_id": "root",
                "m": "test-ticket-action-3"
              },
              {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-output-action-3-test-connector-action-3"
              }
            ],
            "connector": "test-connector-action-3",
            "connector_name": "test-connector-name-action-3",
            "component": "test-component-action-3",
            "resource": "test-resource-action-3"
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
