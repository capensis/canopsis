Feature: get a PBehavior
  I need to be able to get a PBehavior
  Only admin should be able to get a PBehavior

  Scenario: given get all request should return pbehaviors
    When I am admin
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-get-by-name&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-name-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "comments": [
            {
              "_id": "test-pbehavior-to-get-by-name-1-comment-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "ts": 1592215337,
              "message": "test-pbehavior-to-get-by-name-1-comment-1-message"
            },
            {
              "_id": "test-pbehavior-to-get-by-name-1-comment-2",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "ts": 1592215337,
              "message": "test-pbehavior-to-get-by-name-1-comment-2-message"
            }
          ],
          "color": "#FFFFFF",
          "created": 1592215337,
          "updated": 1592215337,
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
                    "_id": "test-type-to-pbh-edit-1",
                    "description": "Pbh edit 1 State type",
                    "icon_name": "test-to-pbh-edit-1-icon",
                    "name": "Pbh edit 1 State",
                    "priority": 11,
                    "type": "active"
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
                  "value": "test-pbehavior-to-get-by-name-1-pattern"
                }
              }
            ]
          ],
          "name": "test-pbehavior-to-get-by-name-1-name",
          "reason": {
            "_id": "test-reason-to-pbh-edit",
            "name": "test-reason-to-pbh-edit-name",
            "description": "test-reason-to-pbh-edit-description"
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
          },
          "last_alarm_date": null,
          "origin": "",
          "editable": true
        },
        {
          "_id": "test-pbehavior-to-get-by-name-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "comments": [],
          "color": "",
          "created": 1592215337,
          "updated": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "old_mongo_query": {
            "$and": [
              {
                "name": "test filter"
              }
            ]
          },
          "name": "test-pbehavior-to-get-by-name-2-name",
          "reason": {
            "_id": "test-reason-to-pbh-edit",
            "name": "test-reason-to-pbh-edit-name",
            "description": "test-reason-to-pbh-edit-description"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": null,
          "type": {
            "_id": "test-type-to-pbh-edit-2",
            "description": "Pbh edit 2 State type",
            "icon_name": "test-to-pbh-edit-2-icon",
            "name": "Pbh edit 2 State",
            "priority": 12,
            "type": "pause"
          },
          "last_alarm_date": null,
          "origin": "",
          "editable": true
        },
        {
          "_id": "test-pbehavior-to-get-by-name-3",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "comments": [],
          "color": "#FFFFFF",
          "created": 1592215337,
          "updated": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-pbehavior-to-get-by-name-3-pattern"
                }
              }
            ]
          ],
          "name": "test-pbehavior-to-get-by-name-3-name",
          "reason": {
            "_id": "test-reason-to-pbh-edit",
            "name": "test-reason-to-pbh-edit-name",
            "description": "test-reason-to-pbh-edit-description"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": null,
          "type": {
            "_id": "test-type-to-pbh-edit-2",
            "description": "Pbh edit 2 State type",
            "icon_name": "test-to-pbh-edit-2-icon",
            "name": "Pbh edit 2 State",
            "priority": 12,
            "type": "pause"
          },
          "last_alarm_date": null,
          "origin": "test-pbehavior-to-get-by-name-3-origin",
          "editable": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given search request by type should return pbehaviors
    When I am admin
    When I do GET /api/v4/pbehaviors?search=test-type-to-get-pbehavior-name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-type",
          "type": {
            "_id": "test-type-to-get-pbehavior",
            "description": "test-type-to-get-pbehavior-description",
            "icon_name": "test-type-to-get-pbehavior-icon",
            "name": "test-type-to-get-pbehavior-name",
            "type": "active"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/pbehaviors?search=type.name="test-type-to-get-pbehavior-name"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-type",
          "type": {
            "_id": "test-type-to-get-pbehavior",
            "description": "test-type-to-get-pbehavior-description",
            "icon_name": "test-type-to-get-pbehavior-icon",
            "name": "test-type-to-get-pbehavior-name",
            "type": "active"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given search request by reason should return pbehaviors
    When I am admin
    When I do GET /api/v4/pbehaviors?search=test-reason-to-pbehavior-get-name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-reason",
          "reason": {
            "_id": "test-reason-to-pbehavior-get",
            "description": "test-reason-to-pbehavior-get-description",
            "name": "test-reason-to-pbehavior-get-name"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/pbehaviors?search=reason.name="test-reason-to-pbehavior-get-name"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-reason",
          "reason": {
            "_id": "test-reason-to-pbehavior-get",
            "description": "test-reason-to-pbehavior-get-description",
            "name": "test-reason-to-pbehavior-get-name"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given sotred get all request should return pbehaviors
    When I am admin
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-get-by-name&sort_by=type.name&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-name-2"
        },
        {
          "_id": "test-pbehavior-to-get-by-name-3"
        },
        {
          "_id": "test-pbehavior-to-get-by-name-1"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/pbehaviors
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/pbehaviors
    Then the response code should be 403

  Scenario: given get request should return pbehavior
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-by-name-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-pbehavior-to-get-by-name-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "comments": [
        {
          "_id": "test-pbehavior-to-get-by-name-1-comment-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-get-by-name-1-comment-1-message"
        },
        {
          "_id": "test-pbehavior-to-get-by-name-1-comment-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "ts": 1592215337,
          "message": "test-pbehavior-to-get-by-name-1-comment-2-message"
        }
      ],
      "color": "#FFFFFF",
      "created": 1592215337,
      "updated": 1592215337,
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
                "_id": "test-type-to-pbh-edit-1",
                "description": "Pbh edit 1 State type",
                "icon_name": "test-to-pbh-edit-1-icon",
                "color": "#2FAB63",
                "name": "Pbh edit 1 State",
                "priority": 11,
                "type": "active"
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
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "color": "#2FAB63",
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
              "value": "test-pbehavior-to-get-by-name-1-pattern"
            }
          }
        ]
      ],
      "name": "test-pbehavior-to-get-by-name-1-name",
      "reason": {
        "_id": "test-reason-to-pbh-edit",
        "name": "test-reason-to-pbh-edit-name",
        "description": "test-reason-to-pbh-edit-description"
      },
      "rrule": "",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1",
        "description": "Pbh edit 1 State type",
        "icon_name": "test-to-pbh-edit-1-icon",
        "color": "#2FAB63",
        "name": "Pbh edit 1 State",
        "priority": 11,
        "type": "active"
      },
      "last_alarm_date": null,
      "origin": ""
    }
    """

  Scenario: given invalid get request should return not found error
    When I am admin
    When I do GET /api/v4/pbehaviors/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get
    Then the response code should be 403
