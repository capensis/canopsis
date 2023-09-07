Feature: update a pbehavior
  I need to be able to update a pbehavior
  Only admin should be able to update a pbehavior

  @concurrent
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

  @concurrent
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

  @concurrent
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

  @concurrent
  Scenario: given update request with rrule should return update rrule end
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-4:
    """json
    {
      "rrule": "FREQ=DAILY;UNTIL=20221108T103000Z",
      "tstart": {{ parseTimeTz "08-10-2022 10:00" }},
      "tstop": {{ parseTimeTz "08-10-2022 11:00" }},
      "enabled": true,
      "name": "test-pbehavior-to-update-4-name",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-4-pattern"
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
      "_id": "test-pbehavior-to-update-4",
      "rrule_end": {{ parseTimeTz "08-11-2022 10:00" }}
    }
    """

  @concurrent
  Scenario: given update request without rrule should return remove rrule end
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-5:
    """json
    {
      "tstart": {{ parseTimeTz "08-10-2022 10:00" }},
      "tstop": {{ parseTimeTz "08-10-2022 11:00" }},
      "enabled": true,
      "name": "test-pbehavior-to-update-5-name",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-5-pattern"
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
      "_id": "test-pbehavior-to-update-5",
      "rrule_end": null
    }
    """
