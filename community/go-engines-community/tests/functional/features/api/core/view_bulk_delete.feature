Feature: Bulk delete views
  I need to be able to delete multiple views
  Only admin should be able to delete multiple views

  Scenario: given bulk delete request should delete view
    When I am admin
    When I do DELETE /api/v4/bulk/views?ids[]=test-view-to-bulk-delete-1&ids[]=test-view-to-bulk-delete-2
    Then the response code should be 204
    Then I do GET /api/v4/views/test-view-to-bulk-delete-1
    Then the response code should be 403
    When I do GET /api/v4/permissions?search=test-view-to-bulk-delete-1
    Then the response code should be 200
    Then the response body should be:
    """
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
    Then I do GET /api/v4/views/test-view-to-bulk-delete-2
    Then the response code should be 403
    When I do GET /api/v4/permissions?search=test-view-to-bulk-delete-2
    Then the response code should be 200
    Then the response body should be:
    """
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

  Scenario: given bulk delete request with not exist view should return error
    When I am admin
    When I do DELETE /api/v4/bulk/views?ids[]=test-view-to-bulk-delete-3&ids[]=test-view-not-found
    Then the response code should be 403

  Scenario: given bulk delete request and no auth user should not allow access
    When I do DELETE /api/v4/bulk/views
    Then the response code should be 401

  Scenario: given bulk delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/views
    Then the response code should be 403

  Scenario: given bulk delete request and auth user without view permission should not allow access
    When I am admin
    When I do DELETE /api/v4/bulk/views?ids[]=test-view-to-bulk-delete-3&ids[]=test-view-to-check-access
    Then the response code should be 403
