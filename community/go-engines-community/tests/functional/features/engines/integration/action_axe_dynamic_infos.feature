Feature: update dynamic infos when alarm is updated by action
  I need to be able to update dynamic infos when action is applied to alarm.

  @concurrent
  Scenario: given dynamic infos and scenarios should update alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-action-axe-dynamic-infos-1-name",
      "description": "test-dynamicinfos-action-axe-dynamic-infos-1-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-action-axe-dynamic-infos-1"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.ack.initiator",
            "cond": {
              "type": "eq",
              "value": "system"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "info1",
          "value": "test-dynamicinfo-action-axe-dynamic-infos-1"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response dynamicInfoRuleId={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-axe-dynamic-infos-1-1-name",
      "priority": 10084,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-axe-dynamic-infos-1"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-output-action-axe-dynamic-infos-1"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-axe-dynamic-infos-1-2-name",
      "priority": 10085,
      "enabled": true,
      "triggers": ["ack"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-axe-dynamic-infos-1"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-axe-dynamic-infos-1\",\"name\":\"{{ `{{ index (index .Alarm.Value.Infos \"` }}{{ .dynamicInfoRuleId }}{{ `\") \"info1\" }}`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "name": "name"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-action-axe-dynamic-infos-1",
      "connector": "test-connector-action-axe-dynamic-infos-1",
      "connector_name": "test-connector-name-action-axe-dynamic-infos-1",
      "component":  "test-component-action-axe-dynamic-infos-1",
      "resource": "test-resource-action-axe-dynamic-infos-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-axe-dynamic-infos-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-action-axe-dynamic-infos-1",
            "infos": {
              "{{ .dynamicInfoRuleId }}": {
                "info1": "test-dynamicinfo-action-axe-dynamic-infos-1"
              }
            },
            "ticket": {
              "ticket": "test-ticket-action-axe-dynamic-infos-1",
              "ticket_data": {
                "name": "test-dynamicinfo-action-axe-dynamic-infos-1"
              }
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
