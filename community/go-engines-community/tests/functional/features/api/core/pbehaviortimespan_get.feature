Feature: get time intervals for periodical behavior
  I need to be able to get time intervals on date, week and month view for periodical behavior.
  Application timezone "Europe/Paris".
  
  Scenario: given periodical behavior on Wednesday, Thursday and Friday and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given active periodical behavior on Wednesday, Thursday and Friday and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }},
      "type": "test-default-active-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "07-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 12:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 12:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 22:00" }},
        "to": {{ parseTime "09-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 12:00" }},
        "to": {{ parseTime "09-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday twice per day and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=HOURLY;BYDAY=WE,TH,FR;BYHOUR=12,16",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 14:00" }},
        "to": {{ parseTime "07-10-2020 16:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 14:00" }},
        "to": {{ parseTime "08-10-2020 16:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 14:00" }},
        "to": {{ parseTime "09-10-2020 16:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given active periodical behavior on Wednesday, Thursday and Friday twice per day and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=HOURLY;BYDAY=WE,TH,FR;BYHOUR=12,16",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }},
      "type": "test-default-active-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 14:00" }},
        "to": {{ parseTime "07-10-2020 16:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "07-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 12:00" }},
        "to": {{ parseTime "07-10-2020 14:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 16:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 14:00" }},
        "to": {{ parseTime "08-10-2020 16:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 12:00" }},
        "to": {{ parseTime "08-10-2020 14:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 16:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 14:00" }},
        "to": {{ parseTime "09-10-2020 16:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 22:00" }},
        "to": {{ parseTime "09-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 12:00" }},
        "to": {{ parseTime "09-10-2020 14:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 16:00" }},
        "to": {{ parseTime "09-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday and Wednesday view should return one timespan for the day
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "06-10-2020 22:00" }},
      "view_to": {{ parseTime "07-10-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given active periodical behavior on Wednesday, Thursday and Friday and Wednesday view should return one timespan for the day
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "06-10-2020 22:00" }},
      "view_to": {{ parseTime "07-10-2020 22:00" }},
      "type": "test-default-active-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "07-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 12:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday and Saturday view should return no timespans
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "02-10-2020 22:00" }},
      "view_to": {{ parseTime "03-10-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    []
	"""

  Scenario: given active periodical behavior on Wednesday, Thursday and Friday and Saturday view should return no timespans
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "02-10-2020 22:00" }},
      "view_to": {{ parseTime "03-10-2020 22:00" }},
      "type": "test-default-active-type"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    []
	"""

  Scenario: given two-day long weekly periodical behavior and second day view should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=WEEKLY",
      "start_at": {{ parseTime "27-07-2020 22:00" }},
      "end_at": {{ parseTime "29-07-2020 22:00" }},
      "view_from": {{ parseTime "11-08-2020 22:00" }},
      "view_to": {{ parseTime "12-08-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "11-08-2020 22:00" }},
        "to": {{ parseTime "12-08-2020 22:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given two-day long weekly active periodical behavior and second day view should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=WEEKLY",
      "start_at": {{ parseTime "27-07-2020 22:00" }},
      "end_at": {{ parseTime "29-07-2020 22:00" }},
      "view_from": {{ parseTime "11-08-2020 22:00" }},
      "view_to": {{ parseTime "12-08-2020 22:00" }},
      "type": "test-default-active-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "11-08-2020 22:00" }},
        "to": {{ parseTime "12-08-2020 22:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      }
    ]
	"""

  Scenario: given two-day affect weekly periodical behavior and second day view should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=WEEKLY",
      "start_at": {{ parseTime "28-07-2020 10:00" }},
      "end_at": {{ parseTime "29-07-2020 12:00" }},
      "view_from": {{ parseTime "11-08-2020 22:00" }},
      "view_to": {{ parseTime "12-08-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "11-08-2020 22:00" }},
        "to": {{ parseTime "12-08-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given two-day affect weekly active periodical behavior and second day view should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=WEEKLY",
      "start_at": {{ parseTime "28-07-2020 10:00" }},
      "end_at": {{ parseTime "29-07-2020 12:00" }},
      "view_from": {{ parseTime "11-08-2020 22:00" }},
      "view_to": {{ parseTime "12-08-2020 22:00" }},
      "type": "test-default-active-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "11-08-2020 22:00" }},
        "to": {{ parseTime "12-08-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "12-08-2020 12:00" }},
        "to": {{ parseTime "12-08-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
	"""

  Scenario: given week long monthly periodical behavior and day view in the middle of the week should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=MONTHLY",
      "start_at": {{ parseTime "02-08-2020 22:00" }},
      "end_at": {{ parseTime "09-08-2020 21:59" }},
      "view_from": {{ parseTime "04-09-2020 22:00" }},
      "view_to": {{ parseTime "05-09-2020 21:59" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "04-09-2020 22:00" }},
        "to": {{ parseTime "05-09-2020 21:59" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday during summer time and Wednesday view during winter should return timespans in the same timezone
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "27-10-2020 23:00" }},
      "view_to": {{ parseTime "28-10-2020 23:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "28-10-2020 11:00" }},
        "to": {{ parseTime "28-10-2020 13:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday without end date and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "09-10-2020 10:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "10-10-2020 22:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday without end date and Wednesday view should return one timespan for the day
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "07-10-2020 10:00" }},
      "view_from": {{ parseTime "06-10-2020 22:00" }},
      "view_to": {{ parseTime "07-10-2020 22:00" }},
      "type": "test-default-maintenance-type"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given daily periodical behavior with exdates and week view should return timespans for week with exdates
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "11-10-2020 22:00" }},
      "type": "test-default-maintenance-type",
      "exdates": [
        {
          "begin": {{ parseTime "04-10-2020 22:00" }},
          "end": {{ parseTime "05-10-2020 22:00" }},
          "type": "test-default-active-type"
        },
        {
          "begin": {{ parseTime "05-10-2020 22:00" }},
          "end": {{ parseTime "06-10-2020 11:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "07-10-2020 11:00" }},
          "end": {{ parseTime "07-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 11:00" }},
          "end": {{ parseTime "08-10-2020 11:30" }},
          "type": "test-default-active-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "09-10-2020 13:00" }},
          "end": {{ parseTime "09-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 09:00" }},
          "end": {{ parseTime "10-10-2020 11:30" }},
          "type": "test-default-active-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 10:30" }},
          "end": {{ parseTime "10-10-2020 18:30" }},
          "type": "test-default-pause-type"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "05-10-2020 10:00" }},
        "to": {{ parseTime "05-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 11:00" }},
        "to": {{ parseTime "06-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 10:00" }},
        "to": {{ parseTime "06-10-2020 11:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 11:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 11:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 11:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:30" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:00" }},
        "to": {{ parseTime "08-10-2020 11:30" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 10:30" }},
        "to": {{ parseTime "10-10-2020 12:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 10:00" }},
        "to": {{ parseTime "10-10-2020 10:30" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "11-10-2020 10:00" }},
        "to": {{ parseTime "11-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given active daily periodical behavior with exdates and week view should return timespans for week with exdates
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "11-10-2020 22:00" }},
      "type": "test-default-active-type",
      "exdates": [
        {
          "begin": {{ parseTime "04-10-2020 22:00" }},
          "end": {{ parseTime "05-10-2020 22:00" }},
          "type": "test-default-maintenance-type"
        },
        {
          "begin": {{ parseTime "05-10-2020 22:00" }},
          "end": {{ parseTime "06-10-2020 11:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "07-10-2020 11:00" }},
          "end": {{ parseTime "07-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 11:00" }},
          "end": {{ parseTime "08-10-2020 11:30" }},
          "type": "test-default-maintenance-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "09-10-2020 13:00" }},
          "end": {{ parseTime "09-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 09:00" }},
          "end": {{ parseTime "10-10-2020 11:30" }},
          "type": "test-default-maintenance-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 10:30" }},
          "end": {{ parseTime "10-10-2020 18:30" }},
          "type": "test-default-pause-type"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "05-10-2020 10:00" }},
        "to": {{ parseTime "05-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "04-10-2020 22:00" }},
        "to": {{ parseTime "05-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "05-10-2020 12:00" }},
        "to": {{ parseTime "05-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 10:00" }},
        "to": {{ parseTime "06-10-2020 11:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 11:00" }},
        "to": {{ parseTime "06-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "05-10-2020 22:00" }},
        "to": {{ parseTime "06-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 12:00" }},
        "to": {{ parseTime "06-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 11:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "07-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 11:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 12:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:00" }},
        "to": {{ parseTime "08-10-2020 11:30" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 11:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:30" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 12:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 22:00" }},
        "to": {{ parseTime "09-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 12:00" }},
        "to": {{ parseTime "09-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 10:00" }},
        "to": {{ parseTime "10-10-2020 11:30" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 11:30" }},
        "to": {{ parseTime "10-10-2020 12:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 22:00" }},
        "to": {{ parseTime "10-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 12:00" }},
        "to": {{ parseTime "10-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "11-10-2020 10:00" }},
        "to": {{ parseTime "11-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 22:00" }},
        "to": {{ parseTime "11-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "11-10-2020 12:00" }},
        "to": {{ parseTime "11-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
	"""

  Scenario: given daily periodical behavior with exceptions and week view should return timespans for week without exceptions
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """json
    {
      "name": "test-exception-pbehavior-timespans-with-exceptions-name-1",
      "description": "test-exception-pbehavior-timespans-with-exceptions-description-1",
      "exdates": [
        {
          "begin": {{ parseTime "04-10-2020 22:00" }},
          "end": {{ parseTime "05-10-2020 22:00" }},
          "type": "test-default-active-type"
        },
        {
          "begin": {{ parseTime "05-10-2020 22:00" }},
          "end": {{ parseTime "06-10-2020 11:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "07-10-2020 11:00" }},
          "end": {{ parseTime "07-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 11:00" }},
          "end": {{ parseTime "08-10-2020 11:30" }},
          "type": "test-default-active-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "09-10-2020 13:00" }},
          "end": {{ parseTime "09-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 09:00" }},
          "end": {{ parseTime "10-10-2020 11:30" }},
          "type": "test-default-active-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 10:30" }},
          "end": {{ parseTime "10-10-2020 18:30" }},
          "type": "test-default-pause-type"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "11-10-2020 22:00" }},
      "type": "test-default-maintenance-type",
      "exceptions": ["{{ .lastResponse._id }}"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "05-10-2020 10:00" }},
        "to": {{ parseTime "05-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 11:00" }},
        "to": {{ parseTime "06-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 10:00" }},
        "to": {{ parseTime "06-10-2020 11:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 11:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 11:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 11:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:30" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:00" }},
        "to": {{ parseTime "08-10-2020 11:30" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 10:30" }},
        "to": {{ parseTime "10-10-2020 12:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 10:00" }},
        "to": {{ parseTime "10-10-2020 10:30" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "11-10-2020 10:00" }},
        "to": {{ parseTime "11-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      }
    ]
	"""

  Scenario: given active daily periodical behavior with exceptions and week view should return timespans for week without exceptions
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """json
    {
      "name": "test-exception-pbehavior-timespans-with-exceptions-name-2",
      "description": "test-exception-pbehavior-timespans-with-exceptions-description-2",
      "exdates": [
        {
          "begin": {{ parseTime "04-10-2020 22:00" }},
          "end": {{ parseTime "05-10-2020 22:00" }},
          "type": "test-default-maintenance-type"
        },
        {
          "begin": {{ parseTime "05-10-2020 22:00" }},
          "end": {{ parseTime "06-10-2020 11:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "07-10-2020 11:00" }},
          "end": {{ parseTime "07-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 11:00" }},
          "end": {{ parseTime "08-10-2020 11:30" }},
          "type": "test-default-maintenance-type"
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }},
          "type": "test-default-pause-type"
        },
        {
          "begin": {{ parseTime "09-10-2020 13:00" }},
          "end": {{ parseTime "09-10-2020 22:00" }},
          "type": "test-default-inactive-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 09:00" }},
          "end": {{ parseTime "10-10-2020 11:30" }},
          "type": "test-default-maintenance-type"
        },
        {
          "begin": {{ parseTime "10-10-2020 10:30" }},
          "end": {{ parseTime "10-10-2020 18:30" }},
          "type": "test-default-pause-type"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-timespans:
    """json
    {
      "rrule": "FREQ=DAILY",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "11-10-2020 22:00" }},
      "type": "test-default-active-type",
      "exceptions": ["{{ .lastResponse._id }}"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "from": {{ parseTime "05-10-2020 10:00" }},
        "to": {{ parseTime "05-10-2020 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "04-10-2020 22:00" }},
        "to": {{ parseTime "05-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "05-10-2020 12:00" }},
        "to": {{ parseTime "05-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 10:00" }},
        "to": {{ parseTime "06-10-2020 11:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 11:00" }},
        "to": {{ parseTime "06-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "05-10-2020 22:00" }},
        "to": {{ parseTime "06-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 12:00" }},
        "to": {{ parseTime "06-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 11:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "07-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 11:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 12:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:00" }},
        "to": {{ parseTime "08-10-2020 11:30" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 11:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 11:30" }},
        "to": {{ parseTime "08-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "07-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 12:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "08-10-2020 22:00" }},
        "to": {{ parseTime "09-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 12:00" }},
        "to": {{ parseTime "09-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 10:00" }},
        "to": {{ parseTime "10-10-2020 11:30" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 11:30" }},
        "to": {{ parseTime "10-10-2020 12:00" }},
        "type": {
          "_id": "test-default-pause-type"
        }
      },
      {
        "from": {{ parseTime "09-10-2020 22:00" }},
        "to": {{ parseTime "10-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 12:00" }},
        "to": {{ parseTime "10-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "11-10-2020 10:00" }},
        "to": {{ parseTime "11-10-2020 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "from": {{ parseTime "10-10-2020 22:00" }},
        "to": {{ parseTime "11-10-2020 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "from": {{ parseTime "11-10-2020 12:00" }},
        "to": {{ parseTime "11-10-2020 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
	"""

  Scenario: given invalid request should return errors
    When I am admin
    When I do POST /api/v4/pbehavior-timespans
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "start_at":"StartAt is missing.",
        "view_from":"ViewFrom is missing.",
        "view_to":"ViewTo is missing.",
        "type":"Type is missing."
      }
    }
    """

  Scenario: given auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/pbehavior-timespans
    Then the response code should be 403
