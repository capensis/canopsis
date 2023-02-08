Feature: create a PBehavior
  I need to be able to create a PBehavior
  Only admin should be able to create a PBehavior

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/pbehaviors
    Then the response code should be 401

  Scenario: given create request and auth user without view permission should not allow access
    When I am noperms
    When I do POST /api/v4/pbehaviors
    Then the response code should be 403

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-create-1",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-1-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin":  1591164001,
          "end":  1591167601,
          "type":  "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions":  ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "enabled": true,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-pbehavior-to-create-1",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type":  {
        "_id": "test-type-to-pbh-edit-1"
      },
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-1-pattern"
            }
          }
        ]
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
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """
    When I do GET /api/v4/pbehaviors/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "comments": [],
      "color": "#FFFFFF",
      "enabled": true,
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit",
          "created": 1592215037,
          "description": "test",
          "exdates": [
            {
              "begin": 15911648001,
              "end": 1591167901,
              "type": {
                "_id": "test-type-to-pbh-edit-1",
                "description": "Pbh edit 1 State type",
                "icon_name": "test-to-pbh-edit-1-icon",
                "name": "Pbh edit 1 State",
                "priority": 11,
                "type": "active"
              }
            }
          ],
          "name": "Exception to pbehavior edit"
        }
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "name": "Pbh edit 1 State",
            "priority": 11,
            "type": "active"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-1-pattern"
            }
          }
        ]
      ],
      "name": "test-pbehavior-to-create-1",
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "description": "test-reason-to-pbh-edit-description",
        "name": "test-reason-to-pbh-edit-name"
      },
      "rrule": "",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1",
        "description": "Pbh edit 1 State type",
        "icon_name": "test-to-pbh-edit-1-icon",
        "name": "Pbh edit 1 State",
        "priority": 11,
        "type": "active"
      }
    }
    """

  Scenario: given create request with corporate pattern should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-create-2",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "exdates": [
        {
          "begin":  1591164001,
          "end":  1591167601,
          "type":  "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions":  ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "enabled": true,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-pbehavior-to-create-2",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type":  {
        "_id": "test-type-to-pbh-edit-1"
      },
      "reason": {
        "_id": "test-reason-to-pbh-edit"
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
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ],
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """

  Scenario: given create request with custom id should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "test-pbehavior-to-create-3",
      "enabled":true,
      "name": "test-pbehavior-to-create-3-name",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-3-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-create-3
    Then the response code should be 200

  Scenario: given create request with the custom ID that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "test-pbehavior-to-check-unique",
      "enabled":true,
      "name": "test-pbehavior-to-create-4",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-4-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request with the name that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled":true,
      "name": "test-pbehavior-to-check-unique-name",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-4-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
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

  Scenario: given create request with pause type and without stop should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled":true,
      "name": "test-pbehavior-to-create-5",
      "tstart": 1591172881,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-3",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-5-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "enabled":true,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-pbehavior-to-create-5",
      "tstart": 1591172881,
      "tstop": null,
      "color": "#FFFFFF",
      "type": {
        "_id": "test-type-to-pbh-edit-3"
      },
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-5-pattern"
            }
          }
        ]
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
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "reason": "Reason is missing.",
        "tstart": "Start is missing.",
        "type": "Type is missing."
      }
    }
    """

  Scenario: given invalid create request with start > stop should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "tstart": 1592172881,
      "tstop": 1591536400
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """

  Scenario: given invalid create request with not existed reason should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "reason": "notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "reason": "Reason doesn't exist."
      }
    }
    """

  Scenario: given invalid create request with invalid exclude date should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "exdates": [
        {
          "begin": 1592164001,
          "end": 1592166001,
          "type": "test-type-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "exdates": "Exdates doesn't exist."
      }
    }
    """

  Scenario: given invalid create request with invalid exception should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "exceptions": ["test-exception-to-pbh-edit", "test-exception-not-exist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """

  Scenario: given invalid create request with not existed type should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "type": "notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type doesn't exist."
      }
    }
    """

  Scenario: given invalid create request with invalid entity pattern should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "entity_pattern": [[]]
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-8-pattern"
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-rule-edit-1",
      "enabled":true,
      "name": "test-pbehavior-to-create-7",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit"
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "corporate_entity_pattern": "test-pattern-not-found",
      "enabled":true,
      "name": "test-pbehavior-to-create-7",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit"
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

  Scenario: given invalid create request with missing stop should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "tstart": 1591172881,
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """

  Scenario: given invalid create request with invalid id should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "invalid/id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "invalidid?key=value"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "$invalidid"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "invalid.id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """

  Scenario: given create request with strange id should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "strange \\id&key=value!*@!'\"-_:;<>",
      "enabled":true,
      "name": "test-pbehavior-to-create-6",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-6-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do DELETE /api/v4/pbehaviors/{{ .lastResponse._id}}
    Then the response code should be 204

  Scenario: given invalid create request with invalid color should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "color": "notcolor"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """
