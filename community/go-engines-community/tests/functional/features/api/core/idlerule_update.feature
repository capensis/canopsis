Feature: Update a idle rule
  I need to be able to update a idle rule
  Only admin should be able to update a idle rule

  Scenario: given update request should update exception
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update:
    """
    {
      "name": "test-idle-rule-to-update-name",
      "description": "test-idle-rule-to-update-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-update-alarm-updated"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-update-resource-updated"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-update-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-idle-rule-to-update",
      "name": "test-idle-rule-to-update-name",
      "description": "test-idle-rule-to-update-description",
      "author": "root",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "created": 1616567033,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-update-alarm-updated"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-update-resource-updated"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "author": "root",
          "output": "test-idle-rule-to-update-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update:
    """
    {
      "name": "test-idle-rule-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """
    
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/idle-rules/notexist
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/idle-rules/notexist
    Then the response code should be 403

  Scenario: given no exist idle rule id should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/notexist:
    """
    {
      "name": "test-idle-rule-notexists-name",
      "description": "test-idle-rule-notexists-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 31,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-notexists-alarm-updated"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-notexists-resource-updated"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-notexists-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 404
