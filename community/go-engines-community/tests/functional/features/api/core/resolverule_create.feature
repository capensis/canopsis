Feature: Create an resolve rule
  I need to be able to create a resolve rule
  Only admin should be able to create a resolve rule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-1-name",
      "description": "test-resolve-rule-to-create-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-create-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
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
      "name": "test-resolve-rule-to-create-1-name",
      "description": "test-resolve-rule-to-create-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-create-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    When I do GET /api/v4/resolve-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-resolve-rule-to-create-1-name",
      "description": "test-resolve-rule-to-create-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-create-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """

  Scenario: given create request should update priority of next rules
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-2-priority-1-name",
      "description": "test-resolve-rule-to-create-2-priority-1-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-create-2-priority-1-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-create-2-priority-1-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-2-priority-2-name",
      "description": "test-resolve-rule-to-create-2-priority-2-description",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-create-2-priority-2-pattern"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-create-2-priority-2-resource"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/resolve-rules?search=test-resolve-rule-to-create-2&sort_by=priority
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resolve-rule-to-create-2-priority-2-name"
        },
        {
          "name": "test-resolve-rule-to-create-2-priority-1-name"
        }
      ]
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "duration.value": "Value is missing.",
        "duration.unit": "Unit is missing.",
        "priority": "Priority is missing."
      }
    }
    """

  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/resolve-rules:
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
        "alarm_patterns": "Invalid alarm pattern list.",
        "entity_patterns": "Invalid entity pattern list."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/resolve-rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/resolve-rules
    Then the response code should be 403

  Scenario: given create request with already exists id and name should return error
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-to-check-unique"
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
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-check-unique-name"
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
