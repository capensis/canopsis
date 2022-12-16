Feature: Create a view tab
  I need to be able to create a view tab
  Only admin should be able to create a view tab

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-1-title",
      "view": "test-view-to-tab-edit"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """
    When I do GET /api/v4/view-tabs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/view-tabs
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/view-tabs
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "view": "View is missing."
      }
    }
    """

  Scenario: given create request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-2-title",
      "view": "test-view-to-tab-check-access"
    }
    """
    Then the response code should be 403

  Scenario: given create request with not exist view should return not allow access error
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-2-title",
      "view": "test-view-not-found"
    }
    """
    Then the response code should be 403
