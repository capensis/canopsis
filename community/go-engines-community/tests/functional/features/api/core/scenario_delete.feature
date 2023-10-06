Feature: Delete a scenario
  I need to be able to delete a scenario
  Only admin should be able to delete a scenario

  @concurrent
  Scenario: given delete request should delete exception
    When I am admin
    When I do DELETE /api/v4/scenarios/test-scenario-to-delete
    Then the response code should be 204
    When I do GET /api/v4/scenarios/test-scenario-to-delete
    Then the response code should be 404

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/scenarios/notexist
    Then the response code should be 401

  @concurrent
  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/scenarios/notexist
    Then the response code should be 403

  @concurrent
  Scenario: given invalid delete request should return error
    When I am admin
    When I do DELETE /api/v4/scenarios/notexist
    Then the response code should be 404
