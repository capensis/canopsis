Feature: update a PBehavior
  I need to be able to patch a PBehavior field individually
  Only admin should be able to patch a PBehavior

  Scenario: given update request and no auth user should not allow access
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 403

  Scenario: given no exist pbehavior id should return error
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-not-exist:
    """json
    {
      "name": "test-pbehavior-not-exist"
    }
    """
    Then the response code should be 404

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

  Scenario: given update request with start, stop and type should update only if all fields are valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-10:
    """json
    {
      "tstart": null,
      "tstop": 1591172881,
      "type": "test-type-to-pbh-edit-1"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-10:
    """json
    {
      "tstart": 1591172681,
      "tstop": null,
      "type": "test-type-to-pbh-edit-1"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-10:
    """json
    {
      "tstart": 1591172881,
      "tstop": 1591172681,
      "type": "test-type-to-pbh-edit-1"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-10:
    """json
    {
      "tstart": 1591172681,
      "tstop": 1591172881,
      "type": "test-type-to-pbh-edit-1"
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
      "tstart": 1591172681,
      "tstop": 1591172881,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-10
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681,
      "tstop": 1591172881,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-10:
    """json
    {
      "tstart": 1591172881,
      "tstop": 1591172681,
      "type": "test-type-to-pbh-edit-3"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-10:
    """json
    {
      "tstart": 1591172681,
      "tstop": null,
      "type": "test-type-to-pbh-edit-3"
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
      "tstart": 1591172681,
      "tstop": null,
      "type": {
        "_id": "test-type-to-pbh-edit-3"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-10
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681,
      "tstop": null,
      "type": {
        "_id": "test-type-to-pbh-edit-3"
      }
    }
    """

  Scenario: given update request with type should update only if type is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-11:
    """json
    {
      "type": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-11:
    """json
    {
      "type": "test-type-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-11:
    """json
    {
      "type": "test-type-to-pbh-edit-2"
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
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-11
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """

  Scenario: given update request with type for pbehavior without stop should update only if type is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-12:
    """json
    {
      "type": "test-type-to-pbh-edit-1"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-12:
    """json
    {
      "type": "test-type-to-pbh-edit-2"
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
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-12
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """

  Scenario: given update request with reason should update only if reason is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "reason": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "reason": "test-reason-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "reason": "Reason doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "reason": "test-reason-to-pbh-edit"
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
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      }
    }
    """

  Scenario: given update request with enabled should update only if enabled is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "enabled": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "enabled": false
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
      "enabled": false
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": false
    }
    """

  Scenario: given update request with rrule should update only if rrule is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "rrule": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "rrule": "FREQ=DAILY"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "rrule": "invalid"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "rrule": "RRule is invalid recurrence rule."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "rrule": "FREQ=WEEKLY"
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
      "rrule": "FREQ=WEEKLY"
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "rrule": "FREQ=WEEKLY"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "rrule": ""
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
      "rrule": ""
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "rrule": ""
    }
    """

  Scenario: given update request with exdates should update only if exdates is empty or is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exdates": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exdates": [{}]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exdates.0.begin": "Begin is missing.",
        "exdates.0.end": "End is missing.",
        "exdates.0.type": "Type is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exdates": "Exdates doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exdates": [
        {
          "begin": 1592164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exdates": [
        {
          "begin": 1591164002,
          "end": 1591167602,
          "type": "test-type-to-pbh-edit-2"
        }
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
      "exdates": [
        {
          "begin": 1591164002,
          "end": 1591167602,
          "type": {
            "_id": "test-type-to-pbh-edit-2"
          }
        }
      ]
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exdates": [
        {
          "begin": 1591164002,
          "end": 1591167602,
          "type": {
            "_id": "test-type-to-pbh-edit-2"
          }
        }
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exdates": []
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
      "exdates": []
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exdates": []
    }
    """

  Scenario: given update request with exceptions should update only if exceptions is empty or is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exceptions": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exceptions": [""]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exceptions": ["test-exception-not-exist"]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exceptions": []
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
      "exceptions": []
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exceptions": []
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-13:
    """json
    {
      "exceptions": ["test-exception-to-pbh-edit"]
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
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """

  Scenario: given update request with color should update only if color is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-14:
    """json
    {
      "color": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "#FFFFFF"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-14:
    """json
    {
      "color": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-14:
    """json
    {
      "color": ""
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": ""
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-14
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": ""
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-14:
    """json
    {
      "color": "#AAAAAA"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "#AAAAAA"
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-14
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "#AAAAAA"
    }
    """
