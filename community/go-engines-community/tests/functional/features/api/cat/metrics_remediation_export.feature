Feature: Export remediation metrics
  I need to be able to export remediation metrics
  Only admin should be able to export remediation metrics

  Scenario: given export request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/remediation?sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&instruction=test-instruction-to-remediation-metrics-get
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
    assigned,executed,ratio,timestamp
    1,1,100,{{ parseTime "20-11-2021 00:00" }}
    1,0,0,{{ parseTime "21-11-2021 00:00" }}
    2,1,50,{{ parseTime "22-11-2021 00:00" }}
    4,2,50,{{ parseTime "23-11-2021 00:00" }}
    0,0,0,{{ parseTime "24-11-2021 00:00" }}

    """

  Scenario: given export request with empty interval should return metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/metrics-export/remediation?sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&instruction=test-instruction-to-remediation-metrics-get
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
    assigned,executed,ratio,timestamp
    0,0,0,{{ parseTime "06-09-2020 00:00" }}
    0,0,0,{{ parseTime "07-09-2020 00:00" }}
    0,0,0,{{ parseTime "08-09-2020 00:00" }}

    """

  Scenario: given export request with invalid query params should return bad request
    When I am admin
    When I do POST /api/v4/cat/metrics-export/remediation
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
    When I do POST /api/v4/cat/metrics-export/remediation?sampling=not-exist
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  Scenario: given export request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics-export/remediation
    Then the response code should be 401

  Scenario: given export request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics-export/remediation
    Then the response code should be 403
