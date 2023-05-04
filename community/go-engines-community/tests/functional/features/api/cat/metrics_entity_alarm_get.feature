Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get request should return metrics
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "entity": "test-entity-to-alarm-metrics-get-1",
      "sampling": "hour",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 01:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 02:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 03:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 04:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 05:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 06:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 07:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 08:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 09:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 10:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 11:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 12:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 13:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 14:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 15:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 16:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 17:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 18:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 19:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 20:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 21:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 22:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 23:00" }},
              "value": 0
            }
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request should return metrics with history
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "entity": "test-entity-to-alarm-metrics-get-2",
      "sampling": "hour",
      "with_history": true,
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 01:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 01:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 02:00" }},
              "value": 1,
              "history_timestamp": {{ parseTimeTz "22-11-2021 02:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 03:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 03:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 04:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 04:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 05:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 05:00" }},
              "history_value": 1
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 06:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 06:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 07:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 07:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 08:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 08:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 09:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 09:00" }},
              "history_value": 1
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 10:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 10:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 11:00" }},
              "value": 1,
              "history_timestamp": {{ parseTimeTz "22-11-2021 11:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 12:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 12:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 13:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 13:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 14:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 14:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 15:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 15:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 16:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 16:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 17:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 17:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 18:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 18:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 19:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 19:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 20:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 20:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 21:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 21:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 22:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 22:00" }},
              "history_value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 23:00" }},
              "value": 0,
              "history_timestamp": {{ parseTimeTz "22-11-2021 23:00" }},
              "history_value": 1
            }
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with empty interval should return metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "entity": "test-entity-to-alarm-metrics-get-1",
      "sampling": "day",
      "from": {{ parseTimeTz "06-09-2020 00:00" }},
      "to": {{ parseTimeTz "08-09-2020 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "06-09-2020 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "07-09-2020 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "08-09-2020 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get request with invalid query params should return bad request
    When I am admin
    When I do POST /api/v4/cat/entity-metrics/alarm
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity": "Entity is missing.",
        "from": "From is missing.",
        "parameters": "Parameters is missing.",
        "sampling": "Sampling is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "entity": "not-exist",
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
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
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "entity": "test-entity-to-alarm-metrics-get-1",
      "parameters": [
        {"metric": "not-exist"}
      ],
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
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
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "entity": "test-entity-to-alarm-metrics-get-1",
      "parameters": [
        {"metric": "total_user_activity"}
      ],
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
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
    When I do POST /api/v4/cat/entity-metrics/alarm:
    """json
    {
      "entity": "test-entity-to-alarm-metrics-get-1",
      "sampling": "not-exist",
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do POST /api/v4/cat/entity-metrics/alarm
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/entity-metrics/alarm
    Then the response code should be 403
