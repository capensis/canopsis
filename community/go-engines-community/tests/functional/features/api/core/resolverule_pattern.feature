Feature: Update and delete corporate pattern should affect resolve rule models
  Scenario: given resolve rule and corporate pattern update and delete actions should update patterns in resolve rule
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """
    {
      "_id": "test-resolve-rule-to-edit-patterns",
      "name": "test-resolve-rule-to-edit-patterns-name",
      "description": "test-resolve-rule-to-edit-patterns-description",
      "corporate_entity_pattern": "test-pattern-to-rule-resolve-corporate-update-1",
      "corporate_alarm_pattern": "test-pattern-to-rule-resolve-corporate-update-2",
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-resolve-rule-to-edit-patterns-name",
      "description": "test-resolve-rule-to-edit-patterns-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-resolve-corporate-update-2",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-resolve-corporate-update-2-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-resolve-corporate-update-2-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-resolve-corporate-update-1",
      "corporate_entity_pattern_title": "test-pattern-to-rule-resolve-corporate-update-1-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-resolve-corporate-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-rule-resolve-corporate-update-1:
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
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-resolve-rule-to-edit-patterns-name",
      "description": "test-resolve-rule-to-edit-patterns-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-resolve-corporate-update-2",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-resolve-corporate-update-2-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-resolve-corporate-update-2-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-resolve-corporate-update-1",
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
      ]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-rule-resolve-corporate-update-2:
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
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-resolve-rule-to-edit-patterns-name",
      "description": "test-resolve-rule-to-edit-patterns-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-resolve-corporate-update-2",
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
      "corporate_entity_pattern": "test-pattern-to-rule-resolve-corporate-update-1",
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
      ]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-resolve-corporate-update-1
    Then the response code should be 204
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-resolve-rule-to-edit-patterns-name",
      "description": "test-resolve-rule-to-edit-patterns-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-resolve-corporate-update-2",
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
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-resolve-corporate-update-2
    Then the response code should be 204
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-edit-patterns
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-resolve-rule-to-edit-patterns-name",
      "description": "test-resolve-rule-to-edit-patterns-description",
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
      ]
    }
    """
    Then the response key "corporate_alarm_pattern" should not exist
    Then the response key "corporate_alarm_pattern_title" should not exist
