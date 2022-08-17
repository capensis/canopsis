Feature: Create a idle rule
  I need to be able to create a idle rule
  Only admin should be able to create a idle rule

  Scenario: given create alarm rule request should return ok
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-1-name",
      "description": "test-idle-rule-to-create-1-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-create-1-alarm"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-create-1-resource"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-1-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-1-name",
      "description": "test-idle-rule-to-create-1-description",
      "author": "root",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-create-1-alarm"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-create-1-resource"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-1-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given create entity rule request should return ok
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-2-name",
      "description": "test-idle-rule-to-create-2-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-create-2-resource"
        }
      ],
      "disable_during_periods": ["pause"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-2-name",
      "description": "test-idle-rule-to-create-2-description",
      "author": "root",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-create-2-resource"
        }
      ],
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-3-name",
      "description": "test-idle-rule-to-create-3-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 22,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-create-3-alarm"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-create-3-resource"
        }
      ],
      "operation": {
        "type": "pbehavior",
        "parameters": {
          "name": "test-idle-rule-to-create-3-operation-name",
          "rrule": "FREQ=DAILY",
          "reason": "test-reason-to-edit-idle-rule",
          "type": "test-type-to-edit-idle-rule",
          "start_on_trigger": true,
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/idle-rules/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-3-name",
      "description": "test-idle-rule-to-create-3-description",
      "author": "root",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 22,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-create-3-alarm"
        }
      ],
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-create-3-resource"
        }
      ],
      "operation": {
        "type": "pbehavior",
        "parameters": {
          "name": "test-idle-rule-to-create-3-operation-name",
          "rrule": "FREQ=DAILY",
          "reason": {
            "_id": "test-reason-to-edit-idle-rule",
            "name": "test-reason-to-edit-idle-rule-name",
            "description": "test-reason-to-edit-idle-rule-description",
            "created": 1592215337
          },
          "type": {
            "_id": "test-type-to-edit-idle-rule",
            "name": "test-type-to-edit-idle-rule-name",
            "description": "test-type-to-edit-idle-rule-description",
            "icon_name": "test-type-to-edit-idle-rule-icon"
          },
          "start_on_trigger": true,
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/idle-rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/idle-rules
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "duration.value": "Value is missing.",
        "duration.unit": "Unit is missing.",
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "priority": "Priority is missing.",
        "type": "Type is missing."
      }
    }
    """

  Scenario: given invalid create request with invalid type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "notexists"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [alarm entity]."
      }
    }
    """

  Scenario: given invalid create request with alarm type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "alarm",
      "alarm_patterns": [],
      "entity_patterns": [],
      "operation": {
        "type": "notexists"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_condition": "AlarmCondition is missing.",
        "alarm_patterns": "AlarmPatterns or EntityPatterns is required.",
        "entity_patterns": "EntityPatterns or AlarmPatterns is required.",
        "operation.type": "Type must be one of [ack ackremove cancel assocticket changestate snooze pbehavior]."
      }
    }
    """

  Scenario: given invalid create request with entity type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "entity",
      "alarm_patterns": [{"_id": "notexists"}],
      "entity_patterns": [],
      "operation": {
        "type": "notexists"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_patterns": "AlarmPatterns is not empty.",
        "entity_patterns": "EntityPatterns is missing.",
        "operation": "Operation is not empty."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """
