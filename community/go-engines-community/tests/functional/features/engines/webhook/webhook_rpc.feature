Feature: execute request on event
  I need to be able to execute request defined by webhook on RPC event

  Scenario: given event should execute request
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-webhook-rpc-1",
      "connector_name" : "test-connector-name-webhook-rpc-1",
      "source_type" : "resource",
      "component" :  "test-component-webhook-rpc-1",
      "resource" : "test-resource-webhook-rpc-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-webhook with alarm test-resource-webhook-rpc-1/test-component-webhook-rpc-1:
    """
    {
      "parameters": {
        "request": {
          "method": "GET",
          "url": "{{ .apiUrl }}/api/v4/scenarios",
          "auth": {
            "username": "root",
            "password": "test"
          }
        }
      }
    }
    """

  Scenario: given event with declare ticket should execute request and update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-webhook-rpc-2",
      "connector_name" : "test-connector-name-webhook-rpc-2",
      "source_type" : "resource",
      "component" :  "test-component-webhook-rpc-2",
      "resource" : "test-resource-webhook-rpc-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-webhook with alarm test-resource-webhook-rpc-2/test-component-webhook-rpc-2:
    """
    {
      "parameters": {
        "request": {
          "method": "POST",
          "url": "{{ .apiUrl }}/api/v4/scenarios",
          "auth": {
            "username": "root",
            "password": "test"
          },
          "headers": {"Content-Type": "application/json"},
          "payload": "{\"priority\": 10068,\"name\":\"test-scenario-webhook-rpc-2\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\":\"eq\",\"value\":\"test-scenario-webhook-rpc-2-alarm\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
        },
        "declare_ticket": {
          "empty_response": false,
          "is_regexp": false,
          "ticket_id": "_id",
          "scenario_name": "name"
        }
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-webhook-rpc-2
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
                "scenario_name": "test-scenario-webhook-rpc-2"
              }
            },
            "connector" : "test-connector-webhook-rpc-2",
            "connector_name" : "test-connector-name-webhook-rpc-2",
            "component" :  "test-component-webhook-rpc-2",
            "resource" : "test-resource-webhook-rpc-2",
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
                "_t": "declareticket"
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

  Scenario: given event with regexp declare ticket should execute request and update alarm
    Given I am admin
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-webhook-rpc-3",
      "connector_name" : "test-connector-name-webhook-rpc-3",
      "source_type" : "resource",
      "component" :  "test-component-webhook-rpc-3",
      "resource" : "test-resource-webhook-rpc-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I call RPC to engine-webhook with alarm test-resource-webhook-rpc-3/test-component-webhook-rpc-3:
    """
    {
      "parameters": {
        "request": {
          "method": "POST",
          "url": "{{ .apiUrl }}/api/v4/scenarios",
          "auth": {
            "username": "root",
            "password": "test"
          },
          "headers": {"Content-Type": "application/json"},
          "payload": "{\"priority\": 10069,\"name\":\"test-scenario-webhook-rpc-3\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\":\"eq\",\"value\":\"test-scenario-webhook-rpc-3-alarm\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
        },
        "declare_ticket": {
          "empty_response": false,
          "is_regexp": true,
          "ticket_id": ".*id.*",
          "scenario_name": ".*name.*"
        }
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-webhook-rpc-3
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
                "scenario_name": "test-scenario-webhook-rpc-3"
              }
            },
            "connector" : "test-connector-webhook-rpc-3",
            "connector_name" : "test-connector-name-webhook-rpc-3",
            "component" :  "test-component-webhook-rpc-3",
            "resource" : "test-resource-webhook-rpc-3",
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
                "_t": "declareticket"
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
