Feature: Export aggregated metrics
  I need to be able to export aggregated metrics
  Only admin should be able to export aggregated metrics

  @concurrent
  Scenario: given export aggregated metrics request should return aggregated metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        },
        {
          "metric": "created_alarms",
          "aggregate_func": "max"
        },
        {
          "metric": "created_alarms",
          "aggregate_func": "min"
        },
        {
          "metric": "created_alarms",
          "aggregate_func": "avg"
        },
        {
          "metric": "active_alarms",
          "aggregate_func": "max"
        },
        {
          "metric": "active_alarms",
          "aggregate_func": "min"
        },
        {
          "metric": "active_alarms",
          "aggregate_func": "avg"
        },
        {
          "metric": "without_ticket_active_alarms",
          "aggregate_func": "avg"
        },
        {
          "metric": "without_ticket_active_alarms",
          "aggregate_func": "max"
        },
        {
          "metric": "without_ticket_active_alarms",
          "aggregate_func": "min"
        },
        {
          "metric": "ratio_tickets"
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
    sum,created_alarms,3
    max,created_alarms,1
    min,created_alarms,1
    avg,created_alarms,1
    max,active_alarms,6
    min,active_alarms,5
    avg,active_alarms,5.58
    avg,without_ticket_active_alarms,4.54
    max,without_ticket_active_alarms,5
    min,without_ticket_active_alarms,4
    ,ratio_tickets,18.65

    """

  @concurrent
  Scenario: given export aggregated metrics request with empty interval should return aggregated metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "06-09-2020 00:00" }},
      "to": {{ parseTimeTz "08-09-2020 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        },
        {
          "metric": "created_alarms",
          "aggregate_func": "max"
        },
        {
          "metric": "created_alarms",
          "aggregate_func": "min"
        },
        {
          "metric": "created_alarms",
          "aggregate_func": "avg"
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
    max,created_alarms,0
    min,created_alarms,0
    avg,created_alarms,0

    """
    
  @concurrent
  Scenario: given export aggregated metrics request without sampling should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "sampling": "Sampling is missing."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with unsupported metric should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "total_user_activity",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "parameters.0.metric": "Metric is not supported."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with not found kpi filter should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-not-found",
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "filter": "Filter doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with not found widget filter should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "widget_filters": [
        "test-widget-filter-not-found"
      ],
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "widget_filters.0": "WidgetFilter doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with both filters should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "widget_filters": [
        "test-widget-filter-to-alarm-metrics-get-1"
      ],
      "parameters": [
        {
          "metric": "created_alarms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget_filters": "Can't be present both WidgetFilters and KpiFilter."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request without aggregate function should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "active_alarms"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.aggregate_func": "AggregateFunc is missing."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with cumulative number metric with sum aggregate function should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "active_alarms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.aggregate_func": "AggregateFunc must be one of [avg, max, min]."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with subtraction type of metric with sum function should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "without_ticket_active_alarms",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.aggregate_func": "AggregateFunc must be one of [avg, max, min]."
      }
    }
    """

  @concurrent
  Scenario: given export aggregated metrics request with invalid function should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "without_ticket_active_alarms",
          "aggregate_func": "qwe"
        }
      ]
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
  Scenario: given export aggregated metrics request with ratio metric with aggregate function should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {
          "metric": "ratio_tickets",
          "aggregate_func": "sum"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.aggregate_func": "AggregateFunc is not empty."
      }
    }
    """    

  @concurrent
  Scenario: given export aggregated metrics request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics-export/aggregate
    Then the response code should be 401

  @concurrent
  Scenario: given export aggregated metrics request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics-export/aggregate
    Then the response code should be 403
