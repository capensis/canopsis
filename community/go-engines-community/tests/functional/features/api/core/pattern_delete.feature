Feature: Delete a saved pattern
  I need to be able to delete a saved pattern
  Only admin should be able to delete a saved pattern

  Scenario: given delete request should return ok
    When I am noperms
    When I do DELETE /api/v4/patterns/test-pattern-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/patterns/test-pattern-to-delete-1
    Then the response code should be 404

  Scenario: given delete request and another user should return not found
    When I am admin
    When I do DELETE /api/v4/patterns/test-pattern-to-delete-2
    Then the response code should be 404

  Scenario: given delete corporate pattern request and another user should return ok
    When I am admin
    When I do DELETE /api/v4/patterns/test-pattern-to-delete-3
    Then the response code should be 204

  Scenario: given delete corporate pattern request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/patterns/test-pattern-to-delete-4
    Then the response code should be 403

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/patterns/test-pattern-notexist
    Then the response code should be 401

  Scenario: given delete request with not exist id should return not found error
    When I am noperms
    When I do DELETE /api/v4/patterns/test-pattern-notexist
    Then the response code should be 404
