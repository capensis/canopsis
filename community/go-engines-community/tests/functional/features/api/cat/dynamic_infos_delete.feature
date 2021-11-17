Feature: Delete a dynamic infos
  I need to be able to delete a dynamic infos
  Only admin should be able to delete a dynamic infos

  Scenario: given delete request should delete dynamic infos
    When I am admin
    When I do DELETE /api/v4/cat/dynamic-infos/test-dynamic-infos-to-delete
    Then the response code should be 204
    Then I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/cat/dynamic-infos/test-dynamic-infos-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/cat/dynamic-infos/test-dynamic-infos-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
