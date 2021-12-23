Feature: Get a filter
  I need to be able to get a filter
  Only admin should be able to get a filter

  Scenario: given search request should return filters
    When I am admin
    When I do GET /api/v4/cat/filters?search=test-filter-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-filter-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "name": "test-filter-to-get-1-name",
          "entity_patterns": [
            {
              "name": "test-filter-to-get-1-resource"
            }
          ],
          "updated": 1619083733
        },
        {
          "_id": "test-filter-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "name": "test-filter-to-get-2-name",
          "entity_patterns": [
            {
              "name": "test-filter-to-get-2-resource"
            }
          ],
          "updated": 1619083733
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

  Scenario: given get request should return filter
    When I am admin
    When I do GET /api/v4/cat/filters/test-filter-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-filter-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "name": "test-filter-to-get-1-name",
      "entity_patterns": [
        {
          "name": "test-filter-to-get-1-resource"
        }
      ],
      "updated": 1619083733
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/filters
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/filters
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/filters/test-filter-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/filters/test-filter-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/filters/test-filter-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
