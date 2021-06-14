Feature: Delete a view group
  I need to be able to delete a view group
  Only admin should be able to delete a view group

  Scenario: given delete request should delete view group
    When I am admin
    When I do DELETE /api/v4/view-groups/test-viewgroup-to-delete
    Then the response code should be 204
    When I do GET /api/v4/view-groups/test-viewgroup-to-delete
    Then the response code should be 404

  Scenario: given request to delete linked group should return error
    When I am admin
    When I do DELETE /api/v4/view-groups/test-viewgroup-to-delete-linked-to-view
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "group is linked to view"
    }
    """

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/view-groups/test-viewgroup-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/view-groups/test-viewgroup-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/view-groups/test-viewgroup-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
