Feature: Delete a scenario
  I need to be able to bulk delete scenarios
  Only admin should be able to bulk delete scenarios

  Scenario: given delete request and no auth scenario should not allow access
    When I do DELETE /api/v4/bulk/scenarios?ids[]=test-scenario-to-bulk-delete-1&ids[]=test-scenario-to-bulk-delete-2
    Then the response code should be 401

  Scenario: given delete request and auth scenario by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/scenarios?ids[]=test-scenario-to-bulk-delete-1&ids[]=test-scenario-to-bulk-delete-2
    Then the response code should be 403

  Scenario: given delete request should delete scenario
    When I am admin
    When I do GET /api/v4/scenarios/test-scenario-to-bulk-delete-1
    Then the response code should be 200
    When I do GET /api/v4/scenarios/test-scenario-to-bulk-delete-2
    Then the response code should be 200
    When I do DELETE /api/v4/bulk/scenarios?ids[]=test-scenario-to-bulk-delete-1&ids[]=test-scenario-to-bulk-delete-2
    Then the response code should be 204
    When I do GET /api/v4/scenarios/test-scenario-to-bulk-delete-1
    Then the response code should be 404
    When I do GET /api/v4/scenarios/test-scenario-to-bulk-delete-2
    Then the response code should be 404

  Scenario: given delete request with empty ids should return error
    When I do DELETE /api/v4/bulk/scenarios
    """json
    {
      "errors": {
        "ids": "IDs is missing."
      }
    }
    """
