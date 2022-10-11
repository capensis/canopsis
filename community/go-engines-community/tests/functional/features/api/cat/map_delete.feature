Feature: Delete a map
  I need to be able to delete a map
  Only admin should be able to delete a map

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/cat/maps/test-map-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/cat/maps/test-map-to-delete
    Then the response code should be 403

  Scenario: given delete request should return ok
    When I am admin
    When I do DELETE /api/v4/cat/maps/test-map-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/cat/maps/test-map-to-delete-1
    Then the response code should be 404

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/cat/maps/test-map-not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given delete linked map request should return error
    When I am admin
    When I do DELETE /api/v4/cat/maps/test-map-to-delete-2
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "map is linked with widget"
    }
    """
