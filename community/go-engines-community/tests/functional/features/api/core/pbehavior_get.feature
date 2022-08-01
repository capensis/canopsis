Feature: get a PBehavior
  I need to be able to get a PBehavior
  Only admin should be able to get a PBehavior

  Scenario: given get all request should return pbehaviors
    When I am admin
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-get-by-name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-name-1",
          "author": "root",
          "comments": [
            {
              "_id": "test-comment-1",
              "author": "root",
              "ts": 1592215337,
              "message": "qwerty"
            },
            {
              "_id": "test-comment-2",
              "author": "root",
              "ts": 1592215337,
              "message": "asdasd"
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
                    "priority": 10,
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
                "priority": 10,
                "type": "active"
              }
            }
          ],
          "filter": {
            "$and": [
              {
                "name": "test filter"
              }
            ]
          },
          "name": "test-pbehavior-to-get-by-name-1-name",
          "reason": {
            "_id": "test-reason-1",
            "name": "test-reason-1-name",
            "description": "test-reason-1-description"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "name": "Pbh edit 1 State",
            "priority": 10,
            "type": "active"
          },
          "last_alarm_date": null
        },
        {
          "_id": "test-pbehavior-to-get-by-name-2",
          "author": "root",
          "comments": [],
          "color": "#FFFFFF",
          "created": 1592215337,
          "updated": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "filter": {
            "$and": [
              {
                "name": "test filter"
              }
            ]
          },
          "name": "test-pbehavior-to-get-by-name-2-name",
          "reason": {
            "_id": "test-reason-1",
            "name": "test-reason-1-name",
            "description": "test-reason-1-description"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-2",
            "description": "Pbh edit 2 State type",
            "icon_name": "test-to-pbh-edit-2-icon",
            "name": "Pbh edit 2 State",
            "priority": 11,
            "type": "active"
          },
          "last_alarm_date": null
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

  Scenario: given get all request should return pbehaviors with test filter pattern in filter field
    When I am admin
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-get-by-filter-filter
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-get-by-filter",
          "author": "root",
          "comments": [
            {
              "_id": "test-comment-1",
              "author": "root",
              "ts": 1592215337,
              "message": "qwerty"
            },
            {
              "_id": "test-comment-2",
              "author": "root",
              "ts": 1592215337,
              "message": "asdasd"
            }
          ],
          "color": "#FFFFFF",
          "created": 1592215338,
          "updated": 1592215338,
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
                    "priority": 10,
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
                "priority": 10,
                "type": "active"
              }
            }
          ],
          "filter": {
            "$and": [
              {
                "name": "test-pbehavior-to-get-by-filter-filter"
              }
            ]
          },
          "name": "test-pbehavior-to-get-by-filter-name",
          "reason": {
            "_id": "test-reason-1",
            "name": "test-reason-1-name",
            "description": "test-reason-1-description"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "name": "Pbh edit 1 State",
            "priority": 10,
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
          "author": "root",
          "comments": [],
          "color": "#FFFFFF",
          "created": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "filter": {"$and": [{"name": "ccccc"}]},
          "name": "test-pbehavior-to-get-by-type-name",
          "reason": {
            "_id": "test-reason-1",
            "description": "test-reason-1-description",
            "name": "test-reason-1-name"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-get-pbehavior",
            "description": "test-type-to-get-pbehavior-description",
            "icon_name": "test-type-to-get-pbehavior-icon",
            "name": "test-type-to-get-pbehavior-name",
            "priority": 25,
            "type": "active"
          },
          "updated": 1592215337
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
          "author": "root",
          "comments": [],
          "color": "#FFFFFF",
          "created": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "filter": {"$and": [{"name": "ccccc"}]},
          "name": "test-pbehavior-to-get-by-type-name",
          "reason": {
            "_id": "test-reason-1",
            "description": "test-reason-1-description",
            "name": "test-reason-1-name"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-get-pbehavior",
            "description": "test-type-to-get-pbehavior-description",
            "icon_name": "test-type-to-get-pbehavior-icon",
            "name": "test-type-to-get-pbehavior-name",
            "priority": 25,
            "type": "active"
          },
          "updated": 1592215337
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
          "author": "root",
          "comments": [],
          "color": "#FFFFFF",
          "created": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "filter": {"$and": [{"name": "ccccc"}]},
          "name": "test-pbehavior-to-get-by-reason-name",
          "reason": {
            "_id": "test-reason-to-pbehavior-get",
            "description": "test-reason-to-pbehavior-get-description",
            "name": "test-reason-to-pbehavior-get-name"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "name": "Pbh edit 1 State",
            "priority": 10,
            "type": "active"
          },
          "updated": 1592215337
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
          "author": "root",
          "comments": [],
          "color": "#FFFFFF",
          "created": 1592215337,
          "enabled": true,
          "exceptions": [],
          "exdates": [],
          "filter": {"$and": [{"name": "ccccc"}]},
          "name": "test-pbehavior-to-get-by-reason-name",
          "reason": {
            "_id": "test-reason-to-pbehavior-get",
            "description": "test-reason-to-pbehavior-get-description",
            "name": "test-reason-to-pbehavior-get-name"
          },
          "rrule": "",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "name": "Pbh edit 1 State",
            "priority": 10,
            "type": "active"
          },
          "updated": 1592215337
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
          "_id": "test-pbehavior-to-get-by-name-1"
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

  Scenario: GET a PBehavior but unauthorized
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get
    Then the response code should be 401

  Scenario: GET a PBehavior but without permissions
    When I am noperms
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get
    Then the response code should be 403

  Scenario: Get a PBehavior with success
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-by-name-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-pbehavior-to-get-by-name-1",
      "author": "root",
      "comments": [
        {
          "_id": "test-comment-1",
          "author": "root",
          "ts": 1592215337,
          "message": "qwerty"
        },
        {
          "_id": "test-comment-2",
          "author": "root",
          "ts": 1592215337,
          "message": "asdasd"
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
                "priority": 10,
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
            "priority": 10,
            "type": "active"
          }
        }
      ],
      "filter": {
        "$and": [
          {
            "name": "test filter"
          }
        ]
      },
      "name": "test-pbehavior-to-get-by-name-1-name",
      "reason": {
        "_id": "test-reason-1",
        "name": "test-reason-1-name",
        "description": "test-reason-1-description"
      },
      "rrule": "",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1",
        "description": "Pbh edit 1 State type",
        "icon_name": "test-to-pbh-edit-1-icon",
        "name": "Pbh edit 1 State",
        "priority": 10,
        "type": "active"
      },
      "last_alarm_date": null
    }
    """

  Scenario: Get a PBehavior with not found response
    When I am admin
    When I do GET /api/v4/pbehaviors/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
