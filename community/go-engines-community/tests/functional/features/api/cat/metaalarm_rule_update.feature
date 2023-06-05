Feature: Update a metaalarmrule
  I need to be able to update a metaalarmrule
  Only admin should be able to update a metaalarmrule

  Scenario: given get request ad no auth user should not allow access
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1
    Then the response code should be 401

  Scenario: given get request ad auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1
    Then the response code should be 403

  Scenario: given update request should update metaalarmrule
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1:
    """json
    {
      "name": "test-metaalarm-to-update-1-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
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
      "_id": "test-metaalarm-to-update-1",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-1-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-1",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-1-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """

  Scenario: given update request with entity pattern should update metaalarmrule
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-2:
    """json
    {
      "name": "test-metaalarm-to-update-2",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-2-pattern-updated"
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
      "_id": "test-metaalarm-to-update-2",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-2-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-2",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-2-pattern-updated"
            }
          }
        ]
      ]
    }
    """

  Scenario: given update request with entity pattern should update metaalarmrule
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-3:
    """json
    {
      "name": "test-metaalarm-to-update-3",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-3-pattern-updated"
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
      "_id": "test-metaalarm-to-update-3",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-3",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-3-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-3",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-3",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-3-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    
  Scenario: given update request with both patterns should update metaalarmrule
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-4:
    """json
    {
      "name": "test-metaalarm-to-update-4",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-4-pattern-updated"
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
              "value": "test-pattern-to-update-4-pattern-updated"
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
      "_id": "test-metaalarm-to-update-4",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-4",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-4-pattern-updated"
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
              "value": "test-pattern-to-update-4-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-4",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-4",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-4-pattern-updated"
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
              "value": "test-pattern-to-update-4-pattern-updated"
            }
          }
        ]
      ]
    }
    """    

  Scenario: given update request with corporate pattern should update metaalarmrule
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-5:
    """json
    {
      "name": "test-metaalarm-to-update-5",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
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
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-5",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-5",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
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
    """json
    {
      "_id": "test-metaalarm-to-update-5",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-5",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
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

  Scenario: given update request with both corporate patterns should update metaalarmrule
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-6:
    """json
    {
      "name": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
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
    """json
    {
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
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

  Scenario: given update request with total entity patterns should be ok
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-6:
    """json
    {
      "name": "test-metaalarm-to-update-6-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-6-pattern"
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
              "value": "test-pattern-to-update-6-pattern-total"
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
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-6-pattern"
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
              "value": "test-pattern-to-update-6-pattern-total"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-6-pattern"
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
              "value": "test-pattern-to-update-6-pattern-total"
            }
          }
        ]
      ]
    }
    """

  Scenario: given update request with corporate total entity patterns should be ok
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-7:
    """json
    {
      "name": "test-metaalarm-to-update-7-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-7-pattern"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-to-update-7",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-7-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-7-pattern"
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
    """json
    {
      "_id": "test-metaalarm-to-update-7",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-7-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-7-pattern"
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

  Scenario: given update request with absent alarm corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "corporate_alarm_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_alarm_pattern": "CorporateAlarmPattern doesn't exist."
      }
    }
    """

  Scenario: given update request with absent entity corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "corporate_entity_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """

  Scenario: given update request with absent total corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
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
      "corporate_total_entity_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_total_entity_pattern": "CorporateTotalEntityPattern doesn't exist."
      }
    }
    """

  Scenario: given update request with wrong alarm pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """

  Scenario: given update request with wrong entity pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
      "entity_pattern": [
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given update request with wrong total entity pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "total_entity_pattern": "TotalEntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given update request with wrong type should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "not-exist",
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
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [relation timebased attribute complex valuegroup corel]."
      }
    }
    """

  Scenario: given update request with wrong type should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
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
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1,
        "value_paths": ["resource.path"]
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "config": "value_paths config can not be in type complex."
      }
    }
    """

  Scenario: given update request with not exist id should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-not-exist:
    """json
    {
      "name": "test-metaalarm-to-update-1-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 10
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
    Then the response body should contain:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given update request with empty patterns should return error
    When I am admin
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required."
      }
    }
    """

  Scenario: given update request should update metaalarmrule without changes in old patterns,
            but should unset old patterns if new patterns are present
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-update:
    """json
    {
      "name": "test-metaalarm-rule-backward-compatibility-to-update-name-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 3
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-update
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-rule-backward-compatibility-to-update",
      "auto_resolve": false,
      "config": {
        "threshold_count": 3
      },
      "name": "test-metaalarm-rule-backward-compatibility-to-update-name-updated",
      "type": "complex",
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-metaalarm-rule-backward-compatibility-component-to-update"
          }
        }
      ],
      "old_entity_patterns": [
        {
          "name": {
            "regex_match": "test-metaalarm-rule-backward-compatibility-component-to-update"
          }
        }
      ],
      "old_total_entity_patterns": [
        {
          "name": {
            "regex_match": "test-metaalarm-rule-backward-compatibility-component-to-update"
          }
        }
      ],
      "old_event_patterns": [
        {
          "resource": {
            "regex_match": "test-metaalarm-rule-backward-compatibility-component-to-update"
          }
        }
      ]
    }
    """
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-update:
    """json
    {
      "name": "test-metaalarm-rule-backward-compatibility-to-update-name-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 3
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
    Then the response code should be 200
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-update
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-rule-backward-compatibility-to-update",
      "auto_resolve": false,
      "config": {
        "threshold_count": 3
      },
      "name": "test-metaalarm-rule-backward-compatibility-to-update-name-updated",
      "type": "complex",
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
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-metaalarm-rule-backward-compatibility-component-to-update"
          }
        }
      ],
      "old_entity_patterns": null,
      "old_event_patterns": null,
      "old_total_entity_patterns": [
        {
          "name": {
            "regex_match": "test-metaalarm-rule-backward-compatibility-component-to-update"
          }
        }
      ]
    }
    """
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-update:
    """json
    {
      "name": "test-metaalarm-rule-backward-compatibility-to-update-name-updated",
      "auto_resolve": false,
      "type": "complex",
      "config": {
        "threshold_count": 3
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
              "value": "test"
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
              "value": "test"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-update
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-rule-backward-compatibility-to-update",
      "auto_resolve": false,
      "config": {
        "threshold_count": 3
      },
      "name": "test-metaalarm-rule-backward-compatibility-to-update-name-updated",
      "type": "complex",
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
              "value": "test"
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
              "value": "test"
            }
          }
        ]
      ],
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_event_patterns": null
    }
    """

  Scenario: given update request with unacceptable alarm pattern and entity pattern fields for metaalarm rules should return error
    When I am admin
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
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
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
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
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
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
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
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
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
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
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
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
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "valuegroup",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.some.value"
        ]
      },
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-metaalarm-rule-to-create-9-pattern"
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
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "total_entity_pattern": "TotalEntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given update request with unacceptable corporate alarm pattern and corporate entity pattern fields for metaalarm rules should exclude invalid patterns
    When I am admin
    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """json
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
      "corporate_entity_pattern": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1",
      "corporate_total_entity_pattern": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1",
      "corporate_alarm_pattern": "test-pattern-to-metaalarm-rule-pattern-to-exclude-2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1
      },
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1",
      "corporate_entity_pattern_title": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1-title",
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
      "corporate_alarm_pattern": "test-pattern-to-metaalarm-rule-pattern-to-exclude-2",
      "corporate_alarm_pattern_title": "test-pattern-to-metaalarm-rule-pattern-to-exclude-2-title",
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1-pattern"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1",
      "corporate_total_entity_pattern_title": "test-pattern-to-metaalarm-rule-pattern-to-exclude-1-title"
    }
    """

  Scenario: given update request should allow to update value group metaalarmrule without patterns
            and all old patterns should be removed
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-valuegroup-rule-rate-backward-compatibility-to-update:
    """json
    {
      "name": "test-valuegroup-rule-rate-backward-compatibility-to-update-1",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.some.value"
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/metaalarmrules/test-valuegroup-rule-rate-backward-compatibility-to-update
    Then the response body should contain:
    """json
    {
      "_id": "test-valuegroup-rule-rate-backward-compatibility-to-update",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "old_event_patterns": null,
      "old_total_entity_patterns": [
        {
          "name": {
            "regex_match": "test-valuegroup-rule-rate-backward-compatibility-to-update"
          }
        }
      ]
    }
    """
