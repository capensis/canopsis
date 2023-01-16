Feature: get a declare ticket status 
  I need to be able to get an execution status

  @concurrent
  Scenario: given declare ticket execution should get success status from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-websocket-1",
      "connector_name" : "test-connector-name-declareticket-execution-websocket-1",
      "component" :  "test-component-declareticket-execution-websocket-1",
      "resource" : "test-resource-declareticket-execution-websocket-1",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-websocket-1
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-websocket-1-name",
      "system_name": "test-declareticketrule-declareticket-execution-websocket-1-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-execution-websocket-1 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-websocket-1\",\"url\":\"https://test/test-ticket-declareticket-execution-websocket-1\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-declareticket-execution-websocket-1"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-execution:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId }}"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    Then I save response executionID={{ (index .lastResponse 0).id }}
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 2,
      "fail_reason": ""
    }
    """

  @concurrent
  Scenario: given declare ticket execution with failed last webhook should get fail status from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-websocket-2",
      "connector_name" : "test-connector-name-declareticket-execution-websocket-2",
      "component" :  "test-component-declareticket-execution-websocket-2",
      "resource" : "test-resource-declareticket-execution-websocket-2",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-websocket-2
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-websocket-2-name",
      "system_name": "test-declareticketrule-declareticket-execution-websocket-2-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/auth-request",
            "method": "POST"
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-websocket-2\",\"url\":\"https://test/test-ticket-declareticket-execution-websocket-2\",\"name\":\"test-ticket-declareticket-execution-websocket-2 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-declareticket-execution-websocket-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-execution:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId }}"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    Then I save response executionID={{ (index .lastResponse 0).id }}
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 3,
      "fail_reason": "url {{ .dummyApiURL }}/webhook/auth-request is unauthorized"
    }
    """

  @concurrent
  Scenario: given declare ticket execution with failed declare ticket webhook should get fail status from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-websocket-3",
      "connector_name" : "test-connector-name-declareticket-execution-websocket-3",
      "component" :  "test-component-declareticket-execution-websocket-3",
      "resource" : "test-resource-declareticket-execution-websocket-3",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-websocket-3
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-websocket-3-name",
      "system_name": "test-declareticketrule-declareticket-execution-websocket-3-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-websocket-3\",\"url\":\"https://test/test-ticket-declareticket-execution-websocket-3\",\"name\":\"test-ticket-declareticket-execution-websocket-3 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "declare_ticket": {
            "ticket_id": "not_exist_field"
          },
          "stop_on_fail": false
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "method": "POST"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-declareticket-execution-websocket-3"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-execution:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId }}"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    Then I save response executionID={{ (index .lastResponse 0).id }}
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 3,
      "fail_reason": "ticket_id is emtpy, response has nothing in not_exist_field"
    }
    """
