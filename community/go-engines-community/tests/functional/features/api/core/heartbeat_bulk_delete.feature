Feature: Bulk delete heartbeats
  I need to be able to delete multiple heartbeats
  Only admin should be able to delete multiple heartbeats

  Scenario: given bulk delete request should delete heartbeat
    When I am admin
    When I do DELETE /api/v4/bulk/heartbeats?ids[]=test-heartbeat-to-bulk-delete-1&ids[]=test-heartbeat-to-bulk-delete-2
    Then the response code should be 204
    Then I do GET /api/v4/heartbeats/test-heartbeat-to-bulk-delete-1
    Then the response code should be 404
    Then I do GET /api/v4/heartbeats/test-heartbeat-to-bulk-delete-2
    Then the response code should be 404

  Scenario: given bulk delete request with not exist ids should return not found error
    When I am admin
    When I do DELETE /api/v4/bulk/heartbeats?ids[]=test-heartbeat-not-found-1&ids[]=test-heartbeat-not-found-2
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "not found test-heartbeat-not-found-1,test-heartbeat-not-found-2"
    }
    """

  Scenario: given bulk delete request with one exits id and one not exist
    should return not found error
    When I am admin
    When I do DELETE /api/v4/bulk/heartbeats?ids[]=test-heartbeat-to-bulk-delete-3&ids[]=test-heartbeat-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "not found test-heartbeat-not-found"
    }
    """

  Scenario: given bulk delete request and no auth user should not allow access
    When I do DELETE /api/v4/bulk/heartbeats?ids[]=test-heartbeat-to-bulk-delete-1
    Then the response code should be 401

  Scenario: given bulk delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/heartbeats?ids[]=test-heartbeat-to-bulk-delete-1
    Then the response code should be 403
