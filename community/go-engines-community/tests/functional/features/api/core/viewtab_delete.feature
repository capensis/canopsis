Feature: Delete a view tab
  I need to be able to delete a view tab
  Only admin should be able to delete a view tab

  Scenario: given delete request should delete tab
    When I am admin
    When I do DELETE /api/v4/view-tabs/test-tab-to-delete
    Then the response code should be 204
    When I do GET /api/v4/view-tabs/test-tab-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/view-tabs/test-tab-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/view-tabs/test-tab-to-delete
    Then the response code should be 403

  Scenario: given delete request and auth user without view permission should not allow access
    When I am admin
    When I do DELETE /api/v4/view-tabs/test-tab-to-check-access
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not allow access error
    When I am admin
    When I do DELETE /api/v4/view-tabs/test-tab-not-found
    Then the response code should be 404

  Scenario: given delete request linked tab should return error
    When I am admin
    When I do DELETE /api/v4/view-tabs/test-tab-to-delete-linked
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "view tab is linked to playlist"
    }
    """
