Feature: Get a broadcast message
  I need to be able to get a broadcast message
  Only admin should be able to get a broadcast message

  Scenario: given search request should return broadcast messages
    When I am admin
    When I do GET /api/v4/broadcast-message?search=test-broadcast-to-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-broadcast-to-get-1",
          "color": "#e75e40",
          "end": 1612296900,
          "message": "broadcast-test-to-get-1",
          "start": 1612139798,
          "created": 1612139798,
          "updated": 1612139798
        },
        {
          "_id": "test-broadcast-to-get-2",
          "color": "rgb(159, 5, 0)",
          "end": 1612296900,
          "message": "broadcast-test-to-get-2",
          "start": 1612139798,
          "created": 1612139798,
          "updated": 1612139798
        },
        {
          "_id": "test-broadcast-to-get-3",
          "color": "#e75e40",
          "end": 1612139798,
          "message": "broadcast-test-to-get-3",
          "start": 1612139798,
          "created": 1612139798,
          "updated": 1612139798
        },
        {
          "_id": "test-broadcast-to-get-4",
          "color": "#e75e40",
          "message": "broadcast-test-to-get-4",
          "created": 1612139798,
          "updated": 1612139798
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  Scenario: given get request should return broadcast message
    When I am admin
    When I do GET /api/v4/broadcast-message/test-broadcast-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-broadcast-to-get-1",
      "color": "#e75e40",
      "end": 1612296900,
      "message": "broadcast-test-to-get-1",
      "start": 1612139798,
      "created": 1612139798,
      "updated": 1612139798
    }
    """

  Scenario: given active request should return broadcast messages
    When I am admin
    When I do GET /api/v4/active-broadcast-message
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "_id": "test-broadcast-to-get-4",
        "color": "#e75e40",
        "message": "broadcast-test-to-get-4",
        "created": 1612139798,
        "updated": 1612139798
      }
    ]
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/broadcast-message
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/broadcast-message
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/broadcast-message/test-broadcast-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/broadcast-message/test-broadcast-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/broadcast-message/test-broadcast-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
