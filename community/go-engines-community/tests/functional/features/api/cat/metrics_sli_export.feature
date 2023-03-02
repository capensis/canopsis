Feature: Export SLI metrics
  I need to be able to export SLI metrics
  Only admin should be able to export SLI metrics

  Scenario: given export request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/sli?sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    downtime,maintenance,timestamp,uptime
    0,0,{{ parseTime "20-11-2021 00:00" }},86400
    30,0,{{ parseTime "21-11-2021 00:00" }},86370
    0,0,{{ parseTime "22-11-2021 00:00" }},86400
    120,30,{{ parseTime "23-11-2021 00:00" }},86250
    0,0,{{ parseTime "24-11-2021 00:00" }},86400

    """
    When I do POST /api/v4/cat/metrics-export/sli?in_percents=true&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    downtime,maintenance,timestamp,uptime
    0,0,{{ parseTime "20-11-2021 00:00" }},100
    0.03,0,{{ parseTime "21-11-2021 00:00" }},99.97
    0,0,{{ parseTime "22-11-2021 00:00" }},100
    0.13,0.03,{{ parseTime "23-11-2021 00:00" }},99.84
    0,0,{{ parseTime "24-11-2021 00:00" }},100

    """

  Scenario: given export request with empty interval should return metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/metrics-export/sli?sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv

    """

  Scenario: given export request with filter by entity infos should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/sli?sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-sli-metrics-get-by-entity-infos
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    downtime,maintenance,timestamp,uptime
    0,0,{{ parseTime "20-11-2021 00:00" }},86400
    0,0,{{ parseTime "21-11-2021 00:00" }},86400
    0,0,{{ parseTime "22-11-2021 00:00" }},86400
    60,60,{{ parseTime "23-11-2021 00:00" }},86280
    0,0,{{ parseTime "24-11-2021 00:00" }},86400

    """

  Scenario: given export request with invalid query params should return bad request
    When I am admin
    When I do POST /api/v4/cat/metrics-export/sli
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
    When I do POST /api/v4/cat/metrics-export/sli?sampling=not-exist
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """
    When I do POST /api/v4/cat/metrics-export/sli?filter=not-exist&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter \"not-exist\" not found."
      }
    }
    """

  Scenario: given export request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics-export/sli
    Then the response code should be 401

  Scenario: given export request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics-export/sli
    Then the response code should be 403
