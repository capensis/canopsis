Feature: test a declare ticket rule
  I need to be able to test a declare ticket rule
  Only admin should be able to test a declare ticket rule

  @concurrent
  Scenario: given test request and no auth user should not allow access
    When I do POST /api/v4/cat/test-declare-ticket-executions
    Then the response code should be 401

  @concurrent
  Scenario: given test request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/test-declare-ticket-executions
    Then the response code should be 403

  @concurrent
  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "webhooks": "Webhooks is missing.",
        "alarms": "Alarms is missing."
      }
    }
    """
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "webhooks": [
        {},
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "webhooks.0.request.url": "URL is missing.",
        "webhooks.0.request.method": "Method is missing.",
        "webhooks.1.request.url": "URL is missing.",
        "webhooks.1.request.method": "Method is missing."
      }
    }
    """
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "webhooks": [
        {
          "declare_ticket": {}
        },
        {
          "declare_ticket": {}
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "webhooks.0.declare_ticket.ticket_id": "TicketID is missing.",
        "webhooks.1.declare_ticket.ticket_id": "TicketID is missing."
      }
    }
    """
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "webhooks": [
        {
          "declare_ticket": {
            "ticket_id": "_id"
          }
        },
        {
          "declare_ticket": {
            "ticket_id": "_id"
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "webhooks.1.declare_ticket.ticket_id": "TicketID is not empty."
      }
    }
    """

  @concurrent
  Scenario: given not exist alarm request should return error
    When I am admin
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "test-alarm-not-exist"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-fail-2-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-fail-2-system-name",
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-fail-2\",\"url\":\"https://test/test-ticket-declareticket-test-execution-fail-2\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
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
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarms": "Alarms don't exist."
      }
    }
    """

  @concurrent
  Scenario: given resolved alarm should return error
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-test-execution-fail-5",
      "connector_name" : "test-connector-name-declareticket-test-execution-fail-5",
      "component" :  "test-component-declareticket-test-execution-fail-5",
      "resource" : "test-resource-declareticket-test-execution-fail-5",
      "source_type" : "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-declareticket-test-execution-fail-5",
      "connector_name" : "test-connector-name-declareticket-test-execution-fail-5",
      "component" :  "test-component-declareticket-test-execution-fail-5",
      "resource" : "test-resource-declareticket-test-execution-fail-5",
      "source_type" : "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "resolve_cancel",
      "connector" : "test-connector-declareticket-test-execution-fail-5",
      "connector_name" : "test-connector-name-declareticket-test-execution-fail-5",
      "component" :  "test-component-declareticket-test-execution-fail-5",
      "resource" : "test-resource-declareticket-test-execution-fail-5",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-test-execution-fail-5
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/test-declare-ticket-executions:
    """json
    {
      "alarms": [
        "{{ .alarmId }}"
      ],
      "name": "test-declareticketrule-declareticket-test-execution-fail-5-name",
      "system_name": "test-declareticketrule-declareticket-test-execution-fail-5-system-name",
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-test-execution-fail-5\",\"url\":\"https://test/test-ticket-declareticket-test-execution-fail-5\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
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
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarms": "Alarms don't exist."
      }
    }
    """
