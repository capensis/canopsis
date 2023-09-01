Feature: update a pbehavior
  I need to be able to patch a pbehavior field individually
  Only admin should be able to patch a pbehavior

  @concurrent
  Scenario: given update request with name should update only if name is not empty and is unique
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """json
    {
        "name": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-pbehavior-to-patch-1-name"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """json
    {
      "name": "test-pbehavior-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """json
    {
      "name": ""
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """json
    {
      "name": "test-pbehavior-to-patch-1-name-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-pbehavior-to-patch-1-name-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "_id": "test-pbehavior-to-patch-1",
      "comments": [
        {
          "_id": "test-pbehavior-to-patch-1-comment-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-patch-1-comment-1-message"
        },
        {
          "_id": "test-pbehavior-to-patch-1-comment-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-patch-1-comment-2-message"
        }
      ],
      "created": 1592215337,
      "enabled": true,
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit",
          "created": 1592215037,
          "name": "Exception to pbehavior edit",
          "description": "test",
          "exdates": [
            {
              "begin": 15911648001,
              "end": 1591167901,
              "type": {
                "_id": "test-type-to-pbh-edit-1"
              }
            }
          ]
        }
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-patch-1-pattern"
            }
          }
        ]
      ],
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      },
      "rrule": "FREQ=DAILY",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      },
      "last_alarm_date": null
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-pbehavior-to-patch-1-name-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "_id": "test-pbehavior-to-patch-1",
      "comments": [
        {
          "_id": "test-pbehavior-to-patch-1-comment-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-patch-1-comment-1-message"
        },
        {
          "_id": "test-pbehavior-to-patch-1-comment-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-patch-1-comment-2-message"
        }
      ],
      "created": 1592215337,
      "enabled": true,
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit",
          "created": 1592215037,
          "name": "Exception to pbehavior edit",
          "description": "test",
          "exdates": [
            {
              "begin": 15911648001,
              "end": 1591167901,
              "type": {
                "_id": "test-type-to-pbh-edit-1"
              }
            }
          ]
        }
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-patch-1-pattern"
            }
          }
        ]
      ],
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      },
      "rrule": "FREQ=DAILY",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      },
      "last_alarm_date": null
    }
    """

  @concurrent
  Scenario: given update request with corporate entity pattern should update only if pattern is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """json
    {
      "corporate_entity_pattern": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-patch-2-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """json
    {
      "corporate_entity_pattern": ""
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """json
    {
      "corporate_entity_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2"
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
      ]
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
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
      ]
    }
    """

  @concurrent
  Scenario: given update request with entity pattern should update only if pattern is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-3:
    """json
    {
      "entity_pattern": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
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
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-3:
    """json
    {
      "entity_pattern": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-3:
    """json
    {
      "entity_pattern": [[]]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-3:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-3-pattern"
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
      "author": {
        "_id": "root",
        "name": "root"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist

  @concurrent
  Scenario: given update request with stop should update only if stop is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-4:
    """json
    {
      "tstop": null
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-4:
    """json
    {
      "tstop": 1591172880
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-4:
    """json
    {
      "tstop": 1591173981
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
      "tstop": 1591173981
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstop": 1591173981
    }
    """

  @concurrent
  Scenario: given update request with stop for pbehavior with pause type should update only if stop is empty or is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-5:
    """json
    {
      "tstop": 1591172880
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-5:
    """json
    {
      "tstop": 1591173982
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
      "tstop": 1591173982
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstop": 1591173982
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-5:
    """json
    {
      "tstop": null
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
      "tstop": null
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstop": null
    }
    """

  @concurrent
  Scenario: given update request with stop and type should update only if stop is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-6:
    """json
    {
      "type": "test-type-to-pbh-edit-1",
      "tstop": null
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-6:
    """json
    {
      "type": "test-type-to-pbh-edit-1",
      "tstop": 1591172880
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-6:
    """json
    {
      "type": "test-type-to-pbh-edit-1",
      "tstop": 1591173981
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
      "tstop": 1591173981
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstop": 1591173981
    }
    """

  @concurrent
  Scenario: given update request with stop for pbehavior with pause type should update only if stop is empty or is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-7:
    """json
    {
      "type": "test-type-to-pbh-edit-3",
      "tstop": 1591172880
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-7:
    """json
    {
      "type": "test-type-to-pbh-edit-3",
      "tstop": null
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
      "tstop": null
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstop": null
    }
    """

  @concurrent
  Scenario: given update request with start should update only if start is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-8:
    """json
    {
      "tstart": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "tstart": 1591172881
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-8:
    """json
    {
      "tstart": 1591172982
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstart": "Start should be less than Stop."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-8:
    """json
    {
      "tstart": 1591172681
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
      "tstart": 1591172681
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681
    }
    """

  @concurrent
  Scenario: given update request with start for pbehavior with pause type should update only if start is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-9:
    """json
    {
      "tstart": 1591172681
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
      "tstart": 1591172681
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681
    }
    """
