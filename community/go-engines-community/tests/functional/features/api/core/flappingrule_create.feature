Feature: Create an flapping rule
  I need to be able to create a flapping rule
  Only admin should be able to create a flapping rule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "_id": "test-flapping-rule-to-create-1",
      "description": "test create 1",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-1-pattern"
          }
        }
      ],
      "flapping_interval": {
        "seconds": 10,
        "unit": "s"
      },
      "flapping_freq_limit": 3,
      "priority": 5
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-flapping-rule-to-create-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 1",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-1-pattern"
          }
        }
      ],
      "flapping_interval": {
        "seconds": 10,
        "unit": "s"
      },
      "flapping_freq_limit": 3,
      "priority": 5
    }
    """
    When I do GET /api/v4/flapping-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-flapping-rule-to-create-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 1",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-1-pattern"
          }
        }
      ],
      "flapping_interval": {
        "seconds": 10,
        "unit": "s"
      },
      "flapping_freq_limit": 3,
      "priority": 5
    }
    """

  Scenario: given create request should update priority of next rules
    When I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "_id": "test-flapping-rule-to-create-2",
      "description": "test create 2",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-2-pattern"
          }
        }
      ],
      "flapping_interval": {
        "seconds": 10,
        "unit": "s"
      },
      "flapping_freq_limit": 3,
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "_id": "test-flapping-rule-to-create-3",
      "description": "test create 3",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-flapping-rule-to-create-3-pattern"
          }
        }
      ],
      "flapping_interval": {
        "seconds": 10,
        "unit": "s"
      },
      "flapping_freq_limit": 3,
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/flapping-rules/test-flapping-rule-to-create-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "priority": 6
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
        "description": "Description is missing.",
        "flapping_freq_limit": "FlappingFreqLimit is missing.",
        "flapping_interval.seconds": "Seconds is missing.",
        "flapping_interval.unit": "Unit is missing.",
        "priority": "Priority is missing."
      }
    }
    """

  Scenario: given create request with wrong alarm patterns format should return bad request
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
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_patterns":"Invalid alarm patterns."
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

  Scenario: given create request with already exists id should return error
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
