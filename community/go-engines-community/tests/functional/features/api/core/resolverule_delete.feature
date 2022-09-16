Feature: Delete a resolve rule
  I need to be able to delete a resolve rule
  Only admin should be able to delete a resolve rule

  Scenario: given delete request should delete resolve rule
    When I am admin
    When I do DELETE /api/v4/resolve-rules/test-resolve-rule-to-delete
    Then the response code should be 204
    Then I do GET /api/v4/resolve-rules/test-resolve-rule-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/resolve-rules/test-resolve-rule-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/resolve-rules/test-resolve-rule-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/resolve-rules/test-resolve-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given delete default rule request should return error
    When I am admin
    When I do DELETE /api/v4/resolve-rules/default_rule
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "cannot delete the default rule"
    }
    """
