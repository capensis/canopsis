Feature: Delete a view
  I need to be able to delete a view
  Only admin should be able to delete a view

  Scenario: given delete request should delete view
    When I am admin
    When I do DELETE /api/v4/views/test-view-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/views/test-view-to-delete-1
    Then the response code should be 403
    When I do GET /api/v4/permissions?search=test-view-to-delete-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/views/test-view-to-delete-1
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/views/test-view-to-delete-1
    Then the response code should be 403

  Scenario: given delete request and auth user without view permission should not allow access
    When I am admin
    When I do DELETE /api/v4/views/test-view-to-check-access
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not allow access error
    When I am admin
    When I do DELETE /api/v4/views/test-view-not-found
    Then the response code should be 403

  Scenario: given delete request should delete tabs and widgets
    When I am admin
    When I do DELETE /api/v4/views/test-view-to-delete-2
    Then the response code should be 204
    When I do GET /api/v4/views/test-view-to-delete-2
    Then the response code should be 403
    When I do GET /api/v4/view-tabs/test-tab-to-view-delete
    Then the response code should be 404
    When I do GET /api/v4/widgets/test-widget-to-view-delete
    Then the response code should be 404
    When I do GET /api/v4/permissions?search=test-view-to-delete-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
