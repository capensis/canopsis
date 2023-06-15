Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get perf data request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }},
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
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "value": 140,
          "aggregate_func": "sum",
          "unit": "%"
        },
        {
          "title": "cpu_%",
          "value": 80,
          "aggregate_func": "last",
          "unit": "%"
        },
        {
          "title": "cpu_%",
          "value": 46.66,
          "aggregate_func": "avg",
          "unit": "%"
        },
        {
          "title": "cpu_%",
          "value": 80,
          "aggregate_func": "max",
          "unit": "%"
        },
        {
          "title": "cpu_%",
          "value": 20,
          "aggregate_func": "min",
          "unit": "%"
        }
      ]
    }
    """

  @concurrent
  Scenario: given not exist perf data request should return empty response
    When I am admin
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }},
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
    Then the response body should contain:
    """json
    {
      "data": []
    }
    """
