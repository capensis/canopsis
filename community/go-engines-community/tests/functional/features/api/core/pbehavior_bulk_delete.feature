Feature: Bulk delete pbehaviors
  I need to be able to delete multiple pbehaviors
  Only admin should be able to delete multiple pbehaviors

  Scenario: given bulk delete request should delete pbehavior
    When I am admin
    When I do DELETE /api/v4/bulk/pbehaviors?ids[]=test-pbehavior-to-bulk-delete-1&ids[]=test-pbehavior-to-bulk-delete-2
    Then the response code should be 204
    Then I do GET /api/v4/pbehaviors/test-pbehavior-to-bulk-delete-1
    Then the response code should be 404
    Then I do GET /api/v4/pbehaviors/test-pbehavior-to-bulk-delete-2
    Then the response code should be 404

  Scenario: given bulk delete request with not exist pbehavior should return not found error
    When I am admin
    When I do DELETE /api/v4/bulk/pbehaviors?ids[]=test-pbehavior-to-bulk-delete-3&ids[]=test-pbehavior-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given bulk delete request and no auth user should not allow access
    When I do DELETE /api/v4/bulk/pbehaviors
    Then the response code should be 401

  Scenario: given bulk delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/pbehaviors
    Then the response code should be 403
