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
    """
    {
      "name": "test-metaalarm-to-update-1-updated",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-1",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-1-updated",
      "author": "root",
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-1",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-1-updated",
      "author": "root",
      "type": "complex",
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
    """
    {
      "name": "test-metaalarm-to-update-2",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-2",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-2",
      "author": "root",
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-2",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-2",
      "author": "root",
      "type": "complex",
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
    """
    {
      "name": "test-metaalarm-to-update-3",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-3",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-3",
      "author": "root",
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-3",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-3",
      "author": "root",
      "type": "complex",
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
    """
    {
      "name": "test-metaalarm-to-update-4",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-4",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-4",
      "author": "root",
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-4",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-4",
      "author": "root",
      "type": "complex",
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
    """
    {
      "name": "test-metaalarm-to-update-5",
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
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-metaalarm-to-update-5",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-5",
      "author": "root",
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
      "_id": "test-metaalarm-to-update-5",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-5",
      "author": "root",
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
    """
    {
      "name": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "type": "complex",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6",
      "author": "root",
      "type": "complex",
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
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6",
      "author": "root",
      "type": "complex",
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
    """
    {
      "name": "test-metaalarm-to-update-6-updated",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6-updated",
      "author": "root",
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-6",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-6-updated",
      "author": "root",
      "type": "complex",
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
    """
    {
      "name": "test-metaalarm-to-update-7-updated",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-7",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-7-updated",
      "author": "root",
      "type": "complex",
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
    """
    {
      "_id": "test-metaalarm-to-update-7",
      "auto_resolve": false,
      "name": "test-metaalarm-to-update-7-updated",
      "author": "root",
      "type": "complex",
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

  Scenario: given update request with corporate total entity patterns and custom total entity pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8-updated",
      "auto_resolve": false,
      "type": "complex",
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-8-pattern"
            }
          }
        ]
      ],
      "corporate_total_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "total_entity_pattern": "Can't be present both TotalEntityPattern and CorporateTotalEntityPattern."
      }
    }
    """

  Scenario: given update request with absent alarm corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
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

  Scenario: given update request with absent entity corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
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

  Scenario: given update request with absent total corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
      "corporate_total_entity_pattern": "test-pattern-not-exist"
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

  Scenario: given update request with wrong alarm pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
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
    """
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """

  Scenario: given update request with wrong entity pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
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
    """
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given update request with wrong total entity pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "errors": {
        "total_entity_pattern": "TotalEntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given update request with corporate entity pattern and custom entity pattern should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
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
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2"
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

  Scenario: given update request with wrong type should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "type": "Type must be one of [relation timebased attribute complex valuegroup corel]."
      }
    }
    """

  Scenario: given update request with wrong type should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-8:
    """
    {
      "name": "test-metaalarm-to-update-8",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "errors": {
        "config": "value_paths config can not be in type complex."
      }
    }
    """

  Scenario: given update request with not exist id should return error
    When I am admin
    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-not-exist:
    """
    {
      "name": "test-metaalarm-to-update-1-updated",
      "auto_resolve": false,
      "type": "complex",
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
    """
    {
      "error": "Not found"
    }
    """
