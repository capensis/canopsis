Feature: Update an icon
  I need to be able to update an icon
  Only admin should be able to update an icon

  Scenario: given update request should return ok
    When I am admin
    When I add form field title=test-icon-to-update-1-title
    And I add form file file=test.svg
    And I do POST /api/v4/icons
    Then the response code should be 201
    When I add form field title=test-icon-to-update-1-title-updated
    And I add form file file=test2.svg
    And I do PUT /api/v4/icons/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-icon-to-update-1-title-updated"
    }
    """
    When I do GET /api/v4/icons/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain file test2.svg

  Scenario: given update request with missing fields should return bad request
    When I am admin
    When I do PUT /api/v4/icons/test-icon-to-update
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
    And I add form field title=test-icon-to-update
    And I do PUT /api/v4/icons/test-icon-to-update
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "file": "Invalid mime type: text/plain; charset=utf-8"
      }
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/icons/test-icon-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/icons/test-icon-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I add form field title=test-icon-to-update-title
    And I add form file file=test.svg
    And I do PUT /api/v4/icons/test-icon-not-exists
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
