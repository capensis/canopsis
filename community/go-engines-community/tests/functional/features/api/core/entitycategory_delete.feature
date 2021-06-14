Feature: Delete a entity category
  I need to be able to delete a entity category
  Only admin should be able to delete a entity category

  Scenario: given delete request should delete entity category
    When I am admin
    When I do DELETE /api/v4/entity-categories/test-category-to-delete
    Then the response code should be 204
    When I do GET /api/v4/entity-categories/test-category-to-delete
    Then the response code should be 404

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/entity-categories/test-category-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/entity-categories/test-category-to-delete
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/entity-categories/test-category-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given request to delete linked category should return validation error
    When I am admin
    When I do DELETE /api/v4/entity-categories/test-category-to-delete-linked-to-entity
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "category is linked to entity"
    }
    """
