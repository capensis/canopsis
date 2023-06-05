Feature: delete an instruction

  Scenario: DELETE an instruction but unauthorized
    When I do DELETE /api/v4/cat/instructions/test-instruction-to-delete
    Then the response code should be 401

  Scenario: DELETE an instruction but without permissions
    When I am noperms
    When I do DELETE /api/v4/cat/instructions/test-instruction-to-delete
    Then the response code should be 403

  Scenario: DELETE an instruction with success
    When I am admin
    When I do DELETE /api/v4/cat/instructions/test-instruction-to-delete
    Then the response code should be 204

  Scenario: DELETE an instruction with not found response
    When I am admin
    When I do DELETE /api/v4/cat/instructions/test-instruction-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
