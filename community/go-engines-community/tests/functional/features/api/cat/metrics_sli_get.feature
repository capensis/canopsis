Feature: Get SLI metrics
  I need to be able to get SLI metrics
  Only admin should be able to get SLI metrics

  Scenario: given get hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=hour&from={{ parseTimeTz "23-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 01:00" }},
          "downtime": 30,
          "maintenance": 0,
          "uptime": 3570
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 02:00" }},
          "downtime": 60,
          "maintenance": 0,
          "uptime": 3540
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 03:00" }},
          "downtime": 30,
          "maintenance": 30,
          "uptime": 3540
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 04:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 05:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 06:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 07:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 08:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 09:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 10:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 11:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 12:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 13:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 14:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 15:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 16:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 17:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 18:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 19:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 20:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 21:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 22:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 23:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 3600
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=hour&from={{ parseTimeTz "23-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 01:00" }},
          "downtime": 0.83,
          "maintenance": 0,
          "uptime": 99.17
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 02:00" }},
          "downtime": 1.66,
          "maintenance": 0,
          "uptime": 98.34
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 03:00" }},
          "downtime": 0.83,
          "maintenance": 0.83,
          "uptime": 98.34
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 04:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 05:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 06:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 07:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 08:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 09:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 10:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 11:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 12:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 13:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 14:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 15:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 16:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 17:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 18:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 19:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 20:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 21:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 22:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 23:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  Scenario: given get day request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=day&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "24-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        },
        {
          "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
          "downtime": 30,
          "maintenance": 0,
          "uptime": 86370
        },
        {
          "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
          "downtime": 120,
          "maintenance": 30,
          "uptime": 86250
        },
        {
          "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=day&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "24-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
          "downtime": 0.03,
          "maintenance": 0,
          "uptime": 99.97
        },
        {
          "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
          "downtime": 0.13,
          "maintenance": 0.03,
          "uptime": 99.84
        },
        {
          "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        }
      ]
    }
    """

  Scenario: given get week request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=week&from={{ parseTimeTz "06-09-2021 00:00" }}&to={{ parseTimeTz "10-10-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "06-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 604800
        },
        {
          "timestamp": {{ parseTimeTz "13-09-2021 00:00" }},
          "downtime": 30,
          "maintenance": 0,
          "uptime": 604770
        },
        {
          "timestamp": {{ parseTimeTz "20-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 604800
        },
        {
          "timestamp": {{ parseTimeTz "27-09-2021 00:00" }},
          "downtime": 60,
          "maintenance": 0,
          "uptime": 604740
        },
        {
          "timestamp": {{ parseTimeTz "04-10-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 604800
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=week&from={{ parseTimeTz "06-09-2021 00:00" }}&to={{ parseTimeTz "10-10-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "06-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "13-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "20-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "27-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "04-10-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        }
      ]
    }
    """

  Scenario: given get month request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=month&from={{ parseTimeTz "01-06-2021 00:00" }}&to={{ parseTimeTz "31-10-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "01-06-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 2592000
        },
        {
          "timestamp": {{ parseTimeTz "01-07-2021 00:00" }},
          "downtime": 30,
          "maintenance": 0,
          "uptime": 2678370
        },
        {
          "timestamp": {{ parseTimeTz "01-08-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 2678400
        },
        {
          "timestamp": {{ parseTimeTz "01-09-2021 00:00" }},
          "downtime": 90,
          "maintenance": 0,
          "uptime": 2591910
        },
        {
          "timestamp": {{ parseTimeTz "01-10-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 2682000
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=month&from={{ parseTimeTz "01-06-2021 00:00" }}&to={{ parseTimeTz "31-10-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "01-06-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "01-07-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "01-08-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "01-09-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        },
        {
          "timestamp": {{ parseTimeTz "01-10-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 100
        }
      ]
    }
    """

  Scenario: given get request with empty interval should return metrics with zeros
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=day&from={{ parseTimeTz "06-09-2020 00:00" }}&to={{ parseTimeTz "08-09-2020 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": []
    }
    """

  Scenario: given get request with filter by entity infos should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=day&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "24-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get-by-entity-infos
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        },
        {
          "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        },
        {
          "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        },
        {
          "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
          "downtime": 60,
          "maintenance": 60,
          "uptime": 86280
        },
        {
          "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
          "downtime": 0,
          "maintenance": 0,
          "uptime": 86400
        }
      ]
    }
    """

  Scenario: given get request with invalid query params should return bad request
    When I am admin
    When I do GET /api/v4/cat/metrics/sli
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
    When I do GET /api/v4/cat/metrics/sli?sampling=not-exist
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/sli?filter=not-exist&from={{ nowDateTz }}&to={{ nowDateTz }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter doesn't exist."
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metrics/sli
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metrics/sli
    Then the response code should be 403
