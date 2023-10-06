Feature: Get event filters
  I need to be able to get event filters
  Only admin should be able to get event filters

  Scenario: given read request should return update counters
    When I am admin
    When I do GET /api/v4/eventfilter/test-eventfilter-to-failure-read-1/failures
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-failure-to-read-2",
          "unread": true
        },
        {
          "_id": "test-eventfilter-failure-to-read-1",
          "unread": true
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
    When I do GET /api/v4/eventfilter/rules?search=test-eventfilter-to-failure-read-1&with_counts=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-to-failure-read-1",
          "failures_count": 2,
          "unread_failures_count": 2
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
    When I do PUT /api/v4/eventfilter/test-eventfilter-to-failure-read-1/failures
    Then the response code should be 204
    When I do GET /api/v4/eventfilter/test-eventfilter-to-failure-read-1/failures
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-failure-to-read-2",
          "unread": false
        },
        {
          "_id": "test-eventfilter-failure-to-read-1",
          "unread": false
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
    When I do GET /api/v4/eventfilter/rules?search=test-eventfilter-to-failure-read-1&with_counts=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-to-failure-read-1",
          "failures_count": 2,
          "unread_failures_count": 0
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

  Scenario: given read request and no auth user should not allow access
    When I do PUT /api/v4/eventfilter/test-eventfilter-not-exist/failures
    Then the response code should be 401

  Scenario: given read request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/eventfilter/test-eventfilter-not-exist/failures
    Then the response code should be 403

  Scenario: given read request with not found event fitter rule should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/test-eventfilter-not-exist/failures
    Then the response code should be 404
