Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get perf data hour request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
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
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "aggregate_func": "sum",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 01:00" }},
              "value": 60
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 03:00" }},
              "value": 80
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 10:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 23:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "cpu_%",
          "aggregate_func": "last",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 01:00" }},
              "value": 40
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 03:00" }},
              "value": 80
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 10:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 23:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "cpu_%",
          "aggregate_func": "avg",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 01:00" }},
              "value": 30
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 03:00" }},
              "value": 80
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 10:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 23:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "cpu_%",
          "aggregate_func": "max",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 01:00" }},
              "value": 40
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 03:00" }},
              "value": 80
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 10:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 23:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "cpu_%",
          "aggregate_func": "min",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 01:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 03:00" }},
              "value": 80
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 10:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2022 23:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get regexp perf data request should return all matched metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
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
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu",
          "aggregate_func": "max",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 70
            }
          ]
        },
        {
          "title": "cpu_%",
          "aggregate_func": "max",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 80
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "max"
        },
        {
          "metric": "cpu*",
          "aggregate_func": "max"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "aggregate_func": "max",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 80
            }
          ]
        },
        {
          "title": "cpu",
          "aggregate_func": "max",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-07-2022 00:00" }},
              "value": 70
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get not exist perf data request should return empty response
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
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
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": []
    }
    """

  @concurrent
  Scenario: given get not exist aggregate func request should return error
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
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
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }}
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
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%"
        }
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "hour",
      "from": {{ parseTimeTz "01-07-2022 00:00" }},
      "to": {{ parseTimeTz "01-07-2022 00:00" }}
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
