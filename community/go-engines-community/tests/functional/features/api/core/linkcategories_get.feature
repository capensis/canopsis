Feature: Get link categories
  I need to be able to get link categories
  Only admin should be able to get link categories

  @concurrent
  Scenario: given search request should return link categories
    When I am admin
    When I do GET /api/v4/link-categories?type=alarm
    Then the response code should be 200
    Then the response array key "categories" should contain:
    """
    [
        "test-category-to-alarm-link-get-1",
        "test-category-to-alarm-link-get-2"
    ]
    """

  @concurrent
  Scenario: given get categories request and no auth user should not allow access
    When I do GET /api/v4/link-categories
    Then the response code should be 401

  @concurrent
  Scenario: given get categories request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/link-categories
    Then the response code should be 403

  @concurrent
  Scenario: given get categories request and auth user without permissions should not allow access
    When I am admin
    When I do GET /api/v4/link-categories?type=unknown-link-type
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "type": "Type must be one of [alarm entity] or empty."
      }
    }
    """

  @concurrent
  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/link-categories?type=entity&limit=3
    Then the response code should be 200
    Then the response array key "categories" should contain only:
    """
    [
        "test-category-to-alarm-link-get-1",
        "test-category-to-alarm-link-get-3",
        "test-link-rule-to-get-3-link-1-category"
    ]
    """
