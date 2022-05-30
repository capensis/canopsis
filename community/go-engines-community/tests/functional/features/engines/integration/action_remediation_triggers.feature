Feature: instruction triggers should trigger an action
  I need to be able to trigger an action by instruction triggers

  Scenario: given scenario should be triggered by instructionfail trigger
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-1-name",
      "enabled": true,
      "priority": 20,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "v": {
                "component": "test-component-action-1"
              }
            }
          ],
          "entity_patterns": [
            {
              "type": "resource"
            }
          ],
          "type": "assocticket",
          "parameters": {
            "forward_author": false,
            "author": "test-scenario-action-1-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "output": "test-scenario-action-1-action-1-output {{ `{{ .Entity.Name }} {{ .Alarm.Value.State.Value }}` }}",
            "ticket": "test-scenario-action-1-action-1-ticket"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
