Feature: Get a permission
  I need to be able to get a permission
  Only admin should be able to get a permission

  Scenario: given search request should return permissions
    When I am admin
    When I do GET /api/v4/permissions?search=test-permission-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-permission-to-get-1",
          "description": "test-permission-to-get-1-description",
          "name": "test-permission-to-get-1",
          "type": ""
        },
        {
          "_id": "test-permission-to-get-2",
          "description": "test-permission-to-get-2-description",
          "name": "test-permission-to-get-2",
          "type": ""
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

  Scenario: given sort request should return sorted permissions
    When I am admin
    When I do GET /api/v4/permissions?search=test-permission-to-get&sort=desc&sort_by=description
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-permission-to-get-2"
        },
        {
          "_id": "test-permission-to-get-1"
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
    When I do GET /api/v4/permissions
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/permissions
    Then the response code should be 403
