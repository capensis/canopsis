Feature: Get a alarm tag
  I need to be able to read a alarm tag
  Only admin should be able to read a alarm tag

  @concurrent
  Scenario: given get all request should return alarm tags
    When I am admin
    When I do GET /api/v4/alarm-tags?search=test-alarm-tag-to-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-tag-to-get-1",
          "type": 1,
          "value": "test-alarm-tag-to-get-1-value",
          "color": "#AABBCC",
          "alarm_pattern": [
            [
              {
                "field": "v.state.val",
                "cond": {
                  "type": "eq",
                  "value": 3
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-alarm-tag-to-get-1-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1612139798,
          "updated": 1612139798
        },
        {
          "_id": "test-alarm-tag-to-get-2",
          "type": 0,
          "value": "test-alarm-tag-to-get-2-value",
          "color": "#AABBCC",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1612139798,
          "updated": 1612139798
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

  @concurrent
  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/alarm-tags
    Then the response code should be 401

  @concurrent
  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarm-tags
    Then the response code should be 403

  @concurrent
  Scenario: given get request should return alarm tag
    When I am admin
    When I do GET /api/v4/alarm-tags/test-alarm-tag-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-alarm-tag-to-get-1",
      "type": 1,
      "value": "test-alarm-tag-to-get-1-value",
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.state.val",
            "cond": {
              "type": "eq",
              "value": 3
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-get-1-pattern"
            }
          }
        ]
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1612139798,
      "updated": 1612139798
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/alarm-tags/notexist
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarm-tags/notexist
    Then the response code should be 403

  @concurrent
  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/alarm-tags/notexist
    Then the response code should be 404
