Feature: Update and delete corporate pattern should affect metaalarm rule models
  Scenario: given metaalarm rule and corporate pattern update and delete actions should update patterns in metaalarm rule
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "_id": "metaalarm-rule-pattern-1",
      "auto_resolve": true,
      "name": "metaalarm-pattern-update-1",
      "type": "complex",
      "output_template": "",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "corporate_alarm_pattern": "test-pattern-to-rule-metaalarm-corporate-update-3",
      "corporate_entity_pattern": "test-pattern-to-rule-metaalarm-corporate-update-2",
      "corporate_total_entity_pattern": "test-pattern-to-rule-metaalarm-corporate-update-1"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/metaalarmrules/metaalarm-rule-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "auto_resolve": true,
      "name": "metaalarm-pattern-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "output_template": "",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "corporate_total_entity_pattern": "test-pattern-to-rule-metaalarm-corporate-update-1",
      "corporate_total_entity_pattern_title": "test-pattern-to-rule-metaalarm-corporate-update-1-title",
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-metaalarm-corporate-update-1-pattern"
            }
          }
        ]
      ],      
      "corporate_entity_pattern": "test-pattern-to-rule-metaalarm-corporate-update-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-metaalarm-corporate-update-2-title",      
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-metaalarm-corporate-update-2-pattern"
            }
          }
        ]
      ],
      "corporate_alarm_pattern": "test-pattern-to-rule-metaalarm-corporate-update-3",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-metaalarm-corporate-update-3-title",      
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-metaalarm-corporate-update-3-pattern"
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
      ]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-rule-metaalarm-corporate-update-1:
    """json
    {
      "title": "new total entity pattern title",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new total entity pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/patterns/test-pattern-to-rule-metaalarm-corporate-update-2:
    """json
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
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/patterns/test-pattern-to-rule-metaalarm-corporate-update-3:
    """json
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
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/metaalarmrules/metaalarm-rule-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "auto_resolve": true,
      "name": "metaalarm-pattern-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "output_template": "",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "corporate_total_entity_pattern": "test-pattern-to-rule-metaalarm-corporate-update-1",
      "corporate_total_entity_pattern_title": "new total entity pattern title",
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new total entity pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-metaalarm-corporate-update-2",
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
      "corporate_alarm_pattern": "test-pattern-to-rule-metaalarm-corporate-update-3",
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
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-metaalarm-corporate-update-1
    Then the response code should be 204
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-metaalarm-corporate-update-2
    Then the response code should be 204
    When I do DELETE /api/v4/patterns/test-pattern-to-rule-metaalarm-corporate-update-3
    Then the response code should be 204
    When I do GET /api/v4/cat/metaalarmrules/metaalarm-rule-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "auto_resolve": true,
      "name": "metaalarm-pattern-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "output_template": "",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new total entity pattern"
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "new alarm pattern"
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_alarm_pattern" should not exist
    Then the response key "corporate_alarm_pattern_title" should not exist
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
    Then the response key "corporate_total_entity_pattern" should not exist
    Then the response key "corporate_total_entity_pattern_title" should not exist
