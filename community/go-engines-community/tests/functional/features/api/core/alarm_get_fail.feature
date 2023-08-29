Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/alarms?filters[]=not-exist
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter doesn't exist."
      }
    }
    """
    When I do GET /api/v4/alarms?time_field=unknown
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "time_field": "TimeField must be one of [t v.creation_date v.resolved v.last_update_date v.last_event_date] or empty."
      }
    }
    """
    When I do GET /api/v4/alarms?multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc&sort_by=v.duration&sort=desc
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "sort_by": "Can't be present both SortBy and MultiSort."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "multi_sort": "Invalid multi_sort value."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,bad
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "multi_sort": "Invalid multi_sort value."
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc,extra
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "multi_sort": "Invalid multi_sort value."
      }
    }
    """

  @concurrent
  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/alarms
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarms
    Then the response code should be 403
