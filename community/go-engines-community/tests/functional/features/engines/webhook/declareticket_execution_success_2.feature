Feature: run a declare ticket rule
  I need to be able to run a declare ticket rule
  Only admin should be able to run a declare ticket rule

  @concurrent
  Scenario: given declare ticket rule and ticket resources should update resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-second-1",
        "connector_name" : "test-connector-name-declareticket-execution-second-1",
        "component" :  "test-component-declareticket-execution-second-1",
        "source_type" : "component"
      },
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-second-1",
        "connector_name" : "test-connector-name-declareticket-execution-second-1",
        "component" :  "test-component-declareticket-execution-second-1",
        "resource" : "test-resource-declareticket-execution-second-1",
        "source_type" : "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response componentAlarmId={{ (index .lastResponse.data 0)._id }}
    When I save response resourceAlarmId={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-second-1-name",
      "system_name": "test-declareticketrule-declareticket-execution-second-1-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-second-1\",\"url\":\"https://test-host\",\"name\":\"{{ `{{ range .Alarms }}{{ .Entity.Name }}{{ end }}`}}\"}",
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
                "test-component-declareticket-execution-second-1"
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
          "{{ .componentAlarmId }}"
        ],
        "comment": "test-comment-declareticket-execution-second-1",
        "ticket_resources": true
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
            "{{ .componentAlarmId }}"
          ],
          "comment": "test-comment-declareticket-execution-second-1",
          "ticket_resources": true
        }
      }
    ]
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "declareticketwebhook",
      "connector" : "test-connector-declareticket-execution-second-1",
      "connector_name" : "test-connector-name-declareticket-execution-second-1",
      "component" :  "test-component-declareticket-execution-second-1",
      "resource" : "test-resource-declareticket-execution-second-1",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-1&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
              "ticket": "test-ticket-declareticket-execution-second-1",
              "ticket_url": "https://test-host",
              "ticket_data": {
                "name": "test-component-declareticket-execution-second-1"
              },
              "ticket_comment": "test-comment-declareticket-execution-second-1",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
                "ticket": "test-ticket-declareticket-execution-second-1",
                "ticket_url": "https://test-host",
                "ticket_data": {
                  "name": "test-component-declareticket-execution-second-1"
                },
                "ticket_comment": "test-comment-declareticket-execution-second-1",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-1",
            "connector_name" : "test-connector-name-declareticket-execution-second-1",
            "component" :  "test-component-declareticket-execution-second-1"
          }
        },
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
              "ticket": "test-ticket-declareticket-execution-second-1",
              "ticket_url": "https://test-host",
              "ticket_data": {
                "name": "test-component-declareticket-execution-second-1"
              },
              "ticket_comment": "test-comment-declareticket-execution-second-1",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
                "ticket": "test-ticket-declareticket-execution-second-1",
                "ticket_url": "https://test-host",
                "ticket_data": {
                  "name": "test-component-declareticket-execution-second-1"
                },
                "ticket_comment": "test-comment-declareticket-execution-second-1",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-1",
            "connector_name" : "test-connector-name-declareticket-execution-second-1",
            "component" :  "test-component-declareticket-execution-second-1",
            "resource" : "test-resource-declareticket-execution-second-1"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .componentAlarmId }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .resourceAlarmId }}",
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1."
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1."
      }
    ]
    """
    Then the response array key "1.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1."
      }
    ]
    """
