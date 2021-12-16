Feature: Update a widget filter
  I need to be able to update a widget filter
  Only admin should be able to update a widget filter

  Scenario: given update request should update filter
    When I am admin
    Then I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-1:
    """json
    {
      "title": "test-widgetfilter-to-update-title-1-updated",
      "query": "{\"test\":\"test-updated\"}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-widgetfilter-to-update-1",
      "title": "test-widgetfilter-to-update-title-1-updated",
      "query": "{\"test\":\"test-updated\"}",
      "author": "root",
      "created": 1611229670
    }
    """

  Scenario: given personal update request should update filter
    When I am test-role-to-widgetfiler-edit
    Then I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-2:
    """json
    {
      "title": "test-widgetfilter-to-update-title-2-updated",
      "query": "{\"test\":\"test-updated\"}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-widgetfilter-to-update-2",
      "title": "test-widgetfilter-to-update-title-2-updated",
      "query": "{\"test\":\"test-updated\"}",
      "author": "test-user-to-widgetfiler-edit",
      "created": 1611229670
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update
    Then the response code should be 403

  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-check-access:
    """json
    {
      "title": "test-widgetfilter-to-check-access-title",
      "query": "{\"test\":\"test\"}"
    }
    """
    Then the response code should be 403

  Scenario: given update request and someone else filter should return error
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-3:
    """json
    {
      "title": "test-widgetfilter-to-update-3-title",
      "query": "{\"test\":\"test\"}"
    }
    """
    Then the response code should be 404

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "query": "Query is missing."
      }
    }
    """

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-not-found:
    """json
    {
      "title": "test-widgetfilter-not-found-title",
      "query": "{\"test\":\"test\"}"
    }
    """
    Then the response code should be 403
