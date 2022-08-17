Feature: Delete a widget
  I need to be able to delete a widget
  Only admin should be able to delete a widget

  Scenario: given delete request should delete widget
    When I am admin
    When I do DELETE /api/v4/widgets/test-widget-to-delete
    Then the response code should be 204
    When I do GET /api/v4/widgets/test-widget-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/widgets/test-widget-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/widgets/test-widget-to-delete
    Then the response code should be 403

  Scenario: given delete request and auth user without view permission should not allow access
    When I am admin
    When I do DELETE /api/v4/widgets/test-widget-to-check-access
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not allow access error
    When I am admin
    When I do DELETE /api/v4/widgets/test-widget-not-found
    Then the response code should be 403
