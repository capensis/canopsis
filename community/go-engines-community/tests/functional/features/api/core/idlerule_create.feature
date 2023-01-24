Feature: Create a idle rule
  I need to be able to create a idle rule
  Only admin should be able to create a idle rule

  Scenario: given create alarm rule request should return ok
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-1-name",
      "description": "test-idle-rule-to-create-1-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-1-alarm"
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
              "value": "test-idle-rule-to-create-1-resource"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-1-operation-output",
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
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-1-name",
      "description": "test-idle-rule-to-create-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-1-alarm"
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
              "value": "test-idle-rule-to-create-1-resource"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-1-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given create entity rule request should return ok
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-2-name",
      "description": "test-idle-rule-to-create-2-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-create-2-resource"
            }
          }
        ]
      ],
      "disable_during_periods": ["pause"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-2-name",
      "description": "test-idle-rule-to-create-2-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-create-2-resource"
            }
          }
        ]
      ],
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-3-name",
      "description": "test-idle-rule-to-create-3-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 22,
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
              "value": "test-idle-rule-to-create-3-alarm"
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
              "value": "test-idle-rule-to-create-3-resource"
            }
          }
        ]
      ],
      "operation": {
        "type": "pbehavior",
        "parameters": {
          "name": "test-idle-rule-to-create-3-operation-name",
          "rrule": "FREQ=DAILY",
          "reason": "test-reason-to-edit-idle-rule",
          "type": "test-type-to-edit-idle-rule",
          "start_on_trigger": true,
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
    When I do GET /api/v4/idle-rules/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-3-name",
      "description": "test-idle-rule-to-create-3-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 22,
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
              "value": "test-idle-rule-to-create-3-alarm"
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
              "value": "test-idle-rule-to-create-3-resource"
            }
          }
        ]
      ],
      "operation": {
        "type": "pbehavior",
        "parameters": {
          "name": "test-idle-rule-to-create-3-operation-name",
          "rrule": "FREQ=DAILY",
          "reason": {
            "_id": "test-reason-to-edit-idle-rule",
            "name": "test-reason-to-edit-idle-rule-name",
            "description": "test-reason-to-edit-idle-rule-description",
            "created": 1592215337
          },
          "type": {
            "_id": "test-type-to-edit-idle-rule",
            "name": "test-type-to-edit-idle-rule-name",
            "description": "test-type-to-edit-idle-rule-description",
            "icon_name": "test-type-to-edit-idle-rule-icon"
          },
          "start_on_trigger": true,
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    
  Scenario: given create request with corporate entity pattern should return success
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-4-name",
      "description": "test-idle-rule-to-create-4-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-4-operation-output",
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
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-4-name",
      "description": "test-idle-rule-to-create-4-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-4-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    
  Scenario: given create request with corporate alarm pattern should return success
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-5-name",
      "description": "test-idle-rule-to-create-5-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-5-operation-output",
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
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-5-name",
      "description": "test-idle-rule-to-create-5-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-5-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    
  Scenario: given create request with both corporate patterns should return success
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-6-name",
      "description": "test-idle-rule-to-create-6-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",      
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-6-operation-output",
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
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-6-name",
      "description": "test-idle-rule-to-create-6-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-6-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    
  Scenario: given create request with absent alarm corporate pattern should return error
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-9-name",
      "description": "test-idle-rule-to-create-9-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-not-exist",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-9-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    
  Scenario: given create request with absent entity corporate pattern should return error
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-10-name",
      "description": "test-idle-rule-to-create-10-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_entity_pattern": "test-pattern-not-exist",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-10-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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

  Scenario: given create request with unacceptable alarm pattern and entity pattern fields for idle rules should return error
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-11-name",
      "description": "test-idle-rule-to-create-11-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-11-pattern"
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-11-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-11-name",
      "description": "test-idle-rule-to-create-11-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-11-pattern"
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-11-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-11-name",
      "description": "test-idle-rule-to-create-11-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-11-pattern"
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-11-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-11-name",
      "description": "test-idle-rule-to-create-11-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-11-pattern"
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-11-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-11-name",
      "description": "test-idle-rule-to-create-11-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
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
              "value": "test-idle-rule-to-create-11-pattern"
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-11-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-11-name",
      "description": "test-idle-rule-to-create-11-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-create-11-resource"
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
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-11-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
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
    
  Scenario: given create request with unacceptable alarm pattern and entity pattern fields for idle rules should return error
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-create-12-name",
      "description": "test-idle-rule-to-create-12-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_entity_pattern": "test-pattern-to-idle-rule-pattern-to-exclude-1",
      "corporate_alarm_pattern": "test-pattern-to-idle-rule-pattern-to-exclude-2",      
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-12-operation-output",
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
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-create-12-name",
      "description": "test-idle-rule-to-create-12-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-idle-rule-pattern-to-exclude-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-idle-rule-pattern-to-exclude-1",
      "corporate_entity_pattern_title": "test-pattern-to-idle-rule-pattern-to-exclude-1-title",
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
      "corporate_alarm_pattern": "test-pattern-to-idle-rule-pattern-to-exclude-2",
      "corporate_alarm_pattern_title": "test-pattern-to-idle-rule-pattern-to-exclude-2-title",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-create-12-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/idle-rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/idle-rules
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "duration.value": "Value is missing.",
        "duration.unit": "Unit is missing.",
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "priority": "Priority is missing.",
        "type": "Type is missing."
      }
    }
    """
    
  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/idle-rules:
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

  Scenario: given invalid create request with invalid type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "notexists"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [alarm entity]."
      }
    }
    """

  Scenario: given invalid create request with alarm type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "alarm",
      "operation": {
        "type": "notexists"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_condition": "AlarmCondition is missing.",
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required.",
        "operation.type": "Type must be one of [ack ackremove cancel assocticket changestate snooze pbehavior]."
      }
    }
    """

  Scenario: given invalid create request with entity type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "entity"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is missing."
      }
    }
    """

  Scenario: given invalid create request with entity type should return errors
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "entity",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "notexists"
            }
          }
        ]
      ],
      "operation": {
        "type": "notexists"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is not empty.",
        "entity_pattern": "EntityPattern is missing.",
        "operation": "Operation is not empty."
      }
    }
    """
    When I do POST /api/v4/idle-rules:
    """json
    {
      "type": "entity",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "operation": {
        "type": "notexists"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is not empty.",
        "entity_pattern": "EntityPattern is missing.",
        "operation": "Operation is not empty."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-to-check-unique-name-name"
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
