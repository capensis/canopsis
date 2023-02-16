Feature: update a PBehavior
  I need to be able to update a PBehavior
  Only admin should be able to update a PBehavior

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 403

  Scenario: given update request should return ok
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-1:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-update-1-name",
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
              "value": "test-pbehavior-to-update-1-pattern"
            }
          }
        ]
      ],
      "exdates":[],
      "exceptions": []
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337,
      "color": "#FFFFFF",
      "enabled": true,
      "exceptions": [],
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "description": "test-reason-to-pbh-edit-description",
        "name": "test-reason-to-pbh-edit-name"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-1-pattern"
            }
          }
        ]
      ],
      "exdates":[],
      "comments": [
        {
          "_id": "test-pbehavior-to-update-1-comment-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-update-1-comment-1-message"
        },
        {
          "_id": "test-pbehavior-to-update-1-comment-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215338,
          "message": "test-pbehavior-to-update-1-comment-2-message"
        }
      ],
      "name": "test-pbehavior-to-update-1-name",
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
      },
      "last_alarm_date": 1592215337
    }
    """

  Scenario: given update request with pause type and without stop should return ok
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-2:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-update-2-name",
      "tstart": 1591172881,
      "tstop": null,
      "color": "",
      "type": "test-type-to-pbh-edit-3",
      "reason": "test-reason-to-pbh-edit",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337,
      "color": "",
      "enabled": true,
      "exceptions": [],
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "description": "test-reason-to-pbh-edit-description",
        "name": "test-reason-to-pbh-edit-name"
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
      "exdates":[],
      "name": "test-pbehavior-to-update-2-name",
      "rrule": "",
      "tstart": 1591172881,
      "tstop": null,
      "type": {
        "_id": "test-type-to-pbh-edit-3",
        "type": "pause"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-update-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337,
      "color": "",
      "enabled": true,
      "exceptions": [],
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "description": "test-reason-to-pbh-edit-description",
        "name": "test-reason-to-pbh-edit-name"
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
      "exdates":[],
      "comments": [],
      "name": "test-pbehavior-to-update-2-name",
      "rrule": "",
      "tstart": 1591172881,
      "tstop": null,
      "type": {
        "_id": "test-type-to-pbh-edit-3",
        "type": "pause"
      }
    }
    """

  Scenario: given update request with old mongo pattern should return ok
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-3:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-update-3-name",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update-3",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337,
      "color": "#FFFFFF",
      "enabled": true,
      "exceptions": [],
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "description": "test-reason-to-pbh-edit-description",
        "name": "test-reason-to-pbh-edit-name"
      },
      "old_mongo_query": {
        "$and": [
          {
            "name": "test filter"
          }
        ]
      },
      "exdates":[],
      "name": "test-pbehavior-to-update-3-name",
      "rrule": "",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-3:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "enabled": true,
      "name": "test-pbehavior-to-update-3-name",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update-3",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337,
      "enabled": true,
      "exceptions": [],
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "description": "test-reason-to-pbh-edit-description",
        "name": "test-reason-to-pbh-edit-name"
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
      "exdates":[],
      "name": "test-pbehavior-to-update-3-name",
      "rrule": "",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    Then the response key "old_mongo_query" should not exist

  Scenario: given update request with the name that already exists should cause dup error
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-1:
    """json
    {
      "enabled": true,
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
              "value": "test-pbehavior-to-update-1-pattern"
            }
          }
        ]
      ]
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

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-1:
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

  Scenario: given no exist pbehavior id should return error
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-not-exist:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-not-exist",
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
              "value": "test-pbehavior-not-exist"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
