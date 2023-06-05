Feature: Delete a role
  I need to be able to delete a role
  Only admin should be able to delete a role

  Scenario: given delete request should delete role
    When I am admin
    When I do DELETE /api/v4/roles/test-role-to-delete
    Then the response code should be 204
    When I do GET /api/v4/roles/test-role-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/roles/test-role-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/roles/test-role-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/roles/test-role-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given delete request with linked to role id should return validation error
    When I am admin
    When I do DELETE /api/v4/roles/test-role-to-delete-linked-to-user
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "role is linked to user"
    }
    """

  Scenario: given delete request for admin should return validation error
    When I am admin
    When I do DELETE /api/v4/roles/admin
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "admin cannot be deleted"
    }
    """
