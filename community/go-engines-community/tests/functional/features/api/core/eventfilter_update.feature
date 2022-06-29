Feature: Update an eventfilter
  I need to be able to update an eventfilter
  Only admin should be able to update an eventfilter

  Scenario: given update request with event_pattern should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-1",
      "author": "root",
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true,
      "created": 1608635535
    }
    """

  Scenario: given update request with entity_pattern should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-2:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern"
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
              "value": "test-eventfilter-update-2-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-2",
      "author": "root",
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern"
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
              "value": "test-eventfilter-update-2-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true,
      "created": 1608635535
    }
    """

  Scenario: given update request with both patterns should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-3:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern-updated"
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
              "value": "test-eventfilter-update-3-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-3",
      "author": "root",
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern-updated"
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
              "value": "test-eventfilter-update-3-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true,
      "created": 1608635535
    }
    """

  Scenario: given update request with corporate pattern should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-4:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-4",
      "author": "root",
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-2-pattern"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true,
      "created": 1608635535
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-not-found:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-update-pattern-updated"
            }
          }
        ]
      ],
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
    
  Scenario: given PUT change_entity rule requests should return error, because of empty config
    Given I am admin
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
      "description": "update change_entity",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "never be used change entity update test"
            }
          }
        ]
      ],
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
      "description": "update change_entity",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "never be used change entity update test"
            }
          }
        ]
      ],
      "config": {
        "component": "",
        "connector": "",
        "resource": "",
        "connector_name": ""
      },
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
      "description": "update change_entity",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "never be used change entity update test"
            }
          }
        ]
      ],
      "config": {},
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """

  Scenario: given update request with bad pattern should return error
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "event_pattern":[[
        {
          "field": "connector_bad",
          "cond": {
            "type": "eq",
            "value": "some"
          }
        }
      ]],
      "priority":0,
      "enabled":true,
      "config": {
        "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
        "on_success":"pass",
        "on_failure":"pass"
      },
      "external_data":{"entity":{"type":"entity"}}
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "event_pattern":[[
        {
          "field": "connector",
          "cond": {
            "type": "gt",
            "value": "some"
          }
        }
      ]],
      "priority":0,
      "enabled":true,
      "config": {
        "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
        "on_success":"pass",
        "on_failure":"pass"
      },
      "external_data":{"entity":{"type":"entity"}}
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "event_pattern":[[
        {
          "field": "extra.test",
          "field_type": "string",
          "cond": {
            "type": "gt",
            "value": "some"
          }
        }
      ]],
      "priority":0,
      "enabled":true,
      "config": {
        "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
        "on_success":"pass",
        "on_failure":"pass"
      },
      "external_data":{"entity":{"type":"entity"}}
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """

  Scenario: given update request with absent corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "corporate_entity_pattern": "test-pattern-not-exist",
      "priority":0,
      "enabled":true,
      "config": {
        "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
        "on_success":"pass",
        "on_failure":"pass"
      },
      "external_data":{"entity":{"type":"entity"}}
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """

  Scenario: given update request should update eventfilter without changes in old patterns,
            but should unset old patterns if new patterns are present
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-backward-compatibility-to-update:
    """
    {
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true,
      "old_patterns": [
        {
          "resource": {
            "regex_match": "test-eventfilter-to-backward-compatibility-to-update"
          }
        }
      ]
    }
    """
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-backward-compatibility-to-update:
    """
    {
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-backward-compatibility-to-update"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true,
      "old_patterns": null,
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-backward-compatibility-to-update"
            }
          }
        ]
      ]
    }
    """
