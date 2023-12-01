Feature: Find a file
  I need to be able to find a file
  Only admin should be able to find a file

  Scenario: given find private request should return ok
    When I am admin
    When I add form file files=test.txt
    And I add form file files=test.txt
    And I do POST /api/v4/file
    Then the response code should be 200
    When I save response fileId1={{ (index .lastResponse 0)._id }}
    When I save response fileId2={{ (index .lastResponse 1)._id }}
    When I do GET /api/v4/file?id={{ .fileId1 }}&id={{ .fileId2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "filename": "test.txt"
      },
      {
        "filename": "test.txt"
      }
    ]
    """
    When I am unauth
    When I do GET /api/v4/file?id={{ .fileId1 }}&id={{ .fileId2 }}
    Then the response code should be 401

  Scenario: given find public request should return ok
    When I am admin
    When I add form file files=test.txt
    And I add form file files=test.txt
    And I do POST /api/v4/file?public=true
    Then the response code should be 200
    When I save response fileId1={{ (index .lastResponse 0)._id }}
    When I save response fileId2={{ (index .lastResponse 1)._id }}
    When I am unauth
    When I do GET /api/v4/file?id={{ .fileId1 }}&id={{ .fileId2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "filename": "test.txt"
      },
      {
        "filename": "test.txt"
      }
    ]
    """

  Scenario: given find request with missing fields should return bad request
    When I am admin
    When I do GET /api/v4/file
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "id": "ID is missing."
      }
    }
    """

  Scenario: given find request with not exist id should return not found error
    When I do GET /api/v4/file?id=test-file-not-exists
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
