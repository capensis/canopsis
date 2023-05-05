Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get request should return aggregated results
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "entity": "test-entity-to-alarm-metrics-get-1",
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
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "value": 1,
          "aggregate_func": "sum"
        },
        {
          "title": "created_alarms",
          "value": 1,
          "aggregate_func": "max"
        },
        {
          "title": "created_alarms",
          "value": 1,
          "aggregate_func": "min"
        },
        {
          "title": "created_alarms",
          "value": 1,
          "aggregate_func": "avg"
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request without sampling should return error
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "entity": "test-entity-to-alarm-metrics-get-1",
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
        "sampling": "Sampling is missing."
      }
    }
    """

  @concurrent
  Scenario: given get request with unsupported metric should return error
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "entity": "test-entity-to-alarm-metrics-get-1",
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
    """json
    {
      "errors": {
        "parameters.0.metric": "Metric is not supported."
      }
    }
    """

  @concurrent
  Scenario: given get request with not found entity should return error
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "entity": "not-found",
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
        "entity": "Entity doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given get request without aggregate function should return error
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "entity": "test-entity-to-alarm-metrics-get-1",
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
  Scenario: given get request with invalid function should return error
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/aggregate:
    """json
    {
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }},
      "sampling": "hour",
      "entity": "test-entity-to-alarm-metrics-get-1",
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
  Scenario: given get request and no auth user should not allow access
    When I do POST /api/v4/cat/entity-metrics/aggregate
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/entity-metrics/aggregate
    Then the response code should be 403
