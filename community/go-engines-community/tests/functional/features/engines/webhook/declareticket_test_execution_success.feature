Feature: test a declare ticket rule
  I need to be able to test a declare ticket rule
  Only admin should be able to test a declare ticket rule

  @concurrent
  Scenario: given test declare ticket request should execute request
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-test-execution-1",
      "connector_name" : "test-connector-name-declareticket-test-execution-1",
      "component" :  "test-component-declareticket-test-execution-1",
      "resource" : "test-resource-declareticket-test-execution-1",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-1
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-1-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-1-system-name",
      "webhooks": [
        {
          "request": {
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-test-execution-1 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}"
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-1\",\"url\":\"https://test/test-ticket-declareticket-test-execution-1\",\"name\":\"{{ `{{ .Response.name }}`}}\"}"
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
    Then the response body should contain:
    """json
    {
      "status": 1,
      "fail_reason": "",
      "webhooks": [
        {
          "status": 0,
          "fail_reason": ""
        },
        {
          "status": 0,
          "fail_reason": ""
        }
      ]
    }
    """
    Then I save response executionId={{ .lastResponse._id }}
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name. Ticket ID: test-ticket-declareticket-test-execution-1. Ticket URL: https://test/test-ticket-declareticket-test-execution-1. Ticket name: test-ticket-declareticket-test-execution-1 test-resource-declareticket-test-execution-1.",
              "ticket": "test-ticket-declareticket-test-execution-1",
              "ticket_url": "https://test/test-ticket-declareticket-test-execution-1",
              "ticket_data": {
                "name": "test-ticket-declareticket-test-execution-1 test-resource-declareticket-test-execution-1"
              },
              "ticket_comment": "Canopsis webhook test",
              "ticket_rule_id": "test",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-test-execution-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name. Ticket ID: test-ticket-declareticket-test-execution-1. Ticket URL: https://test/test-ticket-declareticket-test-execution-1. Ticket name: test-ticket-declareticket-test-execution-1 test-resource-declareticket-test-execution-1.",
                "ticket": "test-ticket-declareticket-test-execution-1",
                "ticket_url": "https://test/test-ticket-declareticket-test-execution-1",
                "ticket_data": {
                  "name": "test-ticket-declareticket-test-execution-1 test-resource-declareticket-test-execution-1"
                },
                "ticket_comment": "Canopsis webhook test",
                "ticket_rule_id": "test",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-test-execution-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-test-execution-1",
            "connector_name" : "test-connector-name-declareticket-test-execution-1",
            "component" :  "test-component-declareticket-test-execution-1",
            "resource" : "test-resource-declareticket-test-execution-1"
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
        "_id": "{{ .alarmId }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name"
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-test-execution-1-name. Ticket ID: test-ticket-declareticket-test-execution-1. Ticket URL: https://test/test-ticket-declareticket-test-execution-1. Ticket name: test-ticket-declareticket-test-execution-1 test-resource-declareticket-test-execution-1."
      }
    ]
    """
