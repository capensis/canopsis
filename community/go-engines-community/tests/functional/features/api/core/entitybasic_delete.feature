Feature: Delete a entity basic
  I need to be able to delete a entity basic
  Only admin should be able to delete a entity basic

  Scenario: given delete request should delete entity basic
    When I am admin
    When I do DELETE /api/v4/entitybasics?_id=test-entitybasic-to-delete-resource/test-entitybasic-to-delete-component
    Then the response code should be 204
    When I do GET /api/v4/entitybasics?_id=test-entitybasic-to-delete-resource/test-entitybasic-to-delete-component
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/entitybasics?_id=test-entitybasic-not-exist
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/entitybasics?_id=test-entitybasic-not-exist
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/entitybasics?_id=test-entitybasic-not-exist
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given request to delete linked entitybasic should return validation error
    When I am admin
    When I do DELETE /api/v4/entitybasics?_id=test-entitybasic-to-delete-linked-to-alarm-resource/test-entitybasic-to-delete-linked-to-alarm-component
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "entity is linked to alarm"
    }
    """
