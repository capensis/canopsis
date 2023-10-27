Feature: Delete a widget
  I need to be able to delete a widget
  Only admin should be able to delete a widget

  @concurrent
  Scenario: given delete request should delete widget
    When I am admin
    When I do DELETE /api/v4/widgets/test-widget-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/widgets/test-widget-to-delete-1
    Then the response code should be 404

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/widgets/test-widget-to-delete-1
    Then the response code should be 401

  @concurrent
  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/widgets/test-widget-to-delete-1
    Then the response code should be 403

  @concurrent
  Scenario: given delete request and auth user without view permission should not allow access
    When I am admin
    When I do DELETE /api/v4/widgets/test-widget-to-check-access
    Then the response code should be 403

  @concurrent
  Scenario: given delete request with not exist id should return not allow access error
    When I am admin
    When I do DELETE /api/v4/widgets/test-widget-not-found
    Then the response code should be 404

  @concurrent
  Scenario: given delete request with owned private widget should be ok
    When I am admin
    When I do DELETE /api/v4/widgets/test-private-widget-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/view-tabs/test-private-widget-to-delete-1
    Then the response code should be 404

  @concurrent
  Scenario: given delete request with not owned private widget should return not allow access error
    When I am admin
    When I do DELETE /api/v4/widgets/test-private-widget-to-delete-2
    Then the response code should be 403

  @concurrent
  Scenario: given delete request with owned private with api_private_view_groups
    but without api_view permissions should be ok
    When I am test-role-to-private-views-without-view-perm
    When I do DELETE /api/v4/widgets/test-private-widget-to-delete-3
    Then the response code should be 204
    When I do GET /api/v4/widgets/test-private-widget-to-delete-3
    Then the response code should be 404

  @concurrent
  Scenario: given delete request with owned private with api_private_view_groups
    but without api_view permissions should be ok
    When I am test-role-to-private-views-without-view-perm
    When I do DELETE /api/v4/widgets/test-widget-to-delete-2
    Then the response code should be 403
