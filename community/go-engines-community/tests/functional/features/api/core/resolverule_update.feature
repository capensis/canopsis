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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-1-pattern-updated"
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
              "value": "test-resolve-rule-to-update-1-resource-updated"
            }
          }
        ]
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-1-pattern-updated"
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
              "value": "test-resolve-rule-to-update-1-resource-updated"
            }
          }
        ]
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-2-pattern-updated"
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
              "value": "test-resolve-rule-to-update-2-resource-updated"
            }
          }
        ]
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
    
  Scenario: given update request with corporate entity pattern should return success
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-2:
    """json
    {
      "name": "test-resolve-rule-to-update-2-name",
      "description": "test-resolve-rule-to-update-2-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-2-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
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
      "name": "test-resolve-rule-to-update-2-name",
      "description": "test-resolve-rule-to-update-2-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-2-pattern"
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
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    
  Scenario: given update request with corporate alarm pattern should return success
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-3:
    """json
    {
      "name": "test-resolve-rule-to-update-3-name",
      "description": "test-resolve-rule-to-update-3-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-1-resource"
            }
          }
        ]
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
      "name": "test-resolve-rule-to-update-3-name",
      "description": "test-resolve-rule-to-update-3-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
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
              "value": "test-resolve-rule-to-update-1-resource"
            }
          }
        ]
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    
  Scenario: given update request with both corporate patterns should return success
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-4:
    """json
    {
      "name": "test-resolve-rule-to-update-4-name",
      "description": "test-resolve-rule-to-update-4-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
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
      "name": "test-resolve-rule-to-update-4-name",
      "description": "test-resolve-rule-to-update-4-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
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
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    
  Scenario: given update request with both corporate entity pattern and custom entity pattern should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-5:
    """json
    {
      "name": "test-resolve-rule-to-update-5-name",
      "description": "test-resolve-rule-to-update-5-description",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-1-resource"
            }
          }
        ]
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "entity_pattern": "Can't be present both EntityPattern and CorporateEntityPattern."
      }
    }
    """
    
  Scenario: given update request with both corporate alarm pattern and custom alarm pattern should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-6:
    """json
    {
      "name": "test-resolve-rule-to-update-6-name",
      "description": "test-resolve-rule-to-update-6-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-1-pattern"
            }
          }
        ]
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "alarm_pattern": "Can't be present both AlarmPattern and CorporateAlarmPattern."
      }
    }
    """
    
  Scenario: given update request with absent alarm corporate pattern should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-7:
    """json
    {
      "name": "test-resolve-rule-to-update-7-name",
      "description": "test-resolve-rule-to-update-7-description",
      "corporate_alarm_pattern": "test-pattern-not-exist",
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "corporate_alarm_pattern": "CorporateAlarmPattern doesn't exist."
      }
    }
    """
    
  Scenario: given update request with absent alarm corporate pattern should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-8:
    """json
    {
      "name": "test-resolve-rule-to-update-8-name",
      "description": "test-resolve-rule-to-update-8-description",
      "corporate_entity_pattern": "test-pattern-not-exist",
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
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
        "priority": "Priority is missing.",
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required."
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

  Scenario: given update requests should update resolve rule without changes in old patterns,
            but should unset old patterns if new patterns are present
    When I am admin
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-description-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-resource-updated"
            }
          }
        ]
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
      "name": "test-resolve-rule-backward-compatibility-to-update-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-description-updated",
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-backward-compatibility-to-update"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-resource-updated"
            }
          }
        ]
      ],
      "old_entity_patterns": null,
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-description-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-resource-updated"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-component-updated"
            }
          }
        ]
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
      "name": "test-resolve-rule-backward-compatibility-to-update-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-description-updated",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-component-updated"
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
              "value": "test-resolve-rule-backward-compatibility-to-update-resource-updated"
            }
          }
        ]
      ],
      "old_entity_patterns": null,
      "old_alarm_patterns": null,
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
