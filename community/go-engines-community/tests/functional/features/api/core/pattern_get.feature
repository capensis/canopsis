Feature: Get a saved pattern
  I need to be able to get a saved pattern
  Only admin should be able to get a saved pattern

  Scenario: given search request should return corporate and not corporate patterns
    When I am noperms
    When I do GET /api/v4/patterns?search=test-pattern-to-get&sort_by=title
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-pattern-to-get-1",
          "title": "test-pattern-to-get-1-title",
          "type": "alarm",
          "is_corporate": false,
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-1-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-pattern-to-get-3",
          "title": "test-pattern-to-get-3-title",
          "type": "entity",
          "is_corporate": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-3-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-pattern-to-get-4",
          "title": "test-pattern-to-get-4-title",
          "type": "pbehavior",
          "is_corporate": true,
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-4-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
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

  Scenario: given search request should return corporate patterns
    When I am noperms
    When I do GET /api/v4/patterns?corporate=true&search=test-pattern-to-get&sort_by=title
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-pattern-to-get-3",
          "title": "test-pattern-to-get-3-title",
          "type": "entity",
          "is_corporate": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-3-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-pattern-to-get-4",
          "title": "test-pattern-to-get-4-title",
          "type": "pbehavior",
          "is_corporate": true,
          "pbehavior_pattern": [
            [
              {
                "field": "pbehavior_info.type",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-4-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
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

  Scenario: given search request should return not corporate patterns
    When I am noperms
    When I do GET /api/v4/patterns?corporate=false&search=test-pattern-to-get&sort_by=title
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-pattern-to-get-1",
          "title": "test-pattern-to-get-1-title",
          "type": "alarm",
          "is_corporate": false,
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-1-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992
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

  Scenario: given filter request should return patterns by type
    When I am noperms
    When I do GET /api/v4/patterns?search=test-pattern-to-get&sort_by=title&type=alarm
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-pattern-to-get-1",
          "title": "test-pattern-to-get-1-title",
          "type": "alarm",
          "is_corporate": false,
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-1-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992
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

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/patterns
    Then the response code should be 401

  Scenario: given get request should return user pattern
    When I am noperms
    When I do GET /api/v4/patterns/test-pattern-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-pattern-to-get-1",
      "title": "test-pattern-to-get-1-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-get-1-pattern"
            }
          }
        ]
      ],
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "created": 1605263992,
      "updated": 1605263992
    }
    """

  Scenario: given get request and another user should return not found
    When I am noperms
    When I do GET /api/v4/patterns/test-pattern-to-get-2
    Then the response code should be 404

  Scenario: given get request should return corporate pattern
    When I am admin
    When I do GET /api/v4/patterns/test-pattern-to-get-3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-pattern-to-get-3",
      "title": "test-pattern-to-get-3-title",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-get-3-pattern"
            }
          }
        ]
      ],
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "created": 1605263992,
      "updated": 1605263992
    }
    """

  Scenario: given get request with not exist id should return not found error
    When I am noperms
    When I do GET /api/v4/patterns/test-pattern-notexist
    Then the response code should be 404

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/patterns/test-pattern-notexist
    Then the response code should be 401
