Feature: Update a idle rule
  I need to be able to update a idle rule
  Only admin should be able to update a idle rule

  Scenario: given update idle rule request should return ok
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-1:
    """json
    {
      "name": "test-idle-rule-to-update-1-name",
      "description": "test-idle-rule-to-update-1-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-update-1-alarm-updated"
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
              "value": "test-idle-rule-to-update-1-resource-updated"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-update-1",
      "name": "test-idle-rule-to-update-1-name",
      "description": "test-idle-rule-to-update-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-update-1-alarm-updated"
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
              "value": "test-idle-rule-to-update-1-resource-updated"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    
  Scenario: given update request with corporate alarm and entity patterns should return success
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-2:
    """json
    {
      "name": "test-idle-rule-to-update-2-name",
      "description": "test-idle-rule-to-update-2-description",
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
          "output": "test-idle-rule-to-update-2-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-update-2-name",
      "description": "test-idle-rule-to-update-2-description",
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
          "output": "test-idle-rule-to-update-2-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given update request with absent corporate alarm pattern should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-not-exist",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-update-3-operation-output",
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

  Scenario: given update request with absent corporate entity pattern should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 20,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "corporate_entity_pattern": "test-pattern-to-rule-edit-not-exist",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-update-3-operation-output",
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
    
  Scenario: given update request with unacceptable alarm pattern and entity pattern fields for idle rules should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
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
              "value": "test-idle-rule-to-update-3-pattern"
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
          "output": "test-idle-rule-to-update-3-operation-output",
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
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
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
              "value": "test-idle-rule-to-update-3-pattern"
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
          "output": "test-idle-rule-to-update-3-operation-output",
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
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
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
              "value": "test-idle-rule-to-update-3-pattern"
            }
          },
          {
            "field": "v.idled",
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
          "output": "test-idle-rule-to-update-3-operation-output",
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
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
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
              "value": "test-idle-rule-to-update-3-pattern"
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
          "output": "test-idle-rule-to-update-3-operation-output",
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
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
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
              "value": "test-idle-rule-to-update-3-pattern"
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
          "output": "test-idle-rule-to-update-3-operation-output",
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
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-3:
    """json
    {
      "name": "test-idle-rule-to-update-3-name",
      "description": "test-idle-rule-to-update-3-description",
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
              "value": "test-idle-rule-to-update-3-resource"
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
          "output": "test-idle-rule-to-update-3-operation-output",
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

  Scenario: given update request with unacceptable alarm pattern and entity pattern fields for idle rules should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-4:
    """json
    {
      "name": "test-idle-rule-to-update-4-name",
      "description": "test-idle-rule-to-update-4-description",
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
          "output": "test-idle-rule-to-update-4-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-idle-rule-to-update-4-name",
      "description": "test-idle-rule-to-update-4-description",
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
          "output": "test-idle-rule-to-update-4-operation-output",
          "duration": {
            "value": 3,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["pause"]
    }
    """

  Scenario: given update request with invalid patterns format should return bad request
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-5:
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

  Scenario: given invalid update request with invalid type should return errors
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-5:
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

  Scenario: given invalid update request with alarm type should return errors
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-5:
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

  Scenario: given invalid update request with entity type should return errors
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update-5:
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
    
  Scenario: given update requests should update alarm idle rule without changes in old patterns,
            but should unset old patterns if new patterns are present
    When I am admin
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-1:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-1-name",
      "description": "test-idle-rule-to-backward-compatibility-update-1-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }        
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-1",
      "name": "test-idle-rule-to-backward-compatibility-update-1-name",
      "description": "test-idle-rule-to-backward-compatibility-update-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "old_alarm_patterns": [
        {
          "_id": "test-idle-rule-to-backward-compatibility-update-1-alarm"
        }
      ],
      "old_entity_patterns": [
        {
          "name": "test-idle-rule-to-backward-compatibility-update-1-resource"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-1:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-1-name",
      "description": "test-idle-rule-to-backward-compatibility-update-1-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-1-alarm-updated"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }        
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-1",
      "name": "test-idle-rule-to-backward-compatibility-update-1-name",
      "description": "test-idle-rule-to-backward-compatibility-update-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-1-alarm-updated"
            }
          }
        ]
      ],
      "old_alarm_patterns": null,
      "old_entity_patterns": [
        {
          "name": "test-idle-rule-to-backward-compatibility-update-1-resource"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """    
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-1:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-1-name",
      "description": "test-idle-rule-to-backward-compatibility-update-1-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-1-alarm-updated"
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
              "value": "test-idle-rule-to-update-1-resource-updated"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }        
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-1",
      "name": "test-idle-rule-to-backward-compatibility-update-1-name",
      "description": "test-idle-rule-to-backward-compatibility-update-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-1-alarm-updated"
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
              "value": "test-idle-rule-to-update-1-resource-updated"
            }
          }
        ]
      ],
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-1-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    
  Scenario: given update requests should update idle rule without changes in old patterns,
            but should unset old patterns if new corporate patterns are present
    When I am admin
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-2:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-2-name",
      "description": "test-idle-rule-to-backward-compatibility-update-2-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-2-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-2",
      "name": "test-idle-rule-to-backward-compatibility-update-2-name",
      "description": "test-idle-rule-to-backward-compatibility-update-2-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
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
      "old_entity_patterns": null,
      "old_alarm_patterns": [
        {
          "_id": "test-idle-rule-to-backward-compatibility-update-2-alarm"
        }
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-2-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-2:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-2-name",
      "description": "test-idle-rule-to-backward-compatibility-update-2-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-2-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-2",
      "name": "test-idle-rule-to-backward-compatibility-update-2-name",
      "description": "test-idle-rule-to-backward-compatibility-update-2-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
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
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-to-backward-compatibility-update-2-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    
  Scenario: given update requests should update entity idle rule without changes in old patterns,
            but should unset old patterns if new patterns are present
    When I am admin
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-3:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-3-name",
      "description": "test-idle-rule-to-backward-compatibility-update-3-description-updated",
      "type": "entity",
      "enabled": true,
      "priority": 3000,
      "duration": {
        "value": 5,
        "unit": "s"
      }
    }    
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-3",
      "name": "test-idle-rule-to-backward-compatibility-update-3-name",
      "description": "test-idle-rule-to-backward-compatibility-update-3-description-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "entity",
      "enabled": true,
      "priority": 3000,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "old_entity_patterns": [
        {
          "name": "test-idle-rule-to-backward-compatibility-update-3-resource"
        }
      ]
    }
    """
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-3:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-3-name",
      "description": "test-idle-rule-to-backward-compatibility-update-3-description-updated",
      "type": "entity",
      "enabled": true,
      "priority": 3000,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-3-resource"
            }
          }
        ]
      ]      
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-3",
      "name": "test-idle-rule-to-backward-compatibility-update-3-name",
      "description": "test-idle-rule-to-backward-compatibility-update-3-description-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "entity",
      "enabled": true,
      "priority": 3000,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "old_entity_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-3-resource"
            }
          }
        ]
      ]
    }
    """
    
  Scenario: given update requests should update entity idle rule without changes in old patterns,
            but should unset old patterns if new corporate patterns are present
    When I am admin
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-4:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-4-name",
      "description": "test-idle-rule-to-backward-compatibility-update-4-description-updated",
      "type": "entity",
      "enabled": true,
      "priority": 3000,
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "duration": {
        "value": 5,
        "unit": "s"
      }
    }    
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-4",
      "name": "test-idle-rule-to-backward-compatibility-update-4-name",
      "description": "test-idle-rule-to-backward-compatibility-update-4-description-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "entity",
      "enabled": true,
      "priority": 3000,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "old_entity_patterns": null,
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
      ]
    }
    """    
    
  Scenario: given update requests should update alarm idle rule, but should unset old_alarm_pattern because of changing to entity type
    When I am admin
    Then I do PUT /api/v4/idle-rules/test-idle-rule-to-backward-compatibility-update-5:
    """json
    {
      "name": "test-idle-rule-to-backward-compatibility-update-5-name",
      "description": "test-idle-rule-to-backward-compatibility-update-5-description",
      "type": "entity",
      "enabled": true,
      "priority": 30,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-backward-compatibility-update-5",
      "name": "test-idle-rule-to-backward-compatibility-update-5-name",
      "description": "test-idle-rule-to-backward-compatibility-update-5-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "entity",
      "enabled": true,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-backward-compatibility-update-5-pattern"
            }
          }
        ]
      ]
    }
    """    

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/test-idle-rule-to-update:
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
    
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/idle-rules/notexist
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/idle-rules/notexist
    Then the response code should be 403

  Scenario: given no exist idle rule id should return error
    When I am admin
    When I do PUT /api/v4/idle-rules/notexist:
    """json
    {
      "name": "test-idle-rule-notexists-name",
      "description": "test-idle-rule-notexists-description",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 31,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-update-alarm-updated"
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
              "value": "test-idle-rule-to-update-resource-updated"
            }
          }
        ]
      ],
      "operation": {
        "type": "snooze",
        "parameters": {
          "output": "test-idle-rule-notexists-operation-output-updated",
          "duration": {
            "value": 5,
            "unit": "s"
          }
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 404
