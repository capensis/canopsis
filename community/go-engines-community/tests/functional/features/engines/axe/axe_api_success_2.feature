Feature: update alarm
  I need to be able to update alarm

  @concurrent
  Scenario: given ticket resources event should update resource alarms
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-api-second-1",
      "connector": "test-connector-axe-api-second-1",
      "connector_name": "test-connector-name-axe-api-second-1",
      "component":  "test-component-axe-api-second-1",
      "resource": "test-resource-axe-api-second-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-api-second-1",
      "connector": "test-connector-axe-api-second-1",
      "connector_name": "test-connector-name-axe-api-second-1",
      "component":  "test-component-axe-api-second-1",
      "source_type": "component"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-axe-api-second-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/assocticket:
    """json
    {
      "ticket_resources": true,
      "ticket": "test-ticket-axe-api-second-1",
      "url": "test-ticket-url-axe-api-second-1",
      "system_name": "test-ticket-system-name-axe-api-second-1",
      "data": {
        "test-ticket-param-axe-api-second-1": "test-ticket-param-val-axe-api-second-1"
      },
      "comment": "test-comment-axe-api-second-1"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "assocticket",
        "connector": "api",
        "connector_name": "api",
        "component":  "test-component-axe-api-second-1",
        "source_type": "component"
      },
      {
        "event_type": "assocticket",
        "connector": "api",
        "connector_name": "api",
        "component":  "test-component-axe-api-second-1",
        "resource":  "test-resource-axe-api-second-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-second-1
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
              "initiator": "user",
              "m": "Ticket ID: test-ticket-axe-api-second-1. Ticket URL: test-ticket-url-axe-api-second-1. Ticket test-ticket-param-axe-api-second-1: test-ticket-param-val-axe-api-second-1.",
              "ticket": "test-ticket-axe-api-second-1",
              "ticket_system_name": "test-ticket-system-name-axe-api-second-1",
              "ticket_url": "test-ticket-url-axe-api-second-1",
              "ticket_data": {
                "test-ticket-param-axe-api-second-1": "test-ticket-param-val-axe-api-second-1"
              },
              "ticket_comment": "test-comment-axe-api-second-1"
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "root",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-api-second-1. Ticket URL: test-ticket-url-axe-api-second-1. Ticket test-ticket-param-axe-api-second-1: test-ticket-param-val-axe-api-second-1.",
                "ticket": "test-ticket-axe-api-second-1",
                "ticket_system_name": "test-ticket-system-name-axe-api-second-1",
                "ticket_url": "test-ticket-url-axe-api-second-1",
                "ticket_data": {
                  "test-ticket-param-axe-api-second-1": "test-ticket-param-val-axe-api-second-1"
                },
                "ticket_comment": "test-comment-axe-api-second-1"
              }
            ],
            "component": "test-component-axe-api-second-1",
            "connector": "test-connector-axe-api-second-1",
            "connector_name": "test-connector-name-axe-api-second-1",
            "resource": "test-resource-axe-api-second-1"
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
                "_t": "assocticket",
                "a": "root",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-api-second-1. Ticket URL: test-ticket-url-axe-api-second-1. Ticket test-ticket-param-axe-api-second-1: test-ticket-param-val-axe-api-second-1.",
                "ticket": "test-ticket-axe-api-second-1",
                "ticket_system_name": "test-ticket-system-name-axe-api-second-1",
                "ticket_url": "test-ticket-url-axe-api-second-1",
                "ticket_data": {
                  "test-ticket-param-axe-api-second-1": "test-ticket-param-val-axe-api-second-1"
                },
                "ticket_comment": "test-comment-axe-api-second-1"
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
