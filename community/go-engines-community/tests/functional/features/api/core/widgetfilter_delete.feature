Feature: Delete a widget filter
  I need to be able to delete a widget filter
  Only admin should be able to delete a widget filter

  Scenario: given delete public filter request should return ok
    When I am test-role-to-widget-filter-edit-2
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-delete-1
    Then the response code should be 403

  Scenario: given delete private filter request should return ok
    When I am test-role-to-widget-filter-edit-1
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-to-delete-2
    Then the response code should be 204
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-delete-2
    Then the response code should be 403

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 403

  Scenario: given delete private filter request and auth user without view permission should not allow access
    When I am admin
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-to-delete-3
    Then the response code should be 403

  Scenario: given delete private filter request and another user should return not found
    When I am admin
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-to-delete-4
    Then the response code should be 404

  Scenario: given delete not exist filter request should return error
    When I am admin
    When I do DELETE /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 404
