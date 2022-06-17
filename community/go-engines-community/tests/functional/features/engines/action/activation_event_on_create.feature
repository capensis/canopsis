Feature: send activation event on create
  I need to be able to trigger rule on alarm activation

  Scenario: given event for new alarm should send event
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-action-activation-event-1",
      "connector_name" : "test-connector-name-action-activation-event-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-1",
      "resource" : "test-resource-action-activation-event-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-1"},{"v.activation_date":{"$exists":true}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-1",
            "connector_name": "test-connector-name-action-activation-event-1",
            "component": "test-component-action-activation-event-1",
            "resource": "test-resource-action-activation-event-1",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"}
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

  Scenario: given event for new alarm and ack action should send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-2-name",
      "enabled": true,
      "priority": 52,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-2"}],
          "type":"ack",
          "parameters":{},
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector" : "test-connector-action-activation-event-2",
      "connector_name" : "test-connector-name-action-activation-event-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-2",
      "resource" : "test-resource-action-activation-event-2",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-2"},{"v.activation_date":{"$exists":true}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-2",
            "connector_name": "test-connector-name-action-activation-event-2",
            "component": "test-component-action-activation-event-2",
            "resource": "test-resource-action-activation-event-2",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "ack"}
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

  Scenario: given event for new alarm and snooze action should not send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-3-name",
      "enabled": true,
      "priority": 53,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-3"}],
          "type":"snooze",
          "parameters":{
            "duration": {
              "value": 1,
              "unit": "h"
            }
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
    """
    {
      "connector" : "test-connector-action-activation-event-3",
      "connector_name" : "test-connector-name-action-activation-event-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-3",
      "resource" : "test-resource-action-activation-event-3",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-3"},{"v.activation_date":{"$exists":false}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-3",
            "connector_name": "test-connector-name-action-activation-event-3",
            "component": "test-component-action-activation-event-3",
            "resource": "test-resource-action-activation-event-3",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "snooze"}
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

  Scenario: given event for new alarm and pbehavior action with start=now should not send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-4-name",
      "enabled": true,
      "priority": 54,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-4"}],
          "type":"pbehavior",
          "parameters":{
            "name": "pbehavior-action-activation-event-4",
            "tstart": {{ now }},
            "tstop": {{ nowAdd "1h" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
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
    """
    {
      "connector" : "test-connector-action-activation-event-4",
      "connector_name" : "test-connector-name-action-activation-event-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-4",
      "resource" : "test-resource-action-activation-event-4",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-4"},{"v.activation_date":{"$exists":false}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-4",
            "connector_name": "test-connector-name-action-activation-event-4",
            "component": "test-component-action-activation-event-4",
            "resource": "test-resource-action-activation-event-4",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"}
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

  Scenario: given event for new alarm and pbehavior action with start on trigger should not send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-1-name",
      "enabled": true,
      "priority": 55,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-5"}],
          "type":"pbehavior",
          "parameters":{
            "name": "pbehavior-action-activation-event-5",
            "start_on_trigger": true,
            "duration": {
              "value": 1,
              "unit": "h"
            },
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
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
    """
    {
      "connector" : "test-connector-action-activation-event-5",
      "connector_name" : "test-connector-name-action-activation-event-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-5",
      "resource" : "test-resource-action-activation-event-5",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-5"},{"v.activation_date":{"$exists":false}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-5",
            "connector_name": "test-connector-name-action-activation-event-5",
            "component": "test-component-action-activation-event-5",
            "resource": "test-resource-action-activation-event-5",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"}
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

  Scenario: given event for new alarm and pbehavior action with start in the future should send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-6-name",
      "enabled": true,
      "priority": 56,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-6"}],
          "type":"pbehavior",
          "parameters":{
            "name": "pbehavior-action-activation-event-6",
            "tstart": {{ nowAdd "1h" }},
            "tstop": {{ nowAdd "2h" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
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
    """
    {
      "connector" : "test-connector-action-activation-event-6",
      "connector_name" : "test-connector-name-action-activation-event-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-6",
      "resource" : "test-resource-action-activation-event-6",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-6"},{"v.activation_date":{"$exists":true}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-6",
            "connector_name": "test-connector-name-action-activation-event-6",
            "component": "test-component-action-activation-event-6",
            "resource": "test-resource-action-activation-event-6",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"}
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

  Scenario: given event for new alarm and pbehavior action with start and stop in the past should send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-7-name",
      "enabled": true,
      "priority": 57,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-7"}],
          "type":"pbehavior",
          "parameters":{
            "name": "pbehavior-action-activation-event-7",
            "tstart": {{ nowAdd "-20m" }},
            "tstop": {{ nowAdd "-10m" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
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
    """
    {
      "connector" : "test-connector-action-activation-event-7",
      "connector_name" : "test-connector-name-action-activation-event-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-7",
      "resource" : "test-resource-action-activation-event-7",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-7"},{"v.activation_date":{"$exists":true}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-7",
            "connector_name": "test-connector-name-action-activation-event-7",
            "component": "test-component-action-activation-event-7",
            "resource": "test-resource-action-activation-event-7",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"}
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

  Scenario: given event for new alarm and pbehavior action with start and stop in the past and rrule should not send event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-activation-8-name",
      "enabled": true,
      "priority": 58,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns":[{"name":"test-resource-action-activation-event-8"}],
          "type":"pbehavior",
          "parameters":{
            "name": "pbehavior-action-activation-event-8",
            "tstart": {{ nowAdd "-24h" }},
            "tstop": {{ nowAdd "-23h" }},
            "rrule": "FREQ=DAILY",
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
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
    """
    {
      "connector" : "test-connector-action-activation-event-8",
      "connector_name" : "test-connector-name-action-activation-event-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-activation-event-8",
      "resource" : "test-resource-action-activation-event-8",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-activation-event-8"},{"v.activation_date":{"$exists":false}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-activation-event-8",
            "connector_name": "test-connector-name-action-activation-event-8",
            "component": "test-component-action-activation-event-8",
            "resource": "test-resource-action-activation-event-8",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"}
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
