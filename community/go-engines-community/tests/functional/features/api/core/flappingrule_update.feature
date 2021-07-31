Feature: Update an flapping rule
  I need to be able to update an flapping rule
  Only admin should be able to update an flapping rule

  Scenario: given update request should update flapping rule
    When I am admin
    Then I do PUT /api/v4/flapping-rules/test-flapping-rule-to-update:
    """json
    {
      "_id": "whatever",
      "description": "updated flapping rule",
      "flapping_interval": {
        "seconds": 200,
        "unit": "m"
      },
      "flapping_freq_limit": 15,
      "alarm_patterns": [
        {
          "v": {
            "connector": "test-flapping-rule-to-update-pattern-updated"
          }
        }
      ],
      "priority": 7
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-flapping-rule-to-update",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "updated flapping rule",
      "flapping_interval": {
        "seconds": 200,
        "unit": "m"
      },
      "flapping_freq_limit": 15,
      "alarm_patterns": [
        {
          "v": {
            "connector": "test-flapping-rule-to-update-pattern-updated"
          }
        }
      ],
      "priority": 7
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/flapping-rules/test-flapping-rule-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/flapping-rules/test-flapping-rule-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/flapping-rules/test-flapping-rule-not-found:
    """json
    {
      "description": "updated flapping rule",
      "flapping_interval": {
        "seconds": 50,
        "unit": "s"
      },
      "flapping_freq_limit": 13,
      "alarm_patterns": [
        {
          "v": {
            "connector": "test-flapping-rule-to-update-pattern-updated"
          }
        }
      ],
      "priority": 1
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
    Then I do PUT /api/v4/flapping-rules/test-flapping-rule-to-update:
    """
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "flapping_freq_limit": "FlappingFreqLimit is missing.",
        "flapping_interval.seconds": "Seconds is missing.",
        "flapping_interval.unit": "Unit is missing.",
        "priority": "Priority is missing."
      }
    }
    """
