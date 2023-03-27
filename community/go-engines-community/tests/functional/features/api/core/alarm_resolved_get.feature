Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get resolved by entity request should return resolved entity alarms
    When I am admin
    When I do GET /api/v4/resolved-alarms?_id=test-resource-to-alarm-get-3/test-component-to-alarm-get&sort_by=_id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-4"
        },
        {
          "_id": "test-alarm-to-get-3"
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
    When I do GET /api/v4/resolved-alarms?_id=test-resource-to-alarm-get-3/test-component-to-alarm-get&tstart=1597030241
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-3"
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
    When I do GET /api/v4/resolved-alarms?_id=test-resource-to-alarm-get-3/test-component-to-alarm-get&tstop=1597030141
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-4"
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

  @concurrent
  Scenario: given get resolved by entity request should return validation error
    When I am admin
    When I do GET /api/v4/resolved-alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID is missing."
      }
    }
    """

  @concurrent
  Scenario: given get resolved by entity request should return not found error
    When I am admin
    When I do GET /api/v4/resolved-alarms?_id=not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given get resolved by entity unauth request should not allow access
    When I do GET /api/v4/resolved-alarms
    Then the response code should be 401

  @concurrent
  Scenario: given get resolved by entity and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/resolved-alarms
    Then the response code should be 403
