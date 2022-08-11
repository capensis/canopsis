Feature: Get map's alarms
  I need to be able to get a map's alarms

  Scenario: given get request should return opened alarms
    When I am admin
    When I do GET /api/v4/cat/map-alarms/test-map-to-alarm-get-1?sort_by=v.resource&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-map-alarm-get-2"
        },
        {
          "_id": "test-alarm-to-map-alarm-get-1"
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

  Scenario: given map without alarms should return emtpy response
    When I am admin
    When I do GET /api/v4/cat/map-alarms/test-map-to-alarm-get-2
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

  Scenario: given get request should return not found error
    When I am admin
    When I do GET /api/v4/cat/map-alarms/test-map-not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/cat/map-alarms/test-map-to-alarm-get
    Then the response code should be 401

  Scenario: given get and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/map-alarms/test-map-to-alarm-get
    Then the response code should be 403
