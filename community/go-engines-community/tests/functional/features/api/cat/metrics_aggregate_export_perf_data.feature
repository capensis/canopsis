Feature: Export aggregated metrics
  I need to be able to export aggregated metrics
  Only admin should be able to export aggregated metrics

  @concurrent
  Scenario: given get perf data request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }},
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        },
        {
          "metric": "cpu_%",
          "aggregate_func": "last"
        },
        {
          "metric": "cpu_%",
          "aggregate_func": "avg"
        },
        {
          "metric": "cpu_%",
          "aggregate_func": "max"
        },
        {
          "metric": "cpu_%",
          "aggregate_func": "min"
        }
      ]
    }
    """
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
    func,metric,unit,value
    sum,cpu_%,%,140
    last,cpu_%,%,80
    avg,cpu_%,%,46.66
    max,cpu_%,%,80
    min,cpu_%,%,20

    """

  @concurrent
  Scenario: given get regexp perf data request should return all matched metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }},
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "cpu*",
          "aggregate_func": "sum"
        }
      ]
    }
    """
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
    func,metric,unit,value
    sum,cpu,,70
    sum,cpu_%,%,140

    """

  @concurrent
  Scenario: given get perf data with base metrics request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }},
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        },
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        }
      ]
    }
    """
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
    func,metric,unit,value
    sum,created_alarms,,0
    sum,cpu_%,%,140

    """
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }},
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        },
        {
          "metric": "cpu_ms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
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
    func,metric,value
    sum,created_alarms,0

    """

  @concurrent
  Scenario: given not exist perf data request should return empty response
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }},
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "cpu_ms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
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
