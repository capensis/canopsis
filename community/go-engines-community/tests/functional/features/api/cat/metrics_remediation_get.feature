Feature: Get remediation metrics
  I need to be able to get remediation metrics
  Only admin should be able to get remediation metrics

  Scenario: given get hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/remediation?sampling=hour&from={{ parseTime "23-11-2021 00:00" }}&to={{ parseTime "23-11-2021 00:00" }}&instruction=test-instruction-to-remediation-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTime "23-11-2021 00:00" }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        },
        {
          "timestamp": {{ parseTime "23-11-2021 01:00" }},
          "assigned": 1,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 02:00" }},
          "assigned": 2,
          "executed": 1,
          "ratio": 50
        },
        {
          "timestamp": {{ parseTime "23-11-2021 03:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 04:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 05:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 06:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 07:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 08:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 09:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 10:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 11:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 12:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 13:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 14:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 15:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 16:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 17:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 18:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 19:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 20:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 21:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 22:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "23-11-2021 23:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        }
      ],
      "meta": {
        "min_date": 1625097600
      }
    }
    """

  Scenario: given get request with empty interval should return metrics with zeros
    When I am admin
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&instruction=test-instruction-to-remediation-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTime "06-09-2020 00:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "07-09-2020 00:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        },
        {
          "timestamp": {{ parseTime "08-09-2020 00:00" }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  Scenario: given get request with invalid query params should return bad request
    When I am admin
    When I do GET /api/v4/cat/metrics/remediation
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "from": "From is missing.",
        "sampling": "Sampling is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=not-exist
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metrics/remediation
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metrics/remediation
    Then the response code should be 403
