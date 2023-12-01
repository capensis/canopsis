Feature: Create an icon
  I need to be able to create an icon
  Only admin should be able to create an icon

  Scenario: given create request should return ok
    When I am admin
    When I add form field title=test-icon-to-create-1-title
    And I add form file file=test.svg
    And I do POST /api/v4/icons
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-icon-to-create-1-title"
    }
    """
    When I do GET /api/v4/icons/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain file test.svg

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/icons
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "file": "File is missing."
      }
    }
    """
    When I add form file file=test.txt
    And I add form field title=test-icon-to-create
    And I do POST /api/v4/icons
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "file": "Invalid mime type: text/plain; charset=utf-8"
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/icons
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/icons
    Then the response code should be 403
