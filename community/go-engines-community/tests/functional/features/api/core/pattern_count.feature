Feature: Count matches
  I need to be able to count matches by patterns
  Only admin should be able to count matches by patterns

  Scenario: given alarms count request should return counts
    When I am noperms
    When I do POST /api/v4/patterns-alarms-count:
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
                "test-resource-to-pattern-count-5"
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
        "count": 2,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      },
      "all": {
        "count": 2,
        "over_limit": false
      },
      "entities": {
        "count": 3,
        "over_limit": false
      }
    }
    """
    When I do POST /api/v4/patterns-alarms-count:
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
      "entity_pattern": []
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
        "count": 0,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      },
      "all": {
        "count": 4,
        "over_limit": false
      },
      "entities": {
        "count": 0,
        "over_limit": false
      }
    }
    """
    When I do POST /api/v4/patterns-alarms-count:
    """json
    {
      "alarm_pattern": [],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-pattern-count-1",
                "test-resource-to-pattern-count-2",
                "test-resource-to-pattern-count-5"
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
        "count": 0,
        "over_limit": false
      },
      "entity_pattern": {
        "count": 2,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      },
      "all": {
        "count": 2,
        "over_limit": false
      },
      "entities": {
        "count": 3,
        "over_limit": false
      }
    }
    """

  Scenario: given entities count request should return counts
    When I am noperms
    When I do POST /api/v4/patterns-entities-count:
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
                "test-resource-to-pattern-count-5"
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
        "count": 3,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      },
      "all": {
        "count": 2,
        "over_limit": false
      }
    }
    """
    When I do POST /api/v4/patterns-entities-count:
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
      "entity_pattern": []
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
        "count": 0,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      },
      "all": {
        "count": 4,
        "over_limit": false
      }
    }
    """
    When I do POST /api/v4/patterns-entities-count:
    """json
    {
      "alarm_pattern": [],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-pattern-count-1",
                "test-resource-to-pattern-count-2",
                "test-resource-to-pattern-count-5"
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
        "count": 0,
        "over_limit": false
      },
      "entity_pattern": {
        "count": 3,
        "over_limit": false
      },
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      },
      "all": {
        "count": 3,
        "over_limit": false
      }
    }
    """

  @standalone
  Scenario: given count request should return counts with over limit
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "max_matched_items": 3
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I am noperms
    When I do POST /api/v4/patterns-alarms-count:
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
      },
      "all": {
        "count": 4,
        "over_limit": true
      },
      "entities": {
        "count": 4,
        "over_limit": true
      }
    }
    """
    When I do POST /api/v4/patterns-entities-count:
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
      },
      "all": {
        "count": 4,
        "over_limit": true
      }
    }
    """

  Scenario: given alarms count request with missing fields should return empty response
    When I am admin
    When I do POST /api/v4/patterns-alarms-count:
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
      },
      "all": {
        "count": 0,
        "over_limit": false
      },
      "entities": {
        "count": 0,
        "over_limit": false
      }
    }
    """

  Scenario: given alarms count request and no auth user should not allow access
    When I do POST /api/v4/patterns-alarms-count
    Then the response code should be 401

  Scenario: given alarms count request with invalid patterns format should return bad request error
    When I am admin
    When I do POST /api/v4/patterns-alarms-count:
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

  Scenario: given entities count request with missing fields should return empty response
    When I am admin
    When I do POST /api/v4/patterns-entities-count:
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
      },
      "all": {
        "count": 0,
        "over_limit": false
      }
    }
    """

  Scenario: given entities count request and no auth user should not allow access
    When I do POST /api/v4/patterns-entities-count
    Then the response code should be 401

  Scenario: given entities count request with invalid patterns format should return bad request error
    When I am admin
    When I do POST /api/v4/patterns-entities-count:
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
