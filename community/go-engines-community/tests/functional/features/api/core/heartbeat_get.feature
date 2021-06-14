Feature: Get a heartbeat
  I need to be able to get a heartbeat
  Only admin should be able to get a heartbeat

  Scenario: given search request should return heartbeats
    When I am admin
    When I do GET /api/v4/heartbeats?search=test-heartbeat-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-heartbeat-to-get-1",
          "name": "test-heartbeat-to-get-1-name",
          "description": "test-heartbeat-to-get-1-description",
          "author": "root",
          "pattern": {
            "name": "test-heartbeat-to-get-1-resource"
          },
          "expected_interval": "1000h",
          "output": "test-heartbeat-to-get-1-output",
          "created": 1592215337,
          "updated": 1592215337
        },
        {
          "_id": "test-heartbeat-to-get-2",
          "name": "test-heartbeat-to-get-2-name",
          "description": "test-heartbeat-to-get-2-description",
          "author": "root",
          "pattern": {
            "name": "test-heartbeat-to-get-2-resource"
          },
          "expected_interval": "1000h",
          "output": "test-heartbeat-to-get-2-output",
          "created": 1592215347,
          "updated": 1592215347
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

  Scenario: given get request should return heartbeat
    When I am admin
    When I do GET /api/v4/heartbeats/test-heartbeat-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-heartbeat-to-get-1",
      "name": "test-heartbeat-to-get-1-name",
      "description": "test-heartbeat-to-get-1-description",
      "author": "root",
      "pattern": {
        "name": "test-heartbeat-to-get-1-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-get-1-output",
      "created": 1592215337,
      "updated": 1592215337
    }
    """

  Scenario: given sort request should return sorted heartbeats
    When I am admin
    When I do GET /api/v4/heartbeats?search=test-heartbeat-to-get&sort=desc&sort_by=updated
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-heartbeat-to-get-2"
        },
        {
          "_id": "test-heartbeat-to-get-1"
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/heartbeats
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/heartbeats
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/heartbeats/test-heartbeat-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/heartbeats/test-heartbeat-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/heartbeats/test-heartbeat-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """