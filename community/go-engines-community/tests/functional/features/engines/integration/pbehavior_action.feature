Feature: no execute action when entity is inactive
  I need to be able to disable action when pause or maintenance pbehavior is in action.

  Scenario: given action and maintenance pbehavior should not update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-1-name",
      "enabled": true,
      "priority": 71,
      "triggers": ["stateinc"],
      "disable_during_periods": ["maintenance"],
      "actions": [
        {
          "entity_patterns": [{"name": "test-resource-pbehavior-action-1"}],
          "type": "ack",
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
      "connector" : "test-connector-pbehavior-action-1",
      "connector_name" : "test-connector-name-pbehavior-action-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-1",
      "resource" : "test-resource-pbehavior-action-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-action-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-action-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-action-1",
      "connector_name" : "test-connector-name-pbehavior-action-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-1",
      "resource" : "test-resource-pbehavior-action-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.ack":{"$exists":false}},{"v.resource":"test-resource-pbehavior-action-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-action-1",
            "connector_name" : "test-connector-name-pbehavior-action-1",
            "component" : "test-component-pbehavior-action-1",
            "resource" : "test-resource-pbehavior-action-1",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": ""
              },
              {"_t": "stateinc"}
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

  Scenario: given delayed action and maintenance pbehavior should update alarm after pbehavior
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-2-name",
      "enabled": true,
      "priority": 72,
      "triggers": ["create"],
      "delay": {
        "value": 10,
        "unit": "s"
      },
      "actions": [
        {
          "entity_patterns": [{"name": "test-resource-pbehavior-action-2"}],
          "type": "ack",
          "parameters": {"output": "test ack"},
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-action-2",
      "tstart": {{ nowAdd "5s" }},
      "tstop": {{ nowAdd "10s" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-action-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-action-2",
      "connector_name" : "test-connector-name-pbehavior-action-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-2",
      "resource" : "test-resource-pbehavior-action-2",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I wait the end of event processing
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-action-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-action-2",
            "connector_name" : "test-connector-name-pbehavior-action-2",
            "component" : "test-component-pbehavior-action-2",
            "resource" : "test-resource-pbehavior-action-2",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": ""
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": ""
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-action-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "system"
            },
            "connector" : "test-connector-pbehavior-action-2",
            "connector_name" : "test-connector-name-pbehavior-action-2",
            "component" : "test-component-pbehavior-action-2",
            "resource" : "test-resource-pbehavior-action-2",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": ""
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": ""
              },
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

  Scenario: given pbehavior action should create pbehavior and update new alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-3-name",
      "enabled": true,
      "priority": 73,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns": [{"name": "test-resource-pbehavior-action-3"}],
          "parameters": {
            "name": "pbehavior-action-3",
            "tstart": {{ now }},
            "tstop": {{ nowAdd "1h" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
          },
          "type": "pbehavior",
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
      "connector" : "test-connector-pbehavior-action-3",
      "connector_name" : "test-connector-name-pbehavior-action-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-3",
      "resource" : "test-resource-pbehavior-action-3",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-action-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-action-3",
            "connector_name" : "test-connector-name-pbehavior-action-3",
            "component" : "test-component-pbehavior-action-3",
            "resource" : "test-resource-pbehavior-action-3",
            "pbehavior_info": {
              "canonical_type": "maintenance"
            },
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": ""
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

  Scenario: given pbehavior action with start on trigger should create pbehavior
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-4-name",
      "enabled": true,
      "priority": 74,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns": [{"name": "test-resource-pbehavior-action-4"}],
          "parameters": {
            "name": "pbehavior-action-4",
            "start_on_trigger": true,
            "duration": {
              "value": 1,
              "unit": "h"
            },
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
          },
          "type": "pbehavior",
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
      "connector" : "test-connector-pbehavior-action-4",
      "connector_name" : "test-connector-name-pbehavior-action-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-4",
      "resource" : "test-resource-pbehavior-action-4",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-action-4"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-action-4",
            "connector_name" : "test-connector-name-pbehavior-action-4",
            "component" : "test-component-pbehavior-action-4",
            "resource" : "test-resource-pbehavior-action-4",
            "pbehavior_info": {
              "canonical_type": "maintenance"
            },
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": ""
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

  Scenario: given pbehavior action should create pbehavior and update old alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-5-name",
      "enabled": true,
      "priority": 75,
      "triggers": ["stateinc"],
      "actions": [
        {
          "_id": "test-action-pbehavior-action-5",
          "enabled": true,
          "entity_patterns": [{"name": "test-resource-pbehavior-action-5"}],
          "parameters": {
            "name": "pbehavior-action-5",
            "tstart": {{ now }},
            "tstop": {{ nowAdd "1h" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
          },
          "type": "pbehavior",
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
      "connector" : "test-connector-pbehavior-action-5",
      "connector_name" : "test-connector-name-pbehavior-action-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-5",
      "resource" : "test-resource-pbehavior-action-5",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-action-5",
      "connector_name" : "test-connector-name-pbehavior-action-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-5",
      "resource" : "test-resource-pbehavior-action-5",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-action-5"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-action-5",
            "connector_name" : "test-connector-name-pbehavior-action-5",
            "component" : "test-component-pbehavior-action-5",
            "resource" : "test-resource-pbehavior-action-5",
            "pbehavior_info": {
              "canonical_type": "maintenance"
            },
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "stateinc"},
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": ""
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

  Scenario: given pbehavior action should create pbehavior and update last alarm date of pbehavior
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-6-name",
      "enabled": true,
      "priority": 76,
      "triggers": ["stateinc"],
      "actions": [
        {
          "_id": "test-action-pbehavior-action-6",
          "enabled": true,
          "entity_patterns": [{"name": "test-resource-pbehavior-action-6"}],
          "parameters": {
            "name": "pbehavior-action-6",
            "tstart": {{ now }},
            "tstop": {{ nowAdd "1h" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
          },
          "type": "pbehavior",
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
      "connector" : "test-connector-pbehavior-action-6",
      "connector_name" : "test-connector-name-pbehavior-action-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-6",
      "resource" : "test-resource-pbehavior-action-6",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-action-6",
      "connector_name" : "test-connector-name-pbehavior-action-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-action-6",
      "resource" : "test-resource-pbehavior-action-6",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/pbehaviors?search=pbehavior-action-6
    Then the response code should be 200
    Then the response key "data.0.last_alarm_date" should not be "null"
