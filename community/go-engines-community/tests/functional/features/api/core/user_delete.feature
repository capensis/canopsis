Feature: Delete a user
  I need to be able to delete a user
  Only admin should be able to delete a user

  Scenario: given delete request should delete user
    When I am admin
    When I do DELETE /api/v4/users/test-user-to-delete
    Then the response code should be 204
    When I do GET /api/v4/users/test-user-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/users/test-user-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/users/test-user-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/users/test-user-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """


