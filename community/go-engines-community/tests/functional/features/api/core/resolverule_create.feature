Feature: Create an resolve rule
  I need to be able to create a resolve rule
  Only admin should be able to create a resolve rule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-1-name",
      "description": "test-resolve-rule-to-create-1-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-1-pattern"
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
              "value": "test-resolve-rule-to-create-1-resource"
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-resolve-rule-to-create-1-name",
      "description": "test-resolve-rule-to-create-1-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-1-pattern"
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
              "value": "test-resolve-rule-to-create-1-resource"
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
    When I do GET /api/v4/resolve-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-resolve-rule-to-create-1-name",
      "description": "test-resolve-rule-to-create-1-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-1-pattern"
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
              "value": "test-resolve-rule-to-create-1-resource"
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

  Scenario: given create request should update priority of next rules
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-2-priority-1-name",
      "description": "test-resolve-rule-to-create-2-priority-1-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-2-priority-1-pattern"
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
              "value": "test-resolve-rule-to-create-2-priority-1-resource"
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
    Then the response code should be 201
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-2-priority-2-name",
      "description": "test-resolve-rule-to-create-2-priority-2-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-2-priority-2-pattern"
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
              "value": "test-resolve-rule-to-create-2-priority-2-resource"
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
    Then the response code should be 201
    When I do GET /api/v4/resolve-rules?search=test-resolve-rule-to-create-2&sort_by=priority
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resolve-rule-to-create-2-priority-2-name"
        },
        {
          "name": "test-resolve-rule-to-create-2-priority-1-name"
        }
      ]
    }
    """
    
  Scenario: given create request with corporate entity pattern should return success
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-3-name",
      "description": "test-resolve-rule-to-create-3-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-3-pattern"
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-resolve-rule-to-create-3-name",
      "description": "test-resolve-rule-to-create-3-description",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-3-pattern"
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
    
  Scenario: given create request with corporate alarm pattern should return success
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-4-name",
      "description": "test-resolve-rule-to-create-4-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-1-resource"
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-resolve-rule-to-create-4-name",
      "description": "test-resolve-rule-to-create-4-description",
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
              "value": "test-resolve-rule-to-create-1-resource"
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
    
  Scenario: given create request with both corporate patterns should return success
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-5-name",
      "description": "test-resolve-rule-to-create-5-description",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-resolve-rule-to-create-5-name",
      "description": "test-resolve-rule-to-create-5-description",
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
    
  Scenario: given create request with absent alarm corporate pattern should return error
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-8-name",
      "description": "test-resolve-rule-to-create-8-description",
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
    
  Scenario: given create request with absent alarm corporate pattern should return error
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-9-name",
      "description": "test-resolve-rule-to-create-9-description",
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
    
  Scenario: given create request with unacceptable alarm pattern and entity pattern fields for resolve rules should return error
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-10-name",
      "description": "test-resolve-rule-to-create-10-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-10-pattern"
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
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-10-name",
      "description": "test-resolve-rule-to-create-10-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-10-pattern"
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
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-10-name",
      "description": "test-resolve-rule-to-create-10-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-10-pattern"
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
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-10-name",
      "description": "test-resolve-rule-to-create-10-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-10-pattern"
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
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-10-name",
      "description": "test-resolve-rule-to-create-10-description",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-10-pattern"
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
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-10-name",
      "description": "test-resolve-rule-to-create-10-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-create-10-resource"
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

  Scenario: given create request with unacceptable corporate alarm pattern and corporate entity pattern fields for resolve rules should exclude invalid patterns
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-to-create-11-name",
      "description": "test-resolve-rule-to-create-11-description",
      "corporate_entity_pattern": "test-pattern-to-resolve-rule-pattern-to-exclude-1",
      "corporate_alarm_pattern": "test-pattern-to-resolve-rule-pattern-to-exclude-2",
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-resolve-rule-to-create-11-name",
      "description": "test-resolve-rule-to-create-11-description",
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

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
    }
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

  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "wrong_field",
            "cond": {
              "type": "eq",
              "value": "ram"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "wrong_field",
            "cond": {
              "type": "eq",
              "value": "ram"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/resolve-rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/resolve-rules
    Then the response code should be 403

  Scenario: given create request with already exists id and name should return error
    When I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-to-check-unique"
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
    When I do POST /api/v4/resolve-rules:
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
