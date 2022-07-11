Feature: Create an metaalarmrule
  I need to be able to create a metaalarmrule
  Only admin should be able to create a metaalarmrule

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/metaalarmrules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metaalarmrules
    Then the response code should be 403

  Scenario: given create request with alarm pattern should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]      
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]      
    }
    """

  Scenario: given create request with entity pattern should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-2",
      "auto_resolve": true,
      "name": "complex-test-2",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-2",
      "auto_resolve": true,
      "name": "complex-test-2",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-2",
      "auto_resolve": true,
      "name": "complex-test-2",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """

  Scenario: given create request with both patterns should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-3",
      "auto_resolve": true,
      "name": "complex-test-3",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-3",
      "auto_resolve": true,
      "name": "complex-test-3",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-3",
      "auto_resolve": true,
      "name": "complex-test-3",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """

  Scenario: given create request with corporate alarm pattern should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-4",
      "auto_resolve": true,
      "name": "complex-test-4",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-4",
      "auto_resolve": true,
      "name": "complex-test-4",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
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
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-4",
      "auto_resolve": true,
      "name": "complex-test-4",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
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
      ]
    }
    """

  Scenario: given create request with both corporate patterns should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-5",
      "auto_resolve": true,
      "name": "complex-test-5",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-5",
      "auto_resolve": true,
      "name": "complex-test-5",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
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
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-5",
      "auto_resolve": true,
      "name": "complex-test-5",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
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
      ]
    }
    """

  Scenario: given create request with total entity pattern should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-6",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "total_test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-6",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "total_test"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-6",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "total_test"
            }
          }
        ]
      ]
    }
    """

  Scenario: given create request with corporate total entity pattern should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-7",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-7",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_total_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "total_entity_pattern": [
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
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-7",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_total_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "total_entity_pattern": [
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

  Scenario: given create request with absent alarm corporate pattern should return error
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-8",
      "auto_resolve": true,
      "name": "complex-test-6",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "corporate_alarm_pattern": "test-pattern-not-exist"
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
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-8",
      "auto_resolve": true,
      "name": "complex-test-6",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "corporate_entity_pattern": "test-pattern-not-exist"
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

  Scenario: given create request with absent corporate total entity pattern should return error
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-8",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-rule-edit-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "corporate_total_entity_pattern": "CorporateTotalEntityPattern doesn't exist."
      }
    }
    """

  Scenario: given create request with wrong alarm patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "wrong-alarm-pattern-1",
      "auto_resolve": false,
      "name": "wrong-alarm-pattern-1",
      "config": {},
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """

  Scenario: given create request with wrong entity patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "address",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given create request with wrong total entity patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "total_entity_pattern": [
        [
          {
            "field": "address",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "total_entity_pattern": "TotalEntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given create request with wrong type should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "wrong-type-1",
      "auto_resolve": false,
      "name": "wrong-type-1",
      "config": {},
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "attribute_path"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "type": "Type must be one of [relation timebased attribute complex valuegroup corel]."
      }
    }
    """

  Scenario: given create request with wrong config type should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-wrong-config-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1,
        "value_paths": ["resource.path"]
      },
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ]      
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "config": "value_paths config can not be in type complex."
      }
    }
    """

  Scenario: given create request with empty patterns should return error
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "auto_resolve": false,
      "name": "test-attribute-type-1",
      "config": {},
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "attribute"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required."
      }
    }
    """

  Scenario: given create request with already exists id should return error
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "test-metaalarm-to-get-1",
      "auto_resolve": false,
      "name": "test-attribute-type-1",
      "config": {},
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test"
            }
          }
        ]
      ],
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "attribute"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """
