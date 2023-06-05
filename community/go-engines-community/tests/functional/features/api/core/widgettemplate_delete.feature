Feature: Delete a widget template
  I need to be able to delete a widget template
  Only admin should be able to delete a widget template

  Scenario: given delete request should return ok
    When I am admin
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/widget-templates/test-widgettemplate-to-delete-1
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 404
