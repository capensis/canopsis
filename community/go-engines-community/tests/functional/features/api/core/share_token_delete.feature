Feature: Delete a share token
  I need to be able to delete a share token
  Only admin should be able to delete a share token

  Scenario: given delete request should return ok
    When I am admin
    When I do DELETE /api/v4/share-tokens/test-share-token-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/share-tokens?search=test-share-token-to-delete-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/share-tokens/test-share-token-notexist
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/share-tokens/test-share-token-notexist
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not found error
    When I am admin
    When I do DELETE /api/v4/share-tokens/test-share-token-notexist
    Then the response code should be 404
