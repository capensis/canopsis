Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given get search request should return alarms only
    with string in connector, connector_name, component or resource fields
    When I am admin
    When I do GET /api/v4/alarm-counters?search=test-alarmcounter-get
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 4,
      "total_active": 4,
      "snooze": 0,
      "ack": 2,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get time inverval request should return alarms which were created
    in this time interval.
    When I am admin
    When I do GET /api/v4/alarm-counters?search=test-alarmcounter-get&tstart=1596931200&tstop=1597017600
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 1,
      "total_active": 1,
      "snooze": 0,
      "ack": 0,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get opened request should return only opened alarms
    When I am admin
    When I do GET /api/v4/alarm-counters?search=test-alarmcounter-get&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 3,
      "total_active": 3,
      "snooze": 0,
      "ack": 1,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get resolved request should return only resolved alarms
    When I am admin
    When I do GET /api/v4/alarm-counters?search=test-alarmcounter-get&resolved=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 1,
      "total_active": 1,
      "snooze": 0,
      "ack": 1,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get filter request should return alarms which are matched to the filter
    When I am admin
    When I do GET /api/v4/alarm-counters?filter={"$or":[{"uid":"test-alarmcounter-get-2"},{"uid":"test-alarmcounter-get-4"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 2,
      "total_active": 2,
      "snooze": 0,
      "ack": 2,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get search by infos.customer.value field request should return alarms
    only with string in entity.infos.customer.value field
    When I am admin
    When I do GET /api/v4/alarm-counters?search=test-alarmcounter-get-customer&active_columns[]=infos.customer.value
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 1,
      "total_active": 1,
      "snooze": 0,
      "ack": 0,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get correlation request should return meta alarms or alarms without parent
    When I am admin
    When I do GET /api/v4/alarm-counters?search=test-alarmcounter-get&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "total": 3,
      "total_active": 3,
      "snooze": 0,
      "ack": 1,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/alarm-counters
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarm-counters
    Then the response code should be 403
