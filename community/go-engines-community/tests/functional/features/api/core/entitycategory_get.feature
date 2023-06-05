Feature: Get a entity category
  I need to be able to get a entity category
  Only admin should be able to get a entity category

  Scenario: given search request should return entity categorys
    When I am admin
    When I do GET /api/v4/entity-categories?search=test-category-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-category-to-get-1",
          "name": "test-category-to-get-1-name",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1592215337,
          "updated": 1592215337
        },
        {
          "_id": "test-category-to-get-2",
          "name": "test-category-to-get-2-name",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1592215337,
          "updated": 1592215337
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

  Scenario: given get request should return category
    When I am admin
    When I do GET /api/v4/entity-categories/test-category-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-category-to-get-1",
      "name": "test-category-to-get-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337,
      "updated": 1592215337
    }
    """

  Scenario: given sort request should return sorted categorys
    When I am admin
    When I do GET /api/v4/entity-categories?search=test-category-to-get&sort=desc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-category-to-get-2"
        },
        {
          "_id": "test-category-to-get-1"
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/entity-categories
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entity-categories
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/entity-categories/test-category-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entity-categories/test-category-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/entity-categories/test-category-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
