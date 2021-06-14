Feature: Delete a idle rule
  I need to be able to delete a idle rule
  Only admin should be able to delete a idle rule

  Scenario: given delete request should delete exception
    When I am admin
    When I do DELETE /api/v4/idle-rules/test-idle-rule-to-delete
    Then the response code should be 204
    When I do GET /api/v4/idle-rules/test-idle-rule-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/idle-rules/notexist
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/idle-rules/notexist
    Then the response code should be 403

  Scenario: given invalid delete request should return error
    When I am admin
    When I do DELETE /api/v4/idle-rules/notexist
    Then the response code should be 404
