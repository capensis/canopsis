Feature: Delete a file
  I need to be able to delete a file
  Only admin should be able to delete a file

  Scenario: given delete request should return ok
    When I am admin
    When I add form file files=test.txt
    And I add form file files=test.txt
    And I do POST /api/v4/file
    Then the response code should be 200
    When I save response fileId1={{ (index .lastResponse 0)._id }}
    When I save response fileId2={{ (index .lastResponse 1)._id }}
    When I do DELETE /api/v4/file/{{ .fileId1 }}
    Then the response code should be 204
    When I do GET /api/v4/file/{{ .fileId1 }}
    Then the response code should be 404
    When I do DELETE /api/v4/file/{{ .fileId2 }}
    Then the response code should be 204
    When I do GET /api/v4/file/{{ .fileId2 }}
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/file/test-file-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/file/test-file-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/file/test-file-not-exists
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
