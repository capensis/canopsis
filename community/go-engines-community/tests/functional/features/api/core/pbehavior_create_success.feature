Feature: create a pbehavior
  I need to be able to create a pbehavior
  Only admin should be able to create a pbehavior

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
  Scenario: given create request with pause type and without stop should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled":true,
      "name": "test-pbehavior-to-create-4",
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "enabled":true,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-pbehavior-to-create-4",
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
              "value": "test-pbehavior-to-create-4-pattern"
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

  @concurrent
  Scenario: given create request with strange id should return ok
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "strange \\id&key=value!*@!'\"-_:;<>",
      "enabled":true,
      "name": "test-pbehavior-to-create-5",
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
              "value": "test-pbehavior-to-create-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do DELETE /api/v4/pbehaviors/{{ .lastResponse._id}}
    Then the response code should be 204

  @concurrent
  Scenario: given create request with rrule with until should return compute rrule end
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "rrule": "FREQ=DAILY;UNTIL=20221108T103000Z",
      "tstart": {{ parseTimeTz "08-10-2022 10:00" }},
      "tstop": {{ parseTimeTz "08-10-2022 11:00" }},
      "enabled": true,
      "name": "test-pbehavior-to-create-6",
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
    Then the response body should contain:
    """json
    {
      "rrule_end": {{ parseTimeTz "08-11-2022 10:00" }}
    }
    """

  @concurrent
  Scenario: given create request with rrule with count should return compute rrule end
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "rrule": "FREQ=DAILY;COUNT=30",
      "tstart": {{ parseTimeTz "08-10-2022 10:00" }},
      "tstop": {{ parseTimeTz "08-10-2022 11:00" }},
      "enabled": true,
      "name": "test-pbehavior-to-create-7",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-7-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "rrule_end": {{ parseTimeTz "06-11-2022 10:00" }}
    }
    """

  @concurrent
  Scenario: given create request with rrule without end should return empty rrule end
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "rrule": "FREQ=DAILY",
      "tstart": {{ parseTimeTz "08-10-2022 10:00" }},
      "tstop": {{ parseTimeTz "08-10-2022 11:00" }},
      "enabled": true,
      "name": "test-pbehavior-to-create-8",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-8-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "rrule_end": null
    }
    """

  @concurrent
  Scenario: given create request with rrule should return next event in get request
    When I am admin
    When I save response start={{ nowAdd "-1m" }}
    When I save response stop={{ nowAdd "1h" }}
    When I save response nextStart={{ nowAdd "1439m" }}
    When I save response nextStop={{ nowAdd "25h" }}
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "tstart": {{ .start }},
      "tstop": {{ .stop }},
      "rrule": "FREQ=DAILY",
      "enabled": true,
      "name": "test-pbehavior-to-create-9-1",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-9-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "tstart": {{ .start }},
      "tstop": {{ .stop }},
      "rrule": "FREQ=YEARLY",
      "enabled": true,
      "name": "test-pbehavior-to-create-9-2",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-9-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-create-9&sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-to-create-9-1",
          "tstart": {{ .nextStart }},
          "tstop": {{ .nextStop }}
        },
        {
          "name": "test-pbehavior-to-create-9-2",
          "tstart": {{ .start }},
          "tstop": {{ .stop }}
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
