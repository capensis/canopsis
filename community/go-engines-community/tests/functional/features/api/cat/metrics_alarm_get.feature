Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get created_alarms hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 01:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "23-11-2021 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 03:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 10:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "23-11-2021 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 23:00" }},
              "value": 0
            }
          ]
        }
      ],
      "meta": {
        "min_date": 1625097600
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms day request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
        "title": "created_alarms",
        "data": [
          {
            "timestamp": {{ parseTime "20-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTime "21-11-2021 00:00" }},
            "value": 1
          },
          {
            "timestamp": {{ parseTime "22-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTime "23-11-2021 00:00" }},
            "value": 3
          },
          {
            "timestamp": {{ parseTime "24-11-2021 00:00" }},
            "value": 0
          }
        ]
      }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms week request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=week&from={{ parseTime "06-09-2021 00:00" }}&to={{ parseTime "10-10-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "06-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "13-09-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "20-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "27-09-2021 00:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "04-10-2021 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms month request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=month&from={{ parseTime "01-06-2021 00:00" }}&to={{ parseTime "31-10-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "01-06-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "01-07-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "01-08-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "01-09-2021 00:00" }},
              "value": 3
            },
            {
              "timestamp": {{ parseTime "01-10-2021 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms request with empty interval should return metrics with zeros
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "06-09-2020 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "07-09-2020 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "08-09-2020 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms request with filter by entity infos should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get-by-entity-infos
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
        "title": "created_alarms",
        "data": [
          {
            "timestamp": {{ parseTime "20-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTime "21-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTime "22-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTime "23-11-2021 00:00" }},
            "value": 2
          },
          {
            "timestamp": {{ parseTime "24-11-2021 00:00" }},
            "value": 0
          }
        ]
      }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=active_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTime "23-11-2021 01:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 02:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 03:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 04:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 05:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 06:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 07:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 08:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 09:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "23-11-2021 10:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 11:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 12:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 13:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 14:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 15:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 16:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 17:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 18:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 19:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 20:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 21:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 22:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "23-11-2021 23:00" }},
              "value": 6
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms day request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=active_alarms&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "20-11-2021 00:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTime "21-11-2021 00:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 6
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms week request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=active_alarms&sampling=week&from={{ parseTime "06-09-2021 00:00" }}&to={{ parseTime "10-10-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
        "title": "active_alarms",
        "data": [
          {
            "timestamp": {{ parseTime "06-09-2021 00:00" }},
            "value": 1
          },
          {
            "timestamp": {{ parseTime "13-09-2021 00:00" }},
            "value": 2
          },
          {
            "timestamp": {{ parseTime "20-09-2021 00:00" }},
            "value": 2
          },
          {
            "timestamp": {{ parseTime "27-09-2021 00:00" }},
            "value": 4
          },
          {
            "timestamp": {{ parseTime "04-10-2021 00:00" }},
            "value": 4
          }
        ]
      }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms month request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=active_alarms&sampling=month&from={{ parseTime "01-06-2021 00:00" }}&to={{ parseTime "31-10-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "01-06-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "01-07-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "01-08-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "01-09-2021 00:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTime "01-10-2021 00:00" }},
              "value": 4
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=ratio_tickets&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTime "23-11-2021 01:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTime "23-11-2021 02:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 03:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 04:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 05:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 06:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 07:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 08:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 09:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTime "23-11-2021 10:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 11:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 12:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 13:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 14:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 15:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 16:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 17:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 18:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 19:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 20:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 21:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 22:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "23-11-2021 23:00" }},
              "value": 16.66
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets day request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=ratio_tickets&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTime "20-11-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTime "21-11-2021 00:00" }},
              "value": 40
            },
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 16.66
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets week request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=ratio_tickets&sampling=week&from={{ parseTime "06-09-2021 00:00" }}&to={{ parseTime "10-10-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTime "06-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "13-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "20-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "27-09-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTime "04-10-2021 00:00" }},
              "value": 25
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets month request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=ratio_tickets&sampling=month&from={{ parseTime "01-06-2021 00:00" }}&to={{ parseTime "31-10-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTime "01-06-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "01-07-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "01-08-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "01-09-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTime "01-10-2021 00:00" }},
              "value": 25
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get request with invalid query params should return bad request
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "from": "From is missing.",
        "parameters[]": "Parameters is missing.",
        "sampling": "Sampling is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter=not-exist&from={{ now }}&to={{ now }}&sampling=day&parameters[]=created_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter \"not-exist\" not found."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=not-exist&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameter.0": "Parameter \"not-exist\" is not supported."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_user_activity&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameter.0": "Parameter \"total_user_activity\" is not supported."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?sampling=not-exist&from={{ now }}&to={{ now }}&parameters[]=created_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metrics/alarm
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metrics/alarm
    Then the response code should be 403

  @concurrent
  Scenario: given get request with all parameters should return all metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&parameters[]=active_alarms&parameters[]=non_displayed_alarms&parameters[]=instruction_alarms&parameters[]=pbehavior_alarms&parameters[]=correlation_alarms&parameters[]=ack_alarms&parameters[]=cancel_ack_alarms&parameters[]=ack_active_alarms&parameters[]=ticket_active_alarms&parameters[]=without_ticket_active_alarms&parameters[]=ratio_correlation&parameters[]=ratio_instructions&parameters[]=ratio_tickets&parameters[]=ratio_non_displayed&parameters[]=average_ack&parameters[]=average_resolve&sampling=day&from={{ parseTime "22-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-all-alarm-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 3
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 3
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 3
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "pbehavior_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "correlation_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "cancel_ack_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "ack_active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 1
            }
          ]
        },
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 1
            }
          ]
        },
        {
          "title": "without_ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 2
            }
          ]
        },
        {
          "title": "ratio_correlation",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "ratio_instructions",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "ratio_non_displayed",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 500
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "average_resolve",
          "data": [
            {
              "timestamp": {{ parseTime "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 1000
            },
            {
              "timestamp": {{ parseTime "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given filter with old pattern should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get-by-old-pattern
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 01:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "23-11-2021 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 03:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 10:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "23-11-2021 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 23:00" }},
              "value": 0
            }
          ]
        }
      ],
      "meta": {
        "min_date": 1625097600
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms hour request should return metrics by widget filter
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&widget_filters[]=test-widget-filter-to-alarm-metrics-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 01:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTime "23-11-2021 02:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "23-11-2021 03:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 10:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "23-11-2021 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 23:00" }},
              "value": 0
            }
          ]
        }
      ],
      "meta": {
        "min_date": 1625097600
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms hour request should return metrics by widget filters
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&widget_filters[]=test-widget-filter-to-alarm-metrics-get-2&widget_filters[]=test-widget-filter-to-alarm-metrics-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTime "23-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 01:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTime "23-11-2021 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 03:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 10:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTime "23-11-2021 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTime "23-11-2021 23:00" }},
              "value": 0
            }
          ]
        }
      ],
      "meta": {
        "min_date": 1625097600
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms hour request with both metrics and widget filters should return error
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&widget_filters[]=test-widget-filter-to-alarm-metrics-get&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget_filters": "Can't be present both WidgetFilters and KpiFilter."
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms hour request with not exist widget filter should return error
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&widget_filters[]=test-widget-filter-not-exist
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget_filters": "filter \"test-widget-filter-not-exist\" not found"
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms hour request with invalid widget filter should return error
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=created_alarms&sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&widget_filters[]=test-widget-filter-to-alarm-metrics-get-3
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget_filters": "filter \"test-widget-filter-to-alarm-metrics-get-3\" cannot be applied"
      }
    }
    """
