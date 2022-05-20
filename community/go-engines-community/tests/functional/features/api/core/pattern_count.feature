Feature: Count matches
  I need to be able to count matches by patterns
  Only admin should be able to count matches by patterns

  Scenario: given count request should return counts
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "max_matched_items": 4
    }
    """
    Then the response code should be 200
    Then I wait 2s
    When I am noperms
    When I do POST /api/v4/patterns-count:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-pattern-count-1",
                "test-resource-to-pattern-count-2",
                "test-resource-to-pattern-count-3",
                "test-resource-to-pattern-count-4"
              ]
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-pattern-count-1",
                "test-resource-to-pattern-count-2",
                "test-resource-to-pattern-count-3",
                "test-resource-to-pattern-count-4"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "alarm_pattern": {
        "count": 4,
        "over_limit": false
      },
      "entity_pattern": {
        "count": 4,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      }
    }
    """

  Scenario: given count request should return counts with over limit
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "max_matched_items": 3
    }
    """
    Then the response code should be 200
    Then I wait 2s
    When I am noperms
    When I do POST /api/v4/patterns-count:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-pattern-count-1",
                "test-resource-to-pattern-count-2",
                "test-resource-to-pattern-count-3",
                "test-resource-to-pattern-count-4"
              ]
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-pattern-count-1",
                "test-resource-to-pattern-count-2",
                "test-resource-to-pattern-count-3",
                "test-resource-to-pattern-count-4"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "alarm_pattern": {
        "count": 4,
        "over_limit": true
      },
      "entity_pattern": {
        "count": 4,
        "over_limit": true
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      }
    }
    """

  Scenario: given count request with missing fields should return empty response
    When I am admin
    When I do POST /api/v4/patterns-count:
    """json
    {
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "alarm_pattern": {
        "count": 0,
        "over_limit": false
      },
      "entity_pattern": {
        "count": 0,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      }
    }
    """

  Scenario: given count request and no auth user should not allow access
    When I do POST /api/v4/patterns-count
    Then the response code should be 401

  Scenario: given count request with invalid patterns format should return bad request error
    When I am admin
    When I do POST /api/v4/patterns-count:
    """json
    {
      "alarm_pattern": [
        []
      ],
      "entity_pattern": [
        []
      ],
      "pbehavior_pattern": [
        []
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "pbehavior_pattern": "PbehaviorPattern is invalid pbehavior pattern."
      }
    }
    """
