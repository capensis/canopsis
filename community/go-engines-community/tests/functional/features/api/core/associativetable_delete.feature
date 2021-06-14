Feature: Delete a associative table
  I need to be able to delete a associative table
  Only admin should be able to delete a associative table

  Scenario: given delete request should delete associative table
    When I am admin
    When I do DELETE /api/v4/associativetable?name=test-associativetable-to-delete
    Then the response code should be 204
    Then I do GET /api/v4/associativetable?name=test-associativetable-to-delete
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "name": "test-associativetable-to-delete",
      "content": null
    }
    """

  Scenario: given delete not exist request should return ok
    When I am admin
    When I do DELETE /api/v4/associativetable?name=test-associativetable-not-exist
    Then the response code should be 204

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/associativetable
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/associativetable
    Then the response code should be 403
