Feature: execute action on trigger
  I need to be able to trigger action on event

  @concurrent
  Scenario: given one scenario with webhook and processing 2 alarms should use 2 different payloads in webhook
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-multiple-alarm-webhook-1-name",
      "priority": 10000,
      "enabled": true,
      "triggers": [
        "create"
      ],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "regexp",
                  "value": "^test-resource-multiple-alarm-webhook-1-\\d"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "forward_author": false,
            "author": "test-scenario-multiple-alarm-webhook-1-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"{{ `test-scenario-action-multiple-alarm-webhook-1!!!{{ .Alarm.Value.Output }}!!!{{ .Alarm.Value.Resource }}|||{{ .Alarm.Value.State.Value }}` }}\"}"
            },
            "declare_ticket": {
              "ticket_id": "_id"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response scenarioId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-multiple-alarm-webhook-1",
        "connector_name": "test-connector-name-multiple-alarm-webhook-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-multiple-alarm-webhook-1",
        "resource": "test-resource-multiple-alarm-webhook-1-1",
        "state": 2,
        "output": "test-output-multiple-alarm-webhook-1-1"
      },
      {
        "connector": "test-connector-multiple-alarm-webhook-1",
        "connector_name": "test-connector-name-multiple-alarm-webhook-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-multiple-alarm-webhook-1",
        "resource": "test-resource-multiple-alarm-webhook-1-2",
        "state": 3,
        "output": "test-output-multiple-alarm-webhook-1-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-component-multiple-alarm-webhook-1&sort_by=d&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
                "ticket": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-1!!!test-resource-multiple-alarm-webhook-1-1|||2",
                "ticket_rule_id": "{{ .scenarioId }}",
                "ticket_rule_name": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
              "ticket": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-1!!!test-resource-multiple-alarm-webhook-1-1|||2",
              "ticket_rule_id": "{{ .scenarioId }}",
              "ticket_rule_name": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
            },
            "connector": "test-connector-multiple-alarm-webhook-1",
            "connector_name": "test-connector-name-multiple-alarm-webhook-1",
            "component": "test-component-multiple-alarm-webhook-1",
            "resource": "test-resource-multiple-alarm-webhook-1-1"
          }
        },
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
                "ticket": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-2!!!test-resource-multiple-alarm-webhook-1-2|||3",
                "ticket_rule_id": "{{ .scenarioId }}",
                "ticket_rule_name": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
              "ticket": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-2!!!test-resource-multiple-alarm-webhook-1-2|||3",
              "ticket_rule_id": "{{ .scenarioId }}",
              "ticket_rule_name": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
            },
            "connector": "test-connector-multiple-alarm-webhook-1",
            "connector_name": "test-connector-name-multiple-alarm-webhook-1",
            "component": "test-component-multiple-alarm-webhook-1",
            "resource": "test-resource-multiple-alarm-webhook-1-2"
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
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
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
        "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
        "user_id": "",
        "m": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
        "user_id": "",
        "m": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
        "user_id": "",
        "ticket": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-1!!!test-resource-multiple-alarm-webhook-1-1|||2",
        "ticket_rule_id": "{{ .scenarioId }}",
        "ticket_rule_name": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
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
        "_t": "webhookstart",
        "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
        "user_id": "",
        "m": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
        "user_id": ""
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
        "user_id": "",
        "ticket": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-2!!!test-resource-multiple-alarm-webhook-1-2|||3",
        "ticket_rule_id": "{{ .scenarioId }}",
        "ticket_rule_name": "Scenario: test-scenario-multiple-alarm-webhook-1-name"
      }
    ]
    """
