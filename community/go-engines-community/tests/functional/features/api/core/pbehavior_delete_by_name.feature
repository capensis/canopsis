Feature: Delete a pbehavior
  I need to be able to delete a pbehavior
  Only admin should be able to delete a pbehavior

  Scenario: given delete request should delete pbehavior
    When I am admin
    When I do DELETE /api/v4/pbehaviors?name=test-pbehavior-to-delete-2-name
    Then the response code should be 204
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-delete-3
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/pbehaviors?name=test-pbehavior-to-delete-name
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/pbehaviors?name=test-pbehavior-to-delete-name
    Then the response code should be 403

  Scenario: given delete request with not exist name should return not found error
    When I am admin
    When I do DELETE /api/v4/pbehaviors?name=test-not-found
    Then the response code should be 404

  Scenario: given invalid delete request should return error
    When I am admin
    When I do DELETE /api/v4/pbehaviors
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name is missing."
      }
    }
    """
