Feature: Delete a private view group
  I need to be able to delete a private view group
  Only user with permission should be able to delete a private view group

  @concurrent
  Scenario: given delete request should delete view group
    When I am admin
    When I do DELETE /api/v4/cat/private-view-groups/test-private-viewgroup-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/view-groups/test-private-viewgroup-to-delete-1
    Then the response code should be 404

  @concurrent
  Scenario: given request to delete linked group should return error
    When I am admin
    When I do DELETE /api/v4/cat/private-view-groups/test-private-viewgroup-to-delete-linked-to-view
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "group is linked to view"
    }
    """

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/cat/private-view-groups/test-private-viewgroup-to-delete-1
    Then the response code should be 401

  @concurrent
  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/cat/private-view-groups/test-viewgroup-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given delete request for not owned private view should not allow access
    When I am admin
    When I do DELETE /api/v4/cat/private-view-groups/test-private-viewgroup-to-delete-2
    Then the response code should be 403
