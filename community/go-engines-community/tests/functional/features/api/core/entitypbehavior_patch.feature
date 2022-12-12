Feature: update a pbehavior
  I need to be able to patch a pbehavior field individually
  Only admin should be able to patch a pbehavior

  Scenario: given update entity pbehavior request should return error
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-entity-patch-1:
    """json
    {
      "name": "test-pbehavior-to-entity-patch-1-name"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "Cannot update a pbehavior with origin."
      }
    }
    """
