Feature: Create a view group
  I need to be able to create a view group
  Only admin should be able to create a view group

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/view-groups:
    """
    {
      "title": "test-viewgroup-to-create-1-title"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "title": "test-viewgroup-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/view-groups:
    """
    {
      "title": "test-viewgroup-to-create-2-title"
    }
    """
    When I do GET /api/v4/view-groups/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "title": "test-viewgroup-to-create-2-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/view-groups
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/view-groups
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/view-groups:
    """
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "title": "Title is missing."
      }
    }
    """

  Scenario: given create request with already exists title should return error
    When I am admin
    When I do POST /api/v4/view-groups:
    """
    {
      "title": "test-viewgroup-to-check-unique-title-title"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "title": "Title already exists."
      }
    }
    """
