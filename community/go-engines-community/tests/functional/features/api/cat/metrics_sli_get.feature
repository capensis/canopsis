Feature: Get SLI metrics
  I need to be able to get SLI metrics
  Only admin should be able to get SLI metrics

  Scenario: given get hour request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=hour&from={{ nowDate }}&to={{ nowDate }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ nowDate }},
        "downtime": 30,
        "maintenance": 0,
        "uptime": 3570
      },
      {
        "timestamp": {{ nowDateAdd "1h" }},
        "downtime": 60,
        "maintenance": 0,
        "uptime": 3540
      },
      {
        "timestamp": {{ nowDateAdd "2h" }},
        "downtime": 30,
        "maintenance": 30,
        "uptime": 3540
      },
      {
        "timestamp": {{ nowDateAdd "3h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "4h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "5h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "6h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "7h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "8h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "9h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "10h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "11h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "12h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "13h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "14h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "15h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "16h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "17h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "18h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "19h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "20h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "21h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "22h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      },
      {
        "timestamp": {{ nowDateAdd "23h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 3600
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=hour&from={{ nowDate }}&to={{ nowDate }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ nowDate }},
        "downtime": 0.83,
        "maintenance": 0,
        "uptime": 99.17
      },
      {
        "timestamp": {{ nowDateAdd "1h" }},
        "downtime": 1.66,
        "maintenance": 0,
        "uptime": 98.34
      },
      {
        "timestamp": {{ nowDateAdd "2h" }},
        "downtime": 0.83,
        "maintenance": 0.83,
        "uptime": 98.34
      },
      {
        "timestamp": {{ nowDateAdd "3h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "4h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "5h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "6h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "7h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "8h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "9h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "10h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "11h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "12h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "13h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "14h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "15h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "16h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "17h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "18h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "19h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "20h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "21h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "22h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "23h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      }
    ]
    """

  Scenario: given get day request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ nowDateAdd "-72h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      },
      {
        "timestamp": {{ nowDateAdd "-48h" }},
        "downtime": 30,
        "maintenance": 0,
        "uptime": 86370
      },
      {
        "timestamp": {{ nowDateAdd "-24h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      },
      {
        "timestamp": {{ nowDate }},
        "downtime": 120,
        "maintenance": 30,
        "uptime": 86250
      },
      {
        "timestamp": {{ nowDateAdd "24h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ nowDateAdd "-72h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDateAdd "-48h" }},
        "downtime": 0.03,
        "maintenance": 0,
        "uptime": 99.97
      },
      {
        "timestamp": {{ nowDateAdd "-24h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ nowDate }},
        "downtime": 0.13,
        "maintenance": 0.03,
        "uptime": 99.84
      },
      {
        "timestamp": {{ nowDateAdd "24h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      }
    ]
    """

  Scenario: given get week request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=week&from={{ parseTime "06-09-2021 00:00" }}&to={{ parseTime "10-10-2021 00:00" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ parseTime "06-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 604800
      },
      {
        "timestamp": {{ parseTime "13-09-2021 00:00" }},
        "downtime": 30,
        "maintenance": 0,
        "uptime": 604770
      },
      {
        "timestamp": {{ parseTime "20-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 604800
      },
      {
        "timestamp": {{ parseTime "27-09-2021 00:00" }},
        "downtime": 60,
        "maintenance": 0,
        "uptime": 604740
      },
      {
        "timestamp": {{ parseTime "04-10-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 604800
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=week&from={{ parseTime "06-09-2021 00:00" }}&to={{ parseTime "10-10-2021 00:00" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ parseTime "06-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "13-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "20-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "27-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "04-10-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      }
    ]
    """

  Scenario: given get month request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=month&from={{ parseTime "01-06-2021 00:00" }}&to={{ parseTime "31-10-2021 00:00" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ parseTime "01-06-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 2592000
      },
      {
        "timestamp": {{ parseTime "01-07-2021 00:00" }},
        "downtime": 30,
        "maintenance": 0,
        "uptime": 2678370
      },
      {
        "timestamp": {{ parseTime "01-08-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 2678400
      },
      {
        "timestamp": {{ parseTime "01-09-2021 00:00" }},
        "downtime": 90,
        "maintenance": 0,
        "uptime": 2591910
      },
      {
        "timestamp": {{ parseTime "01-10-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 2678400
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/sli?in_percents=true&sampling=month&from={{ parseTime "01-06-2021 00:00" }}&to={{ parseTime "31-10-2021 00:00" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ parseTime "01-06-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "01-07-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "01-08-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "01-09-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      },
      {
        "timestamp": {{ parseTime "01-10-2021 00:00" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 100
      }
    ]
    """

  Scenario: given get total_alarms request with empty interval should return metrics with zeros
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&filter=test-filter-to-sli-metrics-get
    Then the response code should be 200
    Then the response body should be:
    """json
    []
    """

  Scenario: given get request with filter by entity infos should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/sli?sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-sli-metrics-get-by-entity-infos
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "timestamp": {{ nowDateAdd "-72h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      },
      {
        "timestamp": {{ nowDateAdd "-48h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      },
      {
        "timestamp": {{ nowDateAdd "-24h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      },
      {
        "timestamp": {{ nowDate }},
        "downtime": 60,
        "maintenance": 60,
        "uptime": 86280
      },
      {
        "timestamp": {{ nowDateAdd "24h" }},
        "downtime": 0,
        "maintenance": 0,
        "uptime": 86400
      }
    ]
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
    When I do GET /api/v4/cat/metrics/sli?filter=not-exist&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "filter \"not-exist\" not found"
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
