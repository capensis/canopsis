Feature: Update a entity category
  I need to be able to update a entity category
  Only admin should be able to update a entity category

  Scenario: given update request should update entity category
    When I am admin
    Then I do PUT /api/v4/entity-categories/test-category-to-update:
    """
    {
      "name": "test-category-to-update-name-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-category-to-update",
      "name": "test-category-to-update-name-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1592215337
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/entity-categories/test-category-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/entity-categories/test-category-to-update
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/entity-categories/test-category-to-update:
    """
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name is missing."
      }
    }
    """

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/entity-categories/test-category-not-found:
    """
    {
      "name": "test-category-not-found-name"
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
