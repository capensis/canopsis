Feature: get a test declare ticket status
  I need to be able to test get an execution status

  @concurrent
  Scenario: given declare ticket execution should get success status from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-test-execution-websocket-1",
      "connector_name" : "test-connector-name-declareticket-test-execution-websocket-1",
      "component" :  "test-component-declareticket-test-execution-websocket-1",
      "resource" : "test-resource-declareticket-test-execution-websocket-1",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-websocket-1
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-websocket-1-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-websocket-1-system-name",
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-test-execution-websocket-1 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/long-request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-websocket-1\",\"url\":\"https://test/test-ticket-declareticket-test-execution-websocket-1\",\"name\":\"{{ `{{ .Response.name }}`}}\"}"
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ]
    }
    """
    Then the response code should be 200
    Then I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 2,
      "fail_reason": "",
      "webhooks": [
        {
          "status": 2,
          "fail_reason": ""
        },
        {
          "status": 2,
          "fail_reason": ""
        }
      ]
    }
    """
    When I do GET /api/v4/cat/declare-ticket-executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 2,
      "fail_reason": "",
      "webhooks": [
        {
          "status": 2,
          "fail_reason": ""
        },
        {
          "status": 2,
          "fail_reason": ""
        }
      ]
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
      "connector" : "test-connector-declareticket-test-execution-websocket-2",
      "connector_name" : "test-connector-name-declareticket-test-execution-websocket-2",
      "component" :  "test-component-declareticket-test-execution-websocket-2",
      "resource" : "test-resource-declareticket-test-execution-websocket-2",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-websocket-2
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-websocket-2-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-websocket-2-system-name",
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/long-auth-request",
            "method": "POST"
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-websocket-2\",\"url\":\"https://test/test-ticket-declareticket-test-execution-websocket-2\",\"name\":\"test-ticket-declareticket-test-execution-websocket-2 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}"
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ]
    }
    """
    Then the response code should be 200
    Then I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 3,
      "fail_reason": "url {{ .dummyApiURL }}/webhook/long-auth-request is unauthorized",
      "webhooks": [
        {
          "status": 3,
          "fail_reason": "url {{ .dummyApiURL }}/webhook/long-auth-request is unauthorized"
        },
        {
          "status": 4,
          "fail_reason": ""
        }
      ]
    }
    """
    When I do GET /api/v4/cat/declare-ticket-executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 3,
      "fail_reason": "url {{ .dummyApiURL }}/webhook/long-auth-request is unauthorized",
      "webhooks": [
        {
          "status": 3,
          "fail_reason": "url {{ .dummyApiURL }}/webhook/long-auth-request is unauthorized"
        },
        {
          "status": 4,
          "fail_reason": ""
        }
      ]
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
      "connector" : "test-connector-declareticket-test-execution-websocket-3",
      "connector_name" : "test-connector-name-declareticket-test-execution-websocket-3",
      "component" :  "test-component-declareticket-test-execution-websocket-3",
      "resource" : "test-resource-declareticket-test-execution-websocket-3",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-websocket-3
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-websocket-3-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-websocket-3-system-name",
      "webhooks": [
        {
          "request": {
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/long-request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-websocket-3\",\"url\":\"https://test/test-ticket-declareticket-test-execution-websocket-3\",\"name\":\"test-ticket-declareticket-test-execution-websocket-3 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}"
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
      ]
    }
    """
    Then the response code should be 200
    Then I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 3,
      "fail_reason": "ticket_id is emtpy, response has nothing in not_exist_field",
      "webhooks": [
        {
          "status": 3,
          "fail_reason": "ticket_id is emtpy, response has nothing in not_exist_field"
        },
        {
          "status": 2,
          "fail_reason": ""
        }
      ]
    }
    """
    When I do GET /api/v4/cat/declare-ticket-executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 3,
      "fail_reason": "ticket_id is emtpy, response has nothing in not_exist_field",
      "webhooks": [
        {
          "status": 3,
          "fail_reason": "ticket_id is emtpy, response has nothing in not_exist_field"
        },
        {
          "status": 2,
          "fail_reason": ""
        }
      ]
    }
    """

  @concurrent
  Scenario: given declare ticket execution and unauth user should return error
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-test-execution-websocket-4",
      "connector_name" : "test-connector-name-declareticket-test-execution-websocket-4",
      "component" :  "test-component-declareticket-test-execution-websocket-4",
      "resource" : "test-resource-declareticket-test-execution-websocket-4",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-websocket-4
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-websocket-4-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-websocket-4-system-name",
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-test-execution-websocket-4 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/long-request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-websocket-4\",\"url\":\"https://test/test-ticket-declareticket-test-execution-websocket-4\",\"name\":\"{{ `{{ .Response.name }}`}}\"}"
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ]
    }
    """
    Then the response code should be 200
    Then I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 401,
      "room": "declareticket/{{ .executionID }}"
    }
    """

  @concurrent
  Scenario: given finished declare ticket execution should not subscribe to websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-test-execution-websocket-5",
      "connector_name" : "test-connector-name-declareticket-test-execution-websocket-5",
      "component" :  "test-component-declareticket-test-execution-websocket-5",
      "resource" : "test-resource-declareticket-test-execution-websocket-5",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-websocket-5
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-websocket-5-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-websocket-5-system-name",
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-test-execution-websocket-5 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/long-request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-websocket-5\",\"url\":\"https://test/test-ticket-declareticket-test-execution-websocket-5\",\"name\":\"{{ `{{ .Response.name }}`}}\"}"
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ]
    }
    """
    Then the response code should be 200
    Then I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait message from websocket room "declareticket/{{ .executionID }}" which contains:
    """json
    {
      "_id": "{{ .executionID }}",
      "status": 2,
      "fail_reason": ""
    }
    """
    When I subscribe to websocket room "declareticket/{{ .executionID }}"
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 404,
      "room": "declareticket/{{ .executionID }}"
    }
    """
