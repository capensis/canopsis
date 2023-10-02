Feature: Get event filters
  I need to be able to get event filters
  Only admin should be able to get event filters

  Scenario: given get request should return failures
    When I am admin
    When I do GET /api/v4/eventfilter/test-eventfilter-to-failure-get-1/failures
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-failure-to-get-2",
          "type": 1,
          "message": "test-eventfilter-failure-to-get-2-message",
          "t": 1592215338,
          "event": {
            "event_type": "check",
            "connector": "test-eventfilter-failure-to-get-2-connector",
            "connector_name": "test-eventfilter-failure-to-get-2-connector-name",
            "component": "test-eventfilter-failure-to-get-2-component",
            "resource": "test-eventfilter-failure-to-get-2-resource"
          },
          "unread": false
        },
        {
          "_id": "test-eventfilter-failure-to-get-1",
          "type": 0,
          "message": "test-eventfilter-failure-to-get-1-message",
          "t": 1592215337,
          "event": null,
          "unread": true
        },
        {
          "_id": "test-eventfilter-failure-to-get-3",
          "type": 2,
          "message": "test-eventfilter-failure-to-get-3-message",
          "t": 1492215338,
          "event": null,
          "unread": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get request with type filter should return filtered failures
    When I am admin
    When I do GET /api/v4/eventfilter/test-eventfilter-to-failure-get-1/failures?type=1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-failure-to-get-2",
          "type": 1,
          "message": "test-eventfilter-failure-to-get-2-message",
          "t": 1592215338,
          "event": {
            "event_type": "check",
            "connector": "test-eventfilter-failure-to-get-2-connector",
            "connector_name": "test-eventfilter-failure-to-get-2-connector-name",
            "component": "test-eventfilter-failure-to-get-2-component",
            "resource": "test-eventfilter-failure-to-get-2-resource"
          },
          "unread": false
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

  Scenario: given get rules request should return failures counts
    When I am admin
    When I do GET /api/v4/eventfilter/rules?search=test-eventfilter-to-failure-get-1&with_counts=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-to-failure-get-1",
          "events_count": 100,
          "failures_count": 2,
          "unread_failures_count": 1
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
    When I do GET /api/v4/eventfilter/test-eventfilter-not-exist/failures
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/eventfilter/test-eventfilter-not-exist/failures
    Then the response code should be 403

  Scenario: given get request with not found event fitter rule should return error
    When I am admin
    When I do GET /api/v4/eventfilter/test-eventfilter-not-exist/failures
    Then the response code should be 404
