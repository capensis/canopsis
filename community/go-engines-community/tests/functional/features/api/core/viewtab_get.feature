Feature: Get a view tab
  I need to be able to get a view tab
  Only admin should be able to get a view tab

  Scenario: given get request should return tab
    When I am admin
    When I do GET /api/v4/view-tabs/test-tab-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-tab-to-get",
      "title": "test-tab-to-get-title",
      "author": "test-author-to-tab-get",
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/view-tabs/test-tab-to-get
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/view-tabs/test-tab-to-get
    Then the response code should be 403

  Scenario: given get request and auth user without view permissions should not allow access
    When I am admin
    When I do GET /api/v4/view-tabs/test-tab-to-check-access
    Then the response code should be 403

  Scenario: given get request with not exist id should return error
    When I am admin
    When I do GET /api/v4/view-tabs/test-tab-not-found
    Then the response code should be 404
