Feature: execute action on trigger
  I need to be able to trigger action on event

  Scenario: given scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-webhook-1-name",
      "enabled": true,
      "priority": 30,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-1"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-webhook-1/test-component-action-webhook-1"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"{{ `{{ .Entity.ID }}` }}\",\"enabled\":true,\"priority\":41,\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"test-scenario-action-webhook-1-alarm\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name": "name"
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
      "connector" : "test-connector-action-webhook-1",
      "connector_name" : "test-connector-name-action-webhook-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-1",
      "resource" : "test-resource-action-webhook-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-action-webhook-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "data": {
                "scenario_name": "test-resource-action-webhook-1/test-component-action-webhook-1"
              }
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "declareticket"
              }
            ],
            "connector": "test-connector-action-webhook-1",
            "connector_name": "test-connector-name-action-webhook-1",
            "component": "test-component-action-webhook-1",
            "resource": "test-resource-action-webhook-1"
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

  Scenario: given scenario and check event should emit declare ticket trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-webhook-2-1-name",
      "enabled": true,
      "priority": 31,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-2"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-webhook-2/test-component-action-webhook-2"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"test-scenario-action-webhook-2\",\"enabled\":true,\"priority\":42,\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"test-scenario-action-webhook-2-alarm\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name": "name"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-webhook-2-2-name",
      "enabled": true,
      "priority": 32,
      "triggers": ["declareticketwebhook"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-2"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-webhook-2/test-component-action-webhook-2"
            }
          ],
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
    """
    {
      "connector" : "test-connector-action-webhook-2",
      "connector_name" : "test-connector-name-action-webhook-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-2",
      "resource" : "test-resource-action-webhook-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-action-webhook-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket"
            },
            "ack": {
              "_t": "ack"
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "declareticket"
              },
              {
                "_t": "ack"
              }
            ],
            "connector": "test-connector-action-webhook-2",
            "connector_name": "test-connector-name-action-webhook-2",
            "component": "test-component-action-webhook-2",
            "resource": "test-resource-action-webhook-2"
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

  Scenario: given scenario and ack resources event should update resource alarms
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-webhook-3-name",
      "enabled": true,
      "priority": 33,
      "triggers": ["ack"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-3"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-component-action-webhook-3"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"test-scenario-action-webhook-3\",\"enabled\":true,\"priority\":43,\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"test-scenario-action-webhook-3-alarm\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name": "name"
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
      "connector" : "test-connector-action-webhook-3",
      "connector_name" : "test-connector-name-action-webhook-3",
      "source_type" : "component",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" : "test-connector-action-webhook-3",
      "connector_name" : "test-connector-name-action-webhook-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-3",
      "resource" : "test-resource-action-webhook-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" : "test-connector-action-webhook-3",
      "connector_name" : "test-connector-name-action-webhook-3",
      "source_type" : "component",
      "event_type" : "ack",
      "component" :  "test-component-action-webhook-3",
      "ack_resources": true,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-component-action-webhook-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket"
            },
            "ack": {
              "_t": "ack"
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "ack"
              },
              {
                "_t": "declareticket"
              }
            ],
            "connector": "test-connector-action-webhook-3",
            "connector_name": "test-connector-name-action-webhook-3",
            "component": "test-component-action-webhook-3"
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
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-action-webhook-3/test-component-action-webhook-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket"
            },
            "ack": {
              "_t": "ack"
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "ack"
              },
              {
                "_t": "declareticket"
              }
            ],
            "connector": "test-connector-action-webhook-3",
            "connector_name": "test-connector-name-action-webhook-3",
            "component": "test-component-action-webhook-3",
            "resource": "test-resource-action-webhook-3"
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

  Scenario: given webhook scenario to test response and header templates
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-webhook-4-name",
      "enabled": true,
      "priority": 34,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-4"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-webhook-4/test-component-action-webhook-4"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "GET",
              "url": "http://api:8082/api/v4/pbehavior-types/test-default-active-type",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"}
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-4"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-webhook-4/test-component-action-webhook-4"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "{{ `{{index .Header \"Content-Type\"}}` }}"},
              "payload": "{\"name\":\"test-scenario-action-webhook-4-webhook {{ `{{index .Response \"icon_name\"}}` }}\",\"enabled\":true,\"priority\":44,\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"{{ `{{index .Response \"_id\"}}` }}\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
      "connector" : "test-connector-action-webhook-4",
      "connector_name" : "test-connector-name-action-webhook-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-4",
      "resource" : "test-resource-action-webhook-4",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-4-webhook
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-4-webhook test-active-icon",
          "actions": [
            {
              "alarm_patterns": [
                {"_id": "test-default-active-type"}
              ]
            }
          ]
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

  Scenario: given scenario and declareticket event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-action-webhook-5-name",
      "enabled": true,
      "priority": 35,
      "triggers": ["declareticket"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-webhook-5"
              }
            }
          ],
          "entity_patterns": [
            {
              "_id": "test-resource-action-webhook-5/test-component-action-webhook-5"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"{{ `{{ .Entity.ID }}` }}\",\"enabled\":true,\"priority\":45,\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"test-scenario-action-webhook-5-alarm\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name": "name"
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
      "connector" : "test-connector-action-webhook-5",
      "connector_name" : "test-connector-name-action-webhook-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-5",
      "resource" : "test-resource-action-webhook-5",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" : "test-connector-action-webhook-5",
      "connector_name" : "test-connector-name-action-webhook-5",
      "source_type" : "resource",
      "event_type" : "declareticket",
      "component" :  "test-component-action-webhook-5",
      "resource" : "test-resource-action-webhook-5",
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-action-webhook-5"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "data": {
                "scenario_name": "test-resource-action-webhook-5/test-component-action-webhook-5"
              }
            },
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "declareticket"
              }
            ],
            "connector": "test-connector-action-webhook-5",
            "connector_name": "test-connector-name-action-webhook-5",
            "component": "test-component-action-webhook-5",
            "resource": "test-resource-action-webhook-5"
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