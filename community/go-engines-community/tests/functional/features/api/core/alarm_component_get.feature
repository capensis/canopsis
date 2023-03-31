Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get by component request should return opened resource alarms
    When I am admin
    When I do GET /api/v4/component-alarms?_id=test-component-to-alarm-get&sort_by=v.resource&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-2"
        },
        {
          "_id": "test-alarm-to-get-1"
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

  @concurrent
  Scenario: given get by component request should return validation error
    When I am admin
    When I do GET /api/v4/component-alarms
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
  Scenario: given get by component request should return not found error
    When I am admin
    When I do GET /api/v4/component-alarms?_id=not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given get by component unauth request should not allow access
    When I do GET /api/v4/component-alarms
    Then the response code should be 401

  @concurrent
  Scenario: given get by component and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/component-alarms
    Then the response code should be 403
