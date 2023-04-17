Feature: Perf data should be stored.
  I need to be able to see metrics.

  @concurrent
  Scenario: given check event with perf data should store it
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "perf_data": "cpu=20%;80;90;0;100",
      "connector": "test-connector-metrics-perf-data-1",
      "connector_name": "test-connector-name-metrics-perf-data-1",
      "component": "test-component-metrics-perf-data-1",
      "resource": "test-resource-metrics-perf-data-1",
      "source_type": "resource"
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 20
            }
          ]
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "perf_data": "cpu=60%;80;90;0;100",
      "connector": "test-connector-metrics-perf-data-1",
      "connector_name": "test-connector-name-metrics-perf-data-1",
      "component": "test-component-metrics-perf-data-1",
      "resource": "test-resource-metrics-perf-data-1",
      "source_type": "resource"
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ nowDate }},
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
          "aggregate_func": "last"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 60
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
          "aggregate_func": "avg"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 40
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
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 60
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
          "aggregate_func": "min"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "unit": "%",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 20
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "sum"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "value": 80,
          "aggregate_func": "sum",
          "unit": "%"
        }
      ]
    }
    """
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "last"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "value": 60,
          "aggregate_func": "last",
          "unit": "%"
        }
      ]
    }
    """
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "avg"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "value": 40,
          "aggregate_func": "avg",
          "unit": "%"
        }
      ]
    }
    """
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "max"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "value": 60,
          "aggregate_func": "max",
          "unit": "%"
        }
      ]
    }
    """
    When I do POST /api/v4/cat/metrics/aggregate:
    """json
    {
      "parameters": [
        {
          "metric": "cpu_%",
          "aggregate_func": "min"
        }
      ],
      "entity": "test-resource-metrics-perf-data-1/test-component-metrics-perf-data-1",
      "sampling": "day",
      "from": {{ nowDate }},
      "to": {{ nowDate }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "cpu_%",
          "value": 20,
          "aggregate_func": "min",
          "unit": "%"
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-metrics-perf-data-1&perf_data[]={{ "cpu_%" | query_escape }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "filtered_perf_data": ["cpu_%"]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/entities?search=test-resource-metrics-perf-data-1&perf_data[]={{ "cpu_%" | query_escape }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "filtered_perf_data": ["cpu_%"]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
