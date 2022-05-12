Feature: execute action on trigger
  I need to be able to trigger action on event

  Scenario: given one scenario with webhook and processing 2 alarms should use 2 different payloads in webhook
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-multiple-alarm-webhook-1",
      "enabled": true,
      "priority": 266,
      "triggers": [
        "create"
      ],
      "actions": [
        {
          "entity_patterns": [
            {
              "name": {
                "regex_match": "^test-resource-multiple-alarm-webhook-1-\\d"
              }
            }
          ],
          "type": "webhook",
          "parameters": {
            "forward_author": false,
            "author": "test-scenario-multiple-alarm-webhook-1-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {
                "Content-Type": "application/json"
              },
              "payload": "{\"name\": \"{{ `test-scenario-action-multiple-alarm-webhook-1!!!{{ .Alarm.Value.Output }}!!!{{ .Alarm.Value.Resource }}|||{{ .Alarm.Value.State.Value }}` }}\", \"enabled\":true,\"priority\":{{ `{{ $state := .Alarm.Value.State.Value }}{{ if (eq $state 2) }}142{{ else }}143{{ end }}` }},\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"test-action-scenario-multiple-alarm-webhook-1-alarm\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    """json
    {
      "connector" : "test-connector-multiple-alarm-webhook-1",
      "connector_name" : "test-connector-name-multiple-alarm-webhook-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-multiple-alarm-webhook-1",
      "resource" : "test-resource-multiple-alarm-webhook-1-1",
      "state" : 2,
      "output" : "test-output-multiple-alarm-webhook-1-1"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-multiple-alarm-webhook-1",
      "connector_name" : "test-connector-name-multiple-alarm-webhook-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-multiple-alarm-webhook-1",
      "resource" : "test-resource-multiple-alarm-webhook-1-2",
      "state" : 3,
      "output" : "test-output-multiple-alarm-webhook-1-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.component":"test-component-multiple-alarm-webhook-1"}]}&with_steps=true&sort_key=d
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
              "data": {
                "scenario_name": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-1!!!test-resource-multiple-alarm-webhook-1-1|||2"
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
                "_t": "declareticket",
                "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-1",
                "user_id": ""
              }
            ],
            "connector": "test-connector-multiple-alarm-webhook-1",
            "connector_name": "test-connector-name-multiple-alarm-webhook-1",
            "component": "test-component-multiple-alarm-webhook-1",
            "resource": "test-resource-multiple-alarm-webhook-1-1"
          }
        },
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
              "data": {
                "scenario_name": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-2!!!test-resource-multiple-alarm-webhook-1-2|||3"
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
                "_t": "declareticket",
                "a": "test-scenario-multiple-alarm-webhook-1-action-1-author test-resource-multiple-alarm-webhook-1-2",
                "user_id": ""
              }
            ],
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
    When I do GET /api/v4/scenarios?search=test-scenario-action-multiple-alarm-webhook-1&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-1!!!test-resource-multiple-alarm-webhook-1-1|||2",
          "enabled": true,
          "triggers": [
            "create"
          ]
        },
        {
          "name": "test-scenario-action-multiple-alarm-webhook-1!!!test-output-multiple-alarm-webhook-1-2!!!test-resource-multiple-alarm-webhook-1-2|||3",
          "enabled": true,
          "triggers": [
            "create"
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
      }
    }
    """
