Feature: Update and delete corporate pattern should affect idle rule models
  Scenario: given idle rule and corporate pattern update and delete actions should update patterns in idle rule
    When I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "_id": "test-idle-rule-to-edit-patterns",
      "name": "test-idle-rule-to-edit-patterns-name",
      "description": "test-idle-rule-to-edit-patterns-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-idle-corporate-update-1",
      "corporate_entity_pattern": "test-pattern-to-rule-idle-corporate-update-2",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-edit-patterns-operation-output",
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
    When I do GET /api/v4/idle-rules/test-idle-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-idle-rule-to-edit-patterns-name",
      "description": "test-idle-rule-to-edit-patterns-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-idle-corporate-update-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-idle-corporate-update-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-idle-corporate-update-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-idle-corporate-update-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-idle-corporate-update-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-idle-corporate-update-2-pattern"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-edit-patterns-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-rule-idle-corporate-update-2:
    """
    {
      "title": "new entity pattern title",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new entity pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/idle-rules/test-idle-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-idle-rule-to-edit-patterns-name",
      "description": "test-idle-rule-to-edit-patterns-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-idle-corporate-update-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-idle-corporate-update-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-idle-corporate-update-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-idle-corporate-update-2",
      "corporate_entity_pattern_title": "new entity pattern title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new entity pattern"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-edit-patterns-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-rule-idle-corporate-update-1:
    """
    {
      "title": "new alarm pattern title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "new alarm pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/idle-rules/test-idle-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-idle-rule-to-edit-patterns-name",
      "description": "test-idle-rule-to-edit-patterns-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-idle-corporate-update-1",
      "corporate_alarm_pattern_title": "new alarm pattern title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "new alarm pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-idle-corporate-update-2",
      "corporate_entity_pattern_title": "new entity pattern title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new entity pattern"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-edit-patterns-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-idle-corporate-update-1
    Then the response code should be 204
    When I do GET /api/v4/idle-rules/test-idle-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-idle-rule-to-edit-patterns-name",
      "description": "test-idle-rule-to-edit-patterns-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "new alarm pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-idle-corporate-update-2",
      "corporate_entity_pattern_title": "new entity pattern title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new entity pattern"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-edit-patterns-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    Then the response key "corporate_alarm_pattern" should not exist
    Then the response key "corporate_alarm_pattern_title" should not exist
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-idle-corporate-update-2
    Then the response code should be 204
    When I do GET /api/v4/idle-rules/test-idle-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-idle-rule-to-edit-patterns-name",
      "description": "test-idle-rule-to-edit-patterns-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "new alarm pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new entity pattern"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-edit-patterns-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
