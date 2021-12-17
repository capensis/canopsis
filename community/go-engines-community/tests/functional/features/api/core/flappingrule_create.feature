Feature: Create an flapping rule
  I need to be able to create a flapping rule
  Only admin should be able to create a flapping rule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "name": "test-flapping-rule-to-create-1-name",
      "description": "test-flapping-rule-to-create-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-flapping-rule-to-create-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "freq_limit": 3,
      "priority": 5
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-flapping-rule-to-create-1-name",
      "description": "test-flapping-rule-to-create-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-flapping-rule-to-create-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "freq_limit": 3,
      "priority": 5
    }
    """
    When I do GET /api/v4/flapping-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-flapping-rule-to-create-1-name",
      "description": "test-flapping-rule-to-create-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-flapping-rule-to-create-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "freq_limit": 3,
      "priority": 5
    }
    """

  Scenario: given create request should update priority of next rules
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "name": "test-flapping-rule-to-create-2-priority-1-name",
      "description": "test-flapping-rule-to-create-2-priority-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-2-priority-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-flapping-rule-to-create-2-priority-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "freq_limit": 3,
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "name": "test-flapping-rule-to-create-2-priority-2-name",
      "description": "test-flapping-rule-to-create-2-priority-2-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-2-priority-2-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-flapping-rule-to-create-2-priority-2-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "freq_limit": 3,
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/flapping-rules?search=test-flapping-rule-to-create-2&sort_by=priority&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-flapping-rule-to-create-2-priority-2-name"
        },
        {
          "name": "test-flapping-rule-to-create-2-priority-1-name"
        }
      ]
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarm_patterns": "AlarmPatterns or EntityPatterns is required.",
        "entity_patterns": "EntityPatterns or AlarmPatterns is required.",
        "name": "Name is missing.",
        "freq_limit": "FreqLimit is missing.",
        "duration.value": "Value is missing.",
        "duration.unit": "Unit is missing.",
        "priority": "Priority is missing."
      }
    }
    """

  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "alarm_patterns": [
        {
          "v": {
            "component_name": "ram"
          }
        }
      ],
      "entity_patterns": [
        {
          "component_name": "ram"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_patterns":"Invalid alarm pattern list.",
        "entity_patterns":"Invalid entity pattern list."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/flapping-rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/flapping-rules
    Then the response code should be 403

  Scenario: given create request with already exists id and name should return error
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "_id": "test-flapping-rule-to-check-unique"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "name": "test-flapping-rule-to-check-unique-name"
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
