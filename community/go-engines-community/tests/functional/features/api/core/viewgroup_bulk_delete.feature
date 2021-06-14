Feature: Bulk delete view groups
  I need to be able to delete multiple view groups
  Only admin should be able to delete multiple view groups

  Scenario: given bulk delete request should delete view group
    When I am admin
    When I do DELETE /api/v4/bulk/view-groups?ids[]=test-viewgroup-to-bulk-delete-1&ids[]=test-viewgroup-to-bulk-delete-2
    Then the response code should be 204
    Then I do GET /api/v4/view-groups/test-viewgroup-to-bulk-delete-1
    Then the response code should be 404
    Then I do GET /api/v4/view-groups/test-viewgroup-to-bulk-delete-2
    Then the response code should be 404

  Scenario: given bulk delete request with not exist group should return not found error
    When I am admin
    When I do DELETE /api/v4/bulk/view-groups?ids[]=test-viewgroup-to-bulk-delete-3&ids[]=test-viewgroup-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given bulk request to delete linked group should return error
    When I am admin
    When I do DELETE /api/v4/bulk/view-groups?ids[]=test-viewgroup-to-bulk-delete-3&ids[]=test-viewgroup-to-delete-linked-to-view
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "group is linked to view"
    }
    """

  Scenario: given bulk delete request and no auth user should not allow access
    When I do DELETE /api/v4/bulk/view-groups
    Then the response code should be 401

  Scenario: given bulk delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/view-groups
    Then the response code should be 403
