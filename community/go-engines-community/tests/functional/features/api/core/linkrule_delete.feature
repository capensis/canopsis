Feature: Delete a link rule
  I need to be able to delete a link rule
  Only admin should be able to delete a link rule

  @concurrent
  Scenario: given delete request should delete link rule
    When I am admin
    When I do DELETE /api/v4/link-rules/test-link-rule-to-delete
    Then the response code should be 204
    Then I do GET /api/v4/link-rules/test-link-rule-to-delete
    Then the response code should be 404

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/link-rules/test-link-rule-to-delete
    Then the response code should be 401

  @concurrent
  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/link-rules/test-link-rule-to-delete
    Then the response code should be 403

  @concurrent
  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/link-rules/test-link-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
