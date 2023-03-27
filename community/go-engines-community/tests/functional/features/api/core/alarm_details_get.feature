Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get details request should return alarm
    When I am admin
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "test-alarm-to-get-1",
        "opened": true,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-1",
        "steps": {
          "page": 2,
          "limit": 2
        }
      },
      {
        "_id": "test-alarm-to-get-1",
        "steps": {
          "limit": 5
        }
      },
      {
        "_id": "test-alarm-to-get-3",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-4",
        "opened": false,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-1",
        "opened": false,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "test-alarm-to-get-3",
        "opened": true,
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "not-exist",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "not-exist"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "_id": "test-alarm-to-get-1",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-1",
                "t": 1597030220,
                "initiator": "user"
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-2",
                "t": 1597030221,
                "initiator": "user"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-1",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-1",
                "t": 1597030220,
                "initiator": "user"
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-2",
                "t": 1597030221,
                "initiator": "user"
              }
            ],
            "meta": {
              "page": 2,
              "page_count": 2,
              "per_page": 2,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-1",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-1-output",
                "t": 1597030219,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-1",
                "t": 1597030220,
                "initiator": "user"
              },
              {
                "_t": "comment",
                "a": "root",
                "m": "test-alarm-to-get-1-comment-2",
                "t": 1597030221,
                "initiator": "user"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 5,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-3",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030221,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030221,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030241,
                "initiator": "external",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-3-output",
                "t": 1597030241,
                "initiator": "external",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 200,
        "_id": "test-alarm-to-get-4",
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030121,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030121,
                "initiator": "external",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030141,
                "initiator": "external",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-default.test-connector-default-name",
                "m": "test-alarm-to-get-4-output",
                "t": 1597030141,
                "initiator": "external",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      },
      {
        "status": 404,
        "_id": "test-alarm-to-get-1",
        "error": "Not found"
      },
      {
        "status": 404,
        "_id": "test-alarm-to-get-3",
        "error": "Not found"
      },
      {
        "status": 404,
        "_id": "not-exist",
        "error": "Not found"
      },
      {
        "status": 400,
        "_id": "not-exist",
        "errors": {
          "steps": "Steps is missing.",
          "children": "Children is missing."
        }
      }
    ]
    """

  @concurrent
  Scenario: given get details unauth request should not allow access
    When I do POST /api/v4/alarm-details
    Then the response code should be 401

  @concurrent
  Scenario: given get details request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/alarm-details
    Then the response code should be 403
