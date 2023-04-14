Feature: Export alarm metrics
  I need to be able to export alarm metrics
  Only admin should be able to export alarm metrics

  @concurrent
  Scenario: given export request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "03-07-2022 00:00" }}
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
    func,metric,timestamp,unit,value
    sum,cpu_%,{{ parseTime "01-07-2022 00:00" }},%,140
    sum,cpu_%,{{ parseTime "02-07-2022 00:00" }},%,0
    sum,cpu_%,{{ parseTime "03-07-2022 00:00" }},%,0

    """

  @concurrent
  Scenario: given get regexp perf data request should return all matched metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu*",
          "aggregate_func": "max"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }}
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
    func,metric,timestamp,unit,value
    max,cpu,{{ parseTime "01-07-2022 00:00" }},,70
    max,cpu_%,{{ parseTime "01-07-2022 00:00" }},%,80

    """

  @concurrent
  Scenario: given get perf data with base metrics request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "max"
        },
        {
          "metric": "cpu_%",
          "aggregate_func": "max"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }}
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
    func,metric,timestamp,unit,value
    ,created_alarms,{{ parseTime "01-07-2022 00:00" }},,0
    max,cpu_%,{{ parseTime "01-07-2022 00:00" }},%,80

    """
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "max"
        },
        {
          "metric": "cpu_ms",
          "aggregate_func": "max"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }}
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
    metric,timestamp,value
    created_alarms,{{ parseTime "01-07-2022 00:00" }},0

    """

  @concurrent
  Scenario: given get not exist perf data request should return empty response
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_ms",
          "aggregate_func": "max"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }}
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

  @concurrent
  Scenario: given get not exist aggregate func request should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "cumulative_sum"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }}
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.aggregate_func": "AggregateFunc must be one of [sum last avg min max] or empty."
      }
    }
    """

  @concurrent
  Scenario: given get empty aggregate func request should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "from": {{ parseTime "01-07-2022 00:00" }},
      "to": {{ parseTime "01-07-2022 00:00" }}
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.metric": "Metric doesn't exist."
      }
    }
    """

