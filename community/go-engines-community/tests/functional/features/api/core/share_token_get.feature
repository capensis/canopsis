Feature: Get a share token
  I need to be able to get a share token
  Only admin should be able to get a share token

  Scenario: given search request should return tokens
    When I am admin
    When I do GET /api/v4/share-tokens?search=test-share-token-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-share-token-to-get-1",
          "value": "**********fk950",
          "description": "test-share-token-to-get-1-description",
          "user": {
            "_id": "root",
            "name": "root"
          },
          "role": {
            "_id": "admin",
            "name": "admin"
          },
          "created": 1619083733,
          "accessed": 1619083733,
          "expired": null
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/share-tokens
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/share-tokens
    Then the response code should be 403
