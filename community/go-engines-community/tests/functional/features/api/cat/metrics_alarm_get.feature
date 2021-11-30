Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  Scenario: given get total_alarms hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&sampling=hour&from={{ nowDate }}&to={{ nowDate }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
        "data": [
          {
            "timestamp": {{ nowDate }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "1h" }},
            "value": 2
          },
          {
            "timestamp": {{ nowDateAdd "2h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "3h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "4h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "5h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "6h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "7h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "8h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "9h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "10h" }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "11h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "12h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "13h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "14h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "15h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "16h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "17h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "18h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "19h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "20h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "21h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "22h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "23h" }},
            "value": 0
          }
        ]
      }
    ]
    """

  Scenario: given get total_alarms day request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-72h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "-48h" }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 3
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      }
    ]
    """

  Scenario: given get total_alarms week request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&sampling=week&from={{ parseTime "06-09-2021 00:00" }}&to={{ parseTime "10-10-2021 00:00" }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
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
    """

  Scenario: given get total_alarms month request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&sampling=month&from={{ parseTime "01-06-2021 00:00" }}&to={{ parseTime "31-10-2021 00:00" }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
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
    """

  Scenario: given get total_alarms request with empty interval should return metrics with zeros
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
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
    """

  Scenario: given get total_alarms request with filter by entity infos should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-total-alarm-metrics-get-by-entity-infos
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-72h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "-48h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 2
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      }
    ]
    """

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
    When I do GET /api/v4/cat/metrics/alarm?filter=not-exist&from={{ now }}&to={{ now }}&sampling=day&parameters[]=total_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "filter \"not-exist\" not found"
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=not-exist&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameter.0": "parameter \"not-exist\" is not supported"
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_user_activity&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameter.0": "parameter \"total_user_activity\" is not supported"
      }
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?sampling=not-exist&from={{ now }}&to={{ now }}&parameters[]=total_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metrics/alarm
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metrics/alarm
    Then the response code should be 403

  Scenario: given get request with all parameters should return all metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?parameters[]=total_alarms&parameters[]=non_displayed_alarms&parameters[]=instruction_alarms&parameters[]=pbehavior_alarms&parameters[]=correlation_alarms&parameters[]=ack_alarms&parameters[]=cancel_ack_alarms&parameters[]=ack_without_cancel_alarms&parameters[]=ticket_alarms&parameters[]=without_ticket_alarms&parameters[]=ratio_correlation&parameters[]=ratio_instructions&parameters[]=ratio_tickets&parameters[]=ratio_non_displayed&parameters[]=average_ack&parameters[]=average_resolve&sampling=day&from={{ nowDateAdd "-1d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-all-alarm-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "title": "total_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 3
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "non_displayed_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "instruction_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "pbehavior_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "correlation_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ack_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 2
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "cancel_ack_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ack_without_cancel_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ticket_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "without_ticket_alarms",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 2
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ratio_correlation",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 33.33
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ratio_instructions",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 33.33
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ratio_tickets",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 33.33
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "ratio_non_displayed",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 33.33
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "average_ack",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 500
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      },
      {
        "title": "average_resolve",
        "data": [
          {
            "timestamp": {{ nowDateAdd "-24h" }},
            "value": 0
          },
          {
            "timestamp": {{ nowDate }},
            "value": 1000
          },
          {
            "timestamp": {{ nowDateAdd "24h" }},
            "value": 0
          }
        ]
      }
    ]
    """
