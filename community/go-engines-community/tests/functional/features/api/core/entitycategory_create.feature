Feature: Create a entity category
  I need to be able to create a entity category
  Only admin should be able to create a entity category

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/entity-categories:
    """json
    {
      "name": "test-category-to-create-1-name"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-category-to-create-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/entity-categories:
    """json
    {
      "name": "test-category-to-create-2-name"
    }
    """
    When I do GET /api/v4/entity-categories/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-category-to-create-2-name",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/entity-categories
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/entity-categories
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/entity-categories:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/entity-categories:
    """json
    {
      "name": "test-category-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """
