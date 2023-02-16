Feature: Update and delete corporate pattern should affect eventfilter models
  Scenario: given eventfilter and corporate pattern update and delete actions should update patterns in eventfiler
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "_id": "eventfilter-pattern-1",
      "description": "test create 4",
      "type": "enrichment",
      "corporate_entity_pattern": "test-pattern-to-rule-eventfilter-corporate-update-1",
      "priority": 0,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "external_data": {
        "test": {
          "type": "mongo"
        }
      }
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/eventfilter/rules/eventfilter-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 4",
      "type": "enrichment",
      "corporate_entity_pattern": "test-pattern-to-rule-eventfilter-corporate-update-1",
      "corporate_entity_pattern_title": "test-pattern-to-rule-eventfilter-corporate-update-1-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-eventfilter-corporate-update-1-pattern"
            }
          }
        ]
      ],
      "priority": 0,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "external_data": {
        "test": {
          "type": "mongo"
        }
      }
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-rule-eventfilter-corporate-update-1:
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
    When I do GET /api/v4/eventfilter/rules/eventfilter-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 4",
      "type": "enrichment",
      "corporate_entity_pattern": "test-pattern-to-rule-eventfilter-corporate-update-1",
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
      "priority": 0,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "external_data": {
        "test": {
          "type": "mongo"
        }
      }
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-eventfilter-corporate-update-1
    Then the response code should be 204
    When I do GET /api/v4/eventfilter/rules/eventfilter-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 4",
      "type": "enrichment",
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
      "priority": 0,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "external_data": {
        "test": {
          "type": "mongo"
        }
      }
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
