Feature: get instruction comments
  I need to be able to get instruction comments
  Only admin should be able to get instruction comments
  
  Scenario: given get request should return instruction comments
    When I am admin
    When I do GET /api/v4/cat/instruction-comments/test-instruction-to-comments-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-rating-to-comments-get-1",
          "comment": "test-rating-to-comments-get-1-comment",
          "rating": 4.5,
          "created": 1596550518,
          "user": {
            "_id": "test-user-author-1-id",
            "name": "test-user-author-1-username"
          }
        },
        {
          "_id": "test-rating-to-comments-get-2",
          "comment": "test-rating-to-comments-get-2-comment",
          "rating": 2,
          "created": 1596550518,
          "user": {
            "_id": "test-user-author-1-id",
            "name": "test-user-author-1-username"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get request should return empty instruction comments
    When I am admin
    When I do GET /api/v4/cat/instruction-comments/test-instruction-to-comments-get-2
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

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/instruction-comments/notexist
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/instruction-comments/notexist
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/instruction-comments/notexist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
