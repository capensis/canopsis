Feature: run a declare ticket rule
  I need to be able to run a declare ticket rule
  Only admin should be able to run a declare ticket rule

  @concurrent
  Scenario: given not exist rule request should return error
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-fail-1",
      "connector_name" : "test-connector-name-declareticket-execution-fail-1",
      "component" :  "test-component-declareticket-execution-fail-1",
      "resource" : "test-resource-declareticket-execution-fail-1",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-fail-1
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "test-declareticketrule-not-exist",
        "alarms": [
          "{{ .alarmId }}"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 400,
        "errors": {
          "_id": "Rule doesn't exist."
        },
        "item": {
          "_id": "test-declareticketrule-not-exist",
          "alarms": [
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """

  @concurrent
  Scenario: given not exist alarm request should return error
    When I am admin
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-fail-2-name",
      "system_name": "test-declareticketrule-declareticket-execution-fail-2-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-fail-2\",\"url\":\"https://test/test-ticket-declareticket-execution-fail-2\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
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
                "test-resource-declareticket-execution-fail-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "test-alarm-not-exist"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 400,
        "errors": {
          "alarms": "Alarms don't exist."
        },
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "test-alarm-not-exist"
          ]
        }
      }
    ]
    """

  @concurrent
  Scenario: given not matched alarm should return error
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-fail-3",
      "connector_name" : "test-connector-name-declareticket-execution-fail-3",
      "component" :  "test-component-declareticket-execution-fail-3",
      "resource" : "test-resource-declareticket-execution-fail-3",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-fail-3
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-fail-3-name",
      "system_name": "test-declareticketrule-declareticket-execution-fail-3-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-execution-fail-3 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-fail-3\",\"url\":\"https://test/test-ticket-declareticket-execution-fail-3\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
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
                "test-resource-not-exist"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
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
    Then the response body should be:
    """json
    [
      {
        "status": 400,
         "errors": {
          "alarms.0": "Alarm doesn't match to rule."
        },
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """

  @concurrent
  Scenario: given resolved alarm should return error
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-fail-4",
      "connector_name" : "test-connector-name-declareticket-execution-fail-4",
      "component" :  "test-component-declareticket-execution-fail-4",
      "resource" : "test-resource-declareticket-execution-fail-4",
      "source_type" : "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-declareticket-execution-fail-4",
      "connector_name" : "test-connector-name-declareticket-execution-fail-4",
      "component" :  "test-component-declareticket-execution-fail-4",
      "resource" : "test-resource-declareticket-execution-fail-4",
      "source_type" : "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "resolve_cancel",
      "connector" : "test-connector-declareticket-execution-fail-4",
      "connector_name" : "test-connector-name-declareticket-execution-fail-4",
      "component" :  "test-component-declareticket-execution-fail-4",
      "resource" : "test-resource-declareticket-execution-fail-4",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-fail-4
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-fail-4-name",
      "system_name": "test-declareticketrule-declareticket-execution-fail-4-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-execution-fail-4 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-fail-4\",\"url\":\"https://test/test-ticket-declareticket-execution-fail-4\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
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
                "test-resource-declareticket-execution-fail-4"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
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
    Then the response body should be:
    """json
    [
      {
        "status": 400,
        "errors": {
          "alarms": "Alarms don't exist."
        },
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """

  @concurrent
  Scenario: given start request and no auth user should not allow access
    When I do POST /api/v4/cat/bulk/declare-ticket-executions
    Then the response code should be 401

  @concurrent
  Scenario: given start request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/bulk/declare-ticket-executions
    Then the response code should be 403
