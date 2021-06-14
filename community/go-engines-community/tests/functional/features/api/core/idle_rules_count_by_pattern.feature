Feature: check matched entities count for idle rules
  I need to be able to create a PBehavior
  Only admin should be able to create a PBehavior

  Scenario: POST a valid count request but unauthorized
    When I do POST /api/v4/idle-rules/count
    Then the response code should be 401

  Scenario: POST a valid count request but without permissions
    When I am noperms
    When I do POST /api/v4/idle-rules/count
    Then the response code should be 403

  Scenario: POST a valid count request where limit by entities is not reached
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [
        {
          "name": "test-idle-rule-count-patterns-resource-1"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-2"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-3"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-4"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": false,
      "total_count_entities": 4,
      "total_count_alarms": 0
    }
    """

  Scenario: POST a valid count request where limit by entities is reached
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-idle-rule-count-patterns-resource"
          }
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": true,
      "total_count_entities": 5,
      "total_count_alarms": 0
    }
    """

  Scenario: POST a valid count request where limit by alarm is not reached
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "alarm_patterns": [
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-1"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-2"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-3"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-4"
          }
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": false,
      "total_count_entities": 0,
      "total_count_alarms": 4
    }
    """

  Scenario: POST a valid count request where limit by alarm is reached
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "alarm_patterns": [
        {
          "v": {
            "resource": {
              "regex_match": "test-idle-rule-count-patterns-resource"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": true,
      "total_count_entities": 0,
      "total_count_alarms": 5
    }
    """

  Scenario: POST a valid count request where limit by alarm is not reached with both patterns
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [
        {
          "name": "test-idle-rule-count-patterns-resource-1"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-2"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-3"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-4"
        }
      ],
      "alarm_patterns": [
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-1"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-2"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-3"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-4"
          }
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": false,
      "total_count_entities": 4,
      "total_count_alarms": 4
    }
    """

  Scenario: POST a valid count request where limit by alarm is reached regardless of patterns
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-idle-rule-count-patterns-resource"
          }
        }
      ],
      "alarm_patterns": [
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-1"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-2"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-3"
          }
        },
        {
          "v": {
            "resource": "test-idle-rule-count-patterns-resource-4"
          }
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": true
    }
    """
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [
        {
          "name": "test-idle-rule-count-patterns-resource-1"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-2"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-3"
        },
        {
          "name": "test-idle-rule-count-patterns-resource-4"
        }
      ],
      "alarm_patterns": [
        {
          "v": {
            "resource": {
              "regex_match": "test-idle-rule-count-patterns-resource"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": true
    }
    """

  Scenario: POST an invalid count request with empty patterns
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [{}],
      "alarm_patterns": [{}]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
        "errors": {
            "alarm_patterns": "alarm pattern list contains an empty pattern.",
            "entity_patterns": "entity pattern list contains an empty pattern."
        }
    }
    """

  Scenario: POST an invalid count request with empty patterns
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [{}],
      "alarm_patterns": [{}]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "alarm_patterns": "alarm pattern list contains an empty pattern.",
        "entity_patterns": "entity pattern list contains an empty pattern."
      }
    }
    """

  Scenario: POST an invalid count request with invalid patterns
    When I am admin
    When I do POST /api/v4/idle-rules/count:
    """
    {
      "entity_patterns": [
        {
          "test": "test"
        }
      ],
      "alarm_patterns": [
        {
          "test": "test"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "alarm_patterns": "Invalid alarm pattern list.",
        "entity_patterns": "Invalid entity pattern list."
      }
    }
    """