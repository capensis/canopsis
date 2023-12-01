Feature: Delete an icon
  I need to be able to delete an icon
  Only admin should be able to delete an icon

  Scenario: given delete request should return ok
    When I am admin
    When I add form field title=test-icon-to-delete-1-title
    And I add form file file=test.svg
    And I do POST /api/v4/icons
    Then the response code should be 201
    When I save response iconId={{ .lastResponse._id }}
    When I do DELETE /api/v4/icons/{{ .iconId }}
    Then the response code should be 204
    When I do GET /api/v4/icons/{{ .iconId }}
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/icons/test-icon-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/icons/test-icon-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    And I do DELETE /api/v4/icons/test-icon-not-exists
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
