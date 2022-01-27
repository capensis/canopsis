Feature: Create an filter
  I need to be able to create a filter
  Only admin should be able to create a filter

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-to-create-1-name",
      "entity_patterns": [
        {
          "name": "test-filter-to-create-1-resource"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-filter-to-create-1-name",
      "entity_patterns": [
        {
          "name": "test-filter-to-create-1-resource"
        }
      ]
    }
    """
    When I do GET /api/v4/cat/filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-filter-to-create-1-name",
      "entity_patterns": [
        {
          "name": "test-filter-to-create-1-resource"
        }
      ]
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_patterns": "EntityPatterns is missing.",
        "name": "Name is missing."
      }
    }
    """

  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "entity_patterns": [
        {
          "component_name": "ram"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_patterns":"Invalid entity pattern list."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/filters
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/filters
    Then the response code should be 403

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """
