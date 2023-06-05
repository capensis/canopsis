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

  Scenario: given update default rule request should update resolve rule
    When I am admin
    Then I do PUT /api/v4/resolve-rules/default_rule:
    """json
    {
      "name": "Default rule",
      "description": "Default rule",
      "duration": {
        "value": 200,
        "unit": "h"
      },
      "priority": 1
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "Default rule",
      "description": "Default rule",
      "duration": {
        "value": 200,
        "unit": "h"
      },
      "priority": 1
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
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update-1:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-1-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-1-description-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-1-resource-updated"
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
      "name": "test-resolve-rule-backward-compatibility-to-update-1-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-1-description-updated",
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-backward-compatibility-to-update-1"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-1-resource-updated"
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
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update-1:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-1-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-1-description-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-1-resource-updated"
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
              "value": "test-resolve-rule-backward-compatibility-to-update-1-component-updated"
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
      "name": "test-resolve-rule-backward-compatibility-to-update-1-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-1-description-updated",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-backward-compatibility-to-update-1-component-updated"
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
              "value": "test-resolve-rule-backward-compatibility-to-update-1-resource-updated"
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
    
  Scenario: given update request should update resolve rule with old patterns if new patterns are absent
    When I am admin
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update-2:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-2-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-2-description-updated",
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
      "name": "test-resolve-rule-backward-compatibility-to-update-2-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-2-description-updated",
      "old_entity_patterns": [
        {
          "name": {
            "regex_match": "test-resolve-rule-backward-compatibility-to-update-2"
          }
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """ 
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update-3:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-3-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-3-description-updated",
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
      "name": "test-resolve-rule-backward-compatibility-to-update-3-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-3-description-updated",
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-backward-compatibility-to-update-3"
          }
        }
      ],
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """    

  Scenario: given update requests should update resolve rule without changes in old patterns,
            but should unset old patterns if new corporate patterns are present
    When I am admin
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update-4:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-4-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-4-description-updated",
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
      "created": 1619083733,
      "name": "test-resolve-rule-backward-compatibility-to-update-4-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-4-description-updated",
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-resolve-rule-backward-compatibility-to-update-4"
          }
        }
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
      "old_entity_patterns": null,
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then I do PUT /api/v4/resolve-rules/test-resolve-rule-backward-compatibility-to-update-4:
    """json
    {
      "name": "test-resolve-rule-backward-compatibility-to-update-4-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-4-description-updated",
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
      "created": 1619083733,
      "name": "test-resolve-rule-backward-compatibility-to-update-4-name-updated",
      "description": "test-resolve-rule-backward-compatibility-to-update-4-description-updated",
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

      "old_entity_patterns": null,
      "old_alarm_patterns": null,
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """

  Scenario: given update request with unacceptable alarm pattern and entity pattern fields for resolve rules should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-9:
    """json
    {
      "name": "test-resolve-rule-to-update-9-name",
      "description": "test-resolve-rule-to-update-9-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-9-pattern"
            }
          },
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
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
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-9:
    """json
    {
      "name": "test-resolve-rule-to-update-9-name",
      "description": "test-resolve-rule-to-update-9-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-9-pattern"
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
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
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-9:
    """json
    {
      "name": "test-resolve-rule-to-update-9-name",
      "description": "test-resolve-rule-to-update-9-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-9-pattern"
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
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
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-9:
    """json
    {
      "name": "test-resolve-rule-to-update-9-name",
      "description": "test-resolve-rule-to-update-9-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-9-pattern"
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
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
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-9:
    """json
    {
      "name": "test-resolve-rule-to-update-9-name",
      "description": "test-resolve-rule-to-update-9-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-9-pattern"
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
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
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-9:
    """json
    {
      "name": "test-resolve-rule-to-update-9-name",
      "description": "test-resolve-rule-to-update-9-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-update-9-resource"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
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
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given update request with unacceptable alarm pattern and entity pattern fields for resolve rules should return error
    When I am admin
    When I do PUT /api/v4/resolve-rules/test-resolve-rule-to-update-10:
    """json
    {
      "name": "test-resolve-rule-to-update-10-name",
      "description": "test-resolve-rule-to-update-10-description",
      "corporate_entity_pattern": "test-pattern-to-resolve-rule-pattern-to-exclude-1",
      "corporate_alarm_pattern": "test-pattern-to-resolve-rule-pattern-to-exclude-2",
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
      "name": "test-resolve-rule-to-update-10-name",
      "description": "test-resolve-rule-to-update-10-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-resolve-rule-pattern-to-exclude-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-resolve-rule-pattern-to-exclude-1",
      "corporate_entity_pattern_title": "test-pattern-to-resolve-rule-pattern-to-exclude-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.state.val",
            "cond": {
              "type": "eq",
              "value": 3
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
      ],
      "corporate_alarm_pattern": "test-pattern-to-resolve-rule-pattern-to-exclude-2",
      "corporate_alarm_pattern_title": "test-pattern-to-resolve-rule-pattern-to-exclude-2-title"
    }
    """
