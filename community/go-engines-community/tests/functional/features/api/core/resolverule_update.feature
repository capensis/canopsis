Feature: Update an resolve rule
  I need to be able to update an resolve rule
  Only admin should be able to update an resolve rule

  Scenario: given update request should update resolve rule
    When I am admin
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-1:
    """json
    {
      "name": "test-resolve-rule-to-update-1-name-updated",
      "description": "test-resolve-rule-to-update-1-description-updated",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-update-1-pattern-updated"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-update-1-resource-updated"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "name": "test-resolve-rule-to-update-1-name-updated",
      "description": "test-resolve-rule-to-update-1-description-updated",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-update-1-pattern-updated"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-update-1-resource-updated"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-not-found:
    """json
    {
      "name": "test-resolve-rule-to-update-2-name-updated",
      "description": "test-resolve-rule-to-update-2-description-updated",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-to-update-2-pattern-updated"
          }
        }
      ],
      "entity_patterns": [
        {
          "name": "test-resolve-rule-to-update-2-resource-updated"
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update:
    """
    {}
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

  Scenario: given create request with already exists id and name should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-1:
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
