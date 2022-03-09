Feature: Update a scenario
  I need to be able to update a scenario
  Only admin should be able to update a scenario

  Scenario: given update request should update exception
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-1:
    """
    {
      "name": "test-scenario-to-update-1-name",
      "enabled": true,
      "priority": 15,
      "triggers": ["create","pbhenter"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-update-1-alarm-updated"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze updated",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-scenario-to-update-1",
      "name": "test-scenario-to-update-1-name",
      "author": "root",
      "enabled": true,
      "priority": 15,
      "delay": null,
      "disable_during_periods": null,
      "triggers": ["create","pbhenter"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-update-1-alarm-updated"
            }
          ],
          "entity_patterns": null,
          "type": "snooze",
          "parameters": {
            "author": "root",
            "user": "root",
            "output": "test snooze updated",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "test comment"
        }
      ],
      "created": 1605263992
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-1:
    """
    {
      "name": "test-scenario-to-check-unique-name-name",
      "enabled": true,
      "priority": 16,
      "triggers": ["create","pbhenter"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-update-1-alarm-updated"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze updated",
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
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  Scenario: given update request with already exists priority should return error
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-1:
    """
    {
      "name": "test-scenario-to-update-1-name",
      "enabled": true,
      "priority": 2,
      "triggers": ["create","pbhenter"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-update-1-alarm-updated"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze updated",
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
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "priority": "Priority already exists."
      }
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/scenarios/notexist
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/scenarios/notexist
    Then the response code should be 403

  Scenario: given no exist scenario id should return error
    When I am admin
    When I do PUT /api/v4/scenarios/notexist:
    """
    {
      "name": "test-scenario-to-update-notexist-name",
      "enabled": true,
      "priority": 20,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-update-notexist-alarm"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
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
    Then the response code should be 404

  Scenario: given create request with custom_id shouldn't update id
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-check-id:
    """
    {
      "_id": "change-id",
      "name": "my_scenario-name-new",
      "enabled": true,
      "priority": 10999,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test-scenario-to-create-1-action-1-output",
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
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-scenario-to-check-id",
      "name": "my_scenario-name-new"
    }
    """
