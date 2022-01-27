Feature: get time intervals for periodical behavior
  I need to be able to get time intervals on date, week and month view
  for periodical behavior.
  Application timezone "Europe/Paris".
  Test month "October 2020", time change date "25 October 2020".

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
  and month view should return timespans for 5 weeks
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "30-09-2020 22:00" }},
      "view_to": {{ parseTime "31-10-2020 23:00" }},
      "by_date": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "30-09-2020 22:00" }},
        "to": {{ parseTime "01-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "13-10-2020 22:00" }},
        "to": {{ parseTime "15-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "20-10-2020 22:00" }},
        "to": {{ parseTime "22-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "27-10-2020 23:00" }},
        "to": {{ parseTime "29-10-2020 23:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
  and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }}
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 12:00" }}
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
  and Wednesday view should return one timespan for the day
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "06-10-2020 22:00" }},
      "view_to": {{ parseTime "07-10-2020 22:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 12:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
  and Saturday view should return no timespans
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "02-10-2020 22:00" }},
      "view_to": {{ parseTime "03-10-2020 22:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    []
	"""

  Scenario: given weekly periodical behavior and month view should return timespans for 5 weeks
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=WEEKLY",
      "start_at": {{ parseTime "03-08-2020 22:00" }},
      "end_at": {{ parseTime "05-08-2020 22:00" }},
      "view_from": {{ parseTime "31-07-2020 22:00" }},
      "view_to": {{ parseTime "31-08-2020 22:00" }},
      "by_date": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "03-08-2020 22:00" }},
        "to": {{ parseTime "05-08-2020 22:00" }}
      },
      {
        "from": {{ parseTime "10-08-2020 22:00" }},
        "to": {{ parseTime "12-08-2020 22:00" }}
      },
      {
        "from": {{ parseTime "17-08-2020 22:00" }},
        "to": {{ parseTime "19-08-2020 22:00" }}
      },
      {
        "from": {{ parseTime "24-08-2020 22:00" }},
        "to": {{ parseTime "26-08-2020 22:00" }}
      },
      {
        "from": {{ parseTime "31-08-2020 22:00" }},
        "to": {{ parseTime "31-08-2020 22:00" }}
      }
    ]
	"""

  Scenario: given monthly periodical behavior and year view should return timespans for 12 months
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=MONTHLY",
      "start_at": {{ parseTime "02-01-2020 23:00" }},
      "end_at": {{ parseTime "08-01-2020 23:00" }},
      "view_from": {{ parseTime "31-12-2019 23:00" }},
      "view_to": {{ parseTime "31-12-2020 23:00" }},
      "by_date": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "02-01-2020 23:00" }},
        "to": {{ parseTime "08-01-2020 23:00" }}
      },
      {
        "from": {{ parseTime "02-02-2020 23:00" }},
        "to": {{ parseTime "08-02-2020 23:00" }}
      },
      {
        "from": {{ parseTime "02-03-2020 23:00" }},
        "to": {{ parseTime "08-03-2020 23:00" }}
      },
      {
        "from": {{ parseTime "02-04-2020 22:00" }},
        "to": {{ parseTime "08-04-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-05-2020 22:00" }},
        "to": {{ parseTime "08-05-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-06-2020 22:00" }},
        "to": {{ parseTime "08-06-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-07-2020 22:00" }},
        "to": {{ parseTime "08-07-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-08-2020 22:00" }},
        "to": {{ parseTime "08-08-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-09-2020 22:00" }},
        "to": {{ parseTime "08-09-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "02-11-2020 23:00" }},
        "to": {{ parseTime "08-11-2020 23:00" }}
      },
      {
        "from": {{ parseTime "02-12-2020 23:00" }},
        "to": {{ parseTime "08-12-2020 23:00" }}
      }
    ]
	"""

  Scenario: given two-day long weekly periodical behavior and second day view should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=WEEKLY",
      "start_at": {{ parseTime "27-07-2020 22:00" }},
      "end_at": {{ parseTime "29-07-2020 21:59" }},
      "view_from": {{ parseTime "11-08-2020 22:00" }},
      "view_to": {{ parseTime "12-08-2020 21:59" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "11-08-2020 22:00" }},
        "to": {{ parseTime "12-08-2020 21:59" }}
      }
    ]
	"""

  Scenario: given week long monthly periodical behavior and day view in the middle of the week should return timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=MONTHLY",
      "start_at": {{ parseTime "02-08-2020 22:00" }},
      "end_at": {{ parseTime "09-08-2020 21:59" }},
      "view_from": {{ parseTime "04-09-2020 22:00" }},
      "view_to": {{ parseTime "05-09-2020 21:59" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "04-09-2020 22:00" }},
        "to": {{ parseTime "05-09-2020 21:59" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
  during summer time and Wednesday view during winter should return timespans in the same timezone
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "27-10-2020 23:00" }},
      "view_to": {{ parseTime "28-10-2020 23:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "28-10-2020 11:00" }},
        "to": {{ parseTime "28-10-2020 13:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
    with exdates and month view should return timespans for 5 weeks
    without exdates
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "30-09-2020 22:00" }},
      "view_to": {{ parseTime "31-10-2020 23:00" }},
      "exdates": [
        {
          "begin": {{ parseTime "07-10-2020 22:00" }},
          "end": {{ parseTime "08-10-2020 22:00" }}
        },
        {
          "begin": {{ parseTime "14-10-2020 22:00" }},
          "end": {{ parseTime "15-10-2020 11:00" }}
        },
        {
          "begin": {{ parseTime "22-10-2020 11:00" }},
          "end": {{ parseTime "22-10-2020 22:00" }}
        },
        {
          "begin": {{ parseTime "29-10-2020 11:00" }},
          "end": {{ parseTime "29-10-2020 11:30" }}
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }}
        },
        {
          "begin": {{ parseTime "16-10-2020 13:00" }},
          "end": {{ parseTime "16-10-2020 22:00" }}
        },
        {
          "begin": {{ parseTime "21-10-2020 09:00" }},
          "end": {{ parseTime "21-10-2020 10:30" }}
        },
        {
          "begin": {{ parseTime "21-10-2020 10:00" }},
          "end": {{ parseTime "21-10-2020 18:30" }}
        }
      ],
      "by_date": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "30-09-2020 22:00" }},
        "to": {{ parseTime "01-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "06-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "08-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "13-10-2020 22:00" }},
        "to": {{ parseTime "15-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "21-10-2020 22:00" }},
        "to": {{ parseTime "22-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "27-10-2020 23:00" }},
        "to": {{ parseTime "29-10-2020 23:00" }}
      }
    ]
	"""

  Scenario: given daily periodical behavior with exdates and week view
    should return timespans for week without exdates
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "11-10-2020 22:00" }},
      "exdates": [
        {
          "begin": {{ parseTime "04-10-2020 22:00" }},
          "end": {{ parseTime "05-10-2020 22:00" }}
        },
        {
          "begin": {{ parseTime "05-10-2020 22:00" }},
          "end": {{ parseTime "06-10-2020 11:00" }}
        },
        {
          "begin": {{ parseTime "07-10-2020 11:00" }},
          "end": {{ parseTime "07-10-2020 22:00" }}
        },
        {
          "begin": {{ parseTime "08-10-2020 11:00" }},
          "end": {{ parseTime "08-10-2020 11:30" }}
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }}
        },
        {
          "begin": {{ parseTime "09-10-2020 13:00" }},
          "end": {{ parseTime "09-10-2020 22:00" }}
        },
        {
          "begin": {{ parseTime "10-10-2020 09:00" }},
          "end": {{ parseTime "10-10-2020 10:30" }}
        },
        {
          "begin": {{ parseTime "10-10-2020 10:00" }},
          "end": {{ parseTime "10-10-2020 18:30" }}
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "06-10-2020 11:00" }},
        "to": {{ parseTime "06-10-2020 12:00" }}
      },
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 11:00" }}
      },
      {
        "from": {{ parseTime "08-10-2020 10:00" }},
        "to": {{ parseTime "08-10-2020 11:00" }}
      },
      {
        "from": {{ parseTime "08-10-2020 11:30" }},
        "to": {{ parseTime "08-10-2020 12:00" }}
      },
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "09-10-2020 12:00" }}
      },
      {
        "from": {{ parseTime "11-10-2020 10:00" }},
        "to": {{ parseTime "11-10-2020 12:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
    with exceptions and month view should return timespans for 5 weeks
    without exceptions
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "Test timespans with exceptions on month view",
      "description": "Test timespans with exceptions on month view",
      "exdates": [
        {
          "begin": {{ parseTime "07-10-2020 22:00" }},
          "end": {{ parseTime "08-10-2020 22:00" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "14-10-2020 22:00" }},
          "end": {{ parseTime "15-10-2020 11:00" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "22-10-2020 11:00" }},
          "end": {{ parseTime "22-10-2020 22:00" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "29-10-2020 11:00" }},
          "end": {{ parseTime "29-10-2020 11:30" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "08-10-2020 22:00" }},
          "end": {{ parseTime "09-10-2020 09:00" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "16-10-2020 13:00" }},
          "end": {{ parseTime "16-10-2020 22:00" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "21-10-2020 09:00" }},
          "end": {{ parseTime "21-10-2020 10:30" }},
          "type": "test-type-to-pbh-edit-1"
        },
        {
          "begin": {{ parseTime "21-10-2020 10:00" }},
          "end": {{ parseTime "21-10-2020 18:30" }},
          "type": "test-type-to-pbh-edit-1"
        }
      ]
    }
    """
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "30-09-2020 22:00" }},
      "view_to": {{ parseTime "31-10-2020 23:00" }},
      "exceptions": ["{{ .lastResponse._id }}"],
      "by_date": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "30-09-2020 22:00" }},
        "to": {{ parseTime "01-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "06-10-2020 22:00" }},
        "to": {{ parseTime "06-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "08-10-2020 22:00" }},
        "to": {{ parseTime "08-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "13-10-2020 22:00" }},
        "to": {{ parseTime "15-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "21-10-2020 22:00" }},
        "to": {{ parseTime "22-10-2020 22:00" }}
      },
      {
        "from": {{ parseTime "27-10-2020 23:00" }},
        "to": {{ parseTime "29-10-2020 23:00" }}
      }
    ]
	"""

#    TODO fix
#  Scenario: given daily periodical behavior with exceptions and week view
#    should return timespans for week without exceptions
#    When I am admin
#    When I do POST /api/v4/pbehavior-exceptions:
#    """
#    {
#      "name": "Test timespans with exceptions on week view",
#      "description": "Test timespans with exceptions on week view",
#      "exdates": [
#        {
#          "begin": {{ parseTime "04-10-2020 22:00" }},
#          "end": {{ parseTime "05-10-2020 22:00" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "05-10-2020 22:00" }},
#          "end": {{ parseTime "06-10-2020 11:00" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "07-10-2020 11:00" }},
#          "end": {{ parseTime "07-10-2020 22:00" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "08-10-2020 11:00" }},
#          "end": {{ parseTime "08-10-2020 11:30" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "08-10-2020 22:00" }},
#          "end": {{ parseTime "09-10-2020 09:00" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "09-10-2020 13:00" }},
#          "end": {{ parseTime "09-10-2020 22:00" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "10-10-2020 09:00" }},
#          "end": {{ parseTime "10-10-2020 10:30" }},
#          "type": "test-type-to-pbh-edit-1"
#        },
#        {
#          "begin": {{ parseTime "10-10-2020 10:00" }},
#          "end": {{ parseTime "10-10-2020 18:30" }},
#          "type": "test-type-to-pbh-edit-1"
#        }
#      ]
#    }
#    """
#    When I do POST /api/v4/pbehavior-timespans:
#    """
#    {
#      "rrule": "FREQ=DAILY",
#      "start_at": {{ parseTime "01-10-2020 10:00" }},
#      "end_at": {{ parseTime "01-10-2020 12:00" }},
#      "view_from": {{ parseTime "04-10-2020 22:00" }},
#      "view_to": {{ parseTime "11-10-2020 22:00" }},
#      "exceptions": ["{{ .lastResponse._id }}"]
#    }
#    """
#    Then the response code should be 200
#    Then the response body should be:
#    """
#    [
#      {
#        "from": {{ parseTime "06-10-2020 11:00" }},
#        "to": {{ parseTime "06-10-2020 12:00" }}
#      },
#      {
#        "from": {{ parseTime "07-10-2020 10:00" }},
#        "to": {{ parseTime "07-10-2020 11:00" }}
#      },
#      {
#        "from": {{ parseTime "08-10-2020 10:00" }},
#        "to": {{ parseTime "08-10-2020 11:00" }}
#      },
#      {
#        "from": {{ parseTime "08-10-2020 11:30" }},
#        "to": {{ parseTime "08-10-2020 12:00" }}
#      },
#      {
#        "from": {{ parseTime "09-10-2020 10:00" }},
#        "to": {{ parseTime "09-10-2020 12:00" }}
#      },
#      {
#        "from": {{ parseTime "11-10-2020 10:00" }},
#        "to": {{ parseTime "11-10-2020 12:00" }}
#      },
#      {
#        "from": {{ parseTime "04-10-2020 22:00" }},
#        "to": {{ parseTime "05-10-2020 22:00" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "05-10-2020 22:00" }},
#        "to": {{ parseTime "06-10-2020 11:00" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "07-10-2020 11:00" }},
#        "to": {{ parseTime "07-10-2020 22:00" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "08-10-2020 11:00" }},
#        "to": {{ parseTime "08-10-2020 11:30" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "08-10-2020 22:00" }},
#        "to": {{ parseTime "09-10-2020 09:00" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "09-10-2020 13:00" }},
#        "to": {{ parseTime "09-10-2020 22:00" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "10-10-2020 09:00" }},
#        "to": {{ parseTime "10-10-2020 10:30" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      },
#      {
#        "from": {{ parseTime "10-10-2020 10:00" }},
#        "to": {{ parseTime "10-10-2020 18:30" }},
#        "type": {
#          "_id": "test-type-to-pbh-edit-1",
#          "description": "Pbh edit 1 State type",
#          "icon_name": "test-to-pbh-edit-1-icon",
#          "name": "Pbh edit 1 State",
#          "priority": 10,
#          "type": "active"
#        }
#      }
#    ]
#	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
    without end date and month view should return one month long timespan
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "view_from": {{ parseTime "30-09-2020 22:00" }},
      "view_to": {{ parseTime "31-10-2020 23:00" }},
      "by_date": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "30-09-2020 22:00" }},
        "to": {{ parseTime "31-10-2020 23:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
    without end date and week view should return timespans for 3 days
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "09-10-2020 10:00" }},
      "view_from": {{ parseTime "04-10-2020 22:00" }},
      "view_to": {{ parseTime "10-10-2020 22:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "09-10-2020 10:00" }},
        "to": {{ parseTime "10-10-2020 22:00" }}
      }
    ]
	"""

  Scenario: given periodical behavior on Wednesday, Thursday and Friday
    without end date and Wednesday view should return one timespan for the day
    When I am admin
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=WE,TH,FR",
      "start_at": {{ parseTime "07-10-2020 10:00" }},
      "view_from": {{ parseTime "06-10-2020 22:00" }},
      "view_to": {{ parseTime "07-10-2020 22:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    [
      {
        "from": {{ parseTime "07-10-2020 10:00" }},
        "to": {{ parseTime "07-10-2020 22:00" }}
      }
    ]
	"""

  Scenario: given invalid request should return errors
    When I am admin
    When I do POST /api/v4/pbehavior-timespans
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "start_at":"StartAt is missing.",
        "view_from":"ViewFrom is missing.",
        "view_to":"ViewTo is missing."
      }
    }
    """

  Scenario: given auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/pbehavior-timespans:
    """
    {
      "rrule": "FREQ=DAILY;BYDAY=MO,TU,WE",
      "start_at": {{ parseTime "01-10-2020 10:00" }},
      "end_at": {{ parseTime "01-10-2020 12:00" }},
      "view_from": {{ parseTime "30-09-2020 22:00" }},
      "view_to": {{ parseTime "31-10-2020 23:00" }}
      "by_date": true
    }
    """
    Then the response code should be 403
