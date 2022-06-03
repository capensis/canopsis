Feature: delete a PBehavior
  I need to be able to delete a PBehavior
  Only admin should be able to delete a PBehavior

  Scenario: DELETE a PBehavior but unauthorized
    When I do DELETE /api/v4/pbehaviors/test-pbehavior-to-delete
    Then the response code should be 401

  Scenario: DELETE a PBehavior but without permissions
    When I am noperms
    When I do DELETE /api/v4/pbehaviors/test-pbehavior-to-delete
    Then the response code should be 403

  Scenario: DELETE a PBehavior with success
    When I am admin
    When I do DELETE /api/v4/pbehaviors/test-pbehavior-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-delete-1
    Then the response code should be 404

  Scenario: DELETE previous PBehavior, should be not found response
    When I am admin
    When I do DELETE /api/v4/pbehaviors/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
