Feature: send activation event on unsnooze
  I need to be able to trigger rule on alarm activation

  Scenario: given event for new alarm and snooze action should send event on unsnooze
    Given I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-axe-action-activation-name",
      "priority": 10058,
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
                  "value": "test-resource-axe-action-activation-event"
                }
              }
            ]
          ],
          "type":"snooze",
          "parameters": {
            "duration": {
              "value": 3,
              "unit": "s"
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
    """
    {
      "connector" : "test-connector-axe-action-activation-event",
      "connector_name" : "test-connector-name-axe-action-activation-event",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-action-activation-event",
      "resource" : "test-resource-axe-action-activation-event",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-action-activation-event
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-axe-action-activation-event",
            "connector_name" : "test-connector-name-axe-action-activation-event",
            "component" : "test-component-axe-action-activation-event",
            "resource" : "test-resource-axe-action-activation-event"
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
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 0,10
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
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "snooze"}
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
