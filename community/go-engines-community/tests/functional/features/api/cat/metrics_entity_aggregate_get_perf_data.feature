Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get perf data request should return metrics
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }},
      "entity": "test-entity-to-alarm-metrics-get-1",
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        },
        {
          "metric": "cpu_ms",
          "aggregate_func": "sum"
        },
        {
          "metric": "memory_GB",
          "aggregate_func": "sum"
        },
        {
          "metric": "cpu*",
          "aggregate_func": "sum"
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
          "title": "cpu",
          "value": 70,
          "aggregate_func": "sum"
        }
      ]
    }
    """
