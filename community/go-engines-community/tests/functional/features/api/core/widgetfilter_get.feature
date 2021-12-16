Feature: Get a widget filter
  I need to be able to get a widget filter
  Only admin should be able to get a widget filter

  Scenario: given get request should return filter
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widgetfilter-to-get-1",
      "title": "test-widgetfilter-to-get-1-title",
      "query": "{\"test\":\"test\"}",
      "author": "test-user-to-widgetfilter-edit",
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  Scenario: given personal get request should return filter
    When I am test-role-to-widgetfiler-edit
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widgetfilter-to-get-2",
      "title": "test-widgetfilter-to-get-2-title",
      "query": "{\"test\":\"test\"}",
      "author": "test-user-to-widgetfilter-edit",
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-1
    Then the response code should be 403

  Scenario: given get request and auth user without view permissions should not allow access
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-check-access
    Then the response code should be 403

  Scenario: given get request and someone else filter should return error
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-2
    Then the response code should be 404

  Scenario: given get request with not exist id should return error
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-not-found
    Then the response code should be 404
