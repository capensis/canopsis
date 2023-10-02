Feature: Get link categories
  I need to be able to get link categories
  Only admin should be able to get link categories

  @concurrent
  Scenario: given search request should return link categories
    When I am admin
    When I do GET /api/v4/link-categories?type=alarm&search=test-category-to-alarm-link-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "categories": [
        "test-category-to-alarm-link-get-1",
        "test-category-to-alarm-link-get-2"
      ]
    }
    """
    When I do GET /api/v4/link-categories?type=entity&search=test-category-to-alarm-link-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "categories": [
        "test-category-to-alarm-link-get-1",
        "test-category-to-alarm-link-get-3"
      ]
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/link-categories
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/link-categories
    Then the response code should be 403

  @concurrent
  Scenario: given invalid request should return error
    When I am admin
    When I do GET /api/v4/link-categories?type=unknown-link-type
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type must be one of [alarm entity] or empty."
      }
    }
    """
