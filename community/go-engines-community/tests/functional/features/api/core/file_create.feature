Feature: Create a file
  I need to be able to create a file
  Only admin should be able to create a file

  Scenario: given private create request should return ok
    When I am admin
    When I add form file texts=test.txt
    And I add form file texts=test.txt
    And I do POST /api/v4/file
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
    When I save response fileId1={{ (index .lastResponse 0)._id }}
    When I save response fileId2={{ (index .lastResponse 1)._id }}
    When I do GET /api/v4/file/{{ .fileId1 }}
    Then the response code should be 200
    Then the response body should contain file test.txt
    When I do GET /api/v4/file/{{ .fileId2 }}
    Then the response code should be 200
    Then the response body should contain file test.txt
    When I am unauth
    When I do GET /api/v4/file/{{ .fileId1 }}
    Then the response code should be 401
    When I do GET /api/v4/file/{{ .fileId2 }}
    Then the response code should be 401

  Scenario: given public create request should return ok
    When I am admin
    When I add form file texts=test.txt
    And I add form file texts=test.txt
    And I do POST /api/v4/file?public=true
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
    When I save response fileId1={{ (index .lastResponse 0)._id }}
    When I save response fileId2={{ (index .lastResponse 1)._id }}
    When I am unauth
    When I do GET /api/v4/file/{{ .fileId1 }}
    Then the response code should be 200
    Then the response body should contain file test.txt
    When I do GET /api/v4/file/{{ .fileId2 }}
    Then the response code should be 200
    Then the response body should contain file test.txt

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/file
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "Files are missing."
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/file
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/file
    Then the response code should be 403
