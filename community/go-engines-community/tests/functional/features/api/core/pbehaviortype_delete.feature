Feature: PBehavior Type delete

  Scenario: DELETE as unauthorized
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete
    Then the response code should be 401

  Scenario: DELETE without permissions
    When I am noperms
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete
    Then the response code should be 403

  Scenario: DELETE with success
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete-1
    Then the response code should be 204

  Scenario: DELETE with success
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete-2
    When I do GET /api/v4/pbehavior-types/test-type-to-delete
    Then the response code should be 404

  Scenario: Given linked type Should return error
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete-linked-to-pbh-1
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "type is linked to pbehavior"
    }
    """

  Scenario: Given linked type Should return error
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete-linked-to-pbh-2
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "type is linked to pbehavior"
    }
    """

  Scenario: Given linked type Should return error
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete-linked-to-exception
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "type is linked to exception"
    }
    """

  Scenario: Given linked type Should return error
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete-linked-to-action
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "type is linked to action"
    }
    """

  Scenario: Given default type Should return error
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-default-inactive-type
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "type is default"
    }
    """

  Scenario: DELETE previous PBehavior Type, should be not found response
    When I am admin
    When I do DELETE /api/v4/pbehavior-types/test-type-to-delete
    Then the response code should be 404
    Then the response body should be:
    """
    {
        "error": "Not found"
    }
    """
