Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get created_alarms hour request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
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
              "value": 1
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
              "value": 1
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

  Scenario: given get created_alarms hour request should return metrics with history
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
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
              "value": 1,
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
  Scenario: given get created_alarms day request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTimeTz "20-11-2021 00:00" }},
      "to": {{ parseTimeTz "24-11-2021 00:00" }}
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
            "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
            "value": 1
          },
          {
            "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
            "value": 3
          },
          {
            "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
            "value": 3
          },
          {
            "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
            "value": 0
          }
        ]
      }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms week request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "week",
      "from": {{ parseTimeTz "06-09-2021 00:00" }},
      "to": {{ parseTimeTz "10-10-2021 00:00" }}
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
              "timestamp": {{ parseTimeTz "06-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "13-09-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "20-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "27-09-2021 00:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTimeTz "04-10-2021 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms month request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "month",
      "from": {{ parseTimeTz "01-06-2021 00:00" }},
      "to": {{ parseTimeTz "31-10-2021 00:00" }}
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
              "timestamp": {{ parseTimeTz "01-06-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "01-08-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-09-2021 00:00" }},
              "value": 3
            },
            {
              "timestamp": {{ parseTimeTz "01-10-2021 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get created_alarms request with empty interval should return metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
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
  Scenario: given get created_alarms request with filter by entity infos should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get-by-entity-infos",
      "sampling": "day",
      "from": {{ parseTimeTz "20-11-2021 00:00" }},
      "to": {{ parseTimeTz "24-11-2021 00:00" }}
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
            "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
            "value": 0
          },
          {
            "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
            "value": 3
          },
          {
            "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
            "value": 2
          },
          {
            "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
            "value": 0
          }
        ]
      }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms hour request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "active_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
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
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 01:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 02:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 03:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 04:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 05:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 06:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 07:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 08:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 09:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 10:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 11:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 12:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 13:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 14:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 15:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 16:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 17:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 18:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 19:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 20:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 21:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 22:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 23:00" }},
              "value": 6
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms day request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "active_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTimeTz "20-11-2021 00:00" }},
      "to": {{ parseTimeTz "24-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 5
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 6
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 6
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms week request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "active_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "week",
      "from": {{ parseTimeTz "06-09-2021 00:00" }},
      "to": {{ parseTimeTz "10-10-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
        "title": "active_alarms",
        "data": [
          {
            "timestamp": {{ parseTimeTz "06-09-2021 00:00" }},
            "value": 1
          },
          {
            "timestamp": {{ parseTimeTz "13-09-2021 00:00" }},
            "value": 2
          },
          {
            "timestamp": {{ parseTimeTz "20-09-2021 00:00" }},
            "value": 2
          },
          {
            "timestamp": {{ parseTimeTz "27-09-2021 00:00" }},
            "value": 4
          },
          {
            "timestamp": {{ parseTimeTz "04-10-2021 00:00" }},
            "value": 4
          }
        ]
      }
      ]
    }
    """

  @concurrent
  Scenario: given get active_alarms month request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "active_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "month",
      "from": {{ parseTimeTz "01-06-2021 00:00" }},
      "to": {{ parseTimeTz "31-10-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-06-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "01-08-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "01-09-2021 00:00" }},
              "value": 4
            },
            {
              "timestamp": {{ parseTimeTz "01-10-2021 00:00" }},
              "value": 4
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets hour request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "ratio_tickets"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
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
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 01:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 02:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 03:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 04:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 05:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 06:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 07:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 08:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 09:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 10:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 11:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 12:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 13:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 14:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 15:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 16:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 17:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 18:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 19:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 20:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 21:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 22:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 23:00" }},
              "value": 16.66
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets day request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "ratio_tickets"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTimeTz "20-11-2021 00:00" }},
      "to": {{ parseTimeTz "24-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTimeTz "20-11-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTimeTz "21-11-2021 00:00" }},
              "value": 40
            },
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 20
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 16.66
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 16.66
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets week request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "ratio_tickets"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "week",
      "from": {{ parseTimeTz "06-09-2021 00:00" }},
      "to": {{ parseTimeTz "10-10-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTimeTz "06-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "13-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "20-09-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "27-09-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTimeTz "04-10-2021 00:00" }},
              "value": 25
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get ratio_tickets month request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "ratio_tickets"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "sampling": "month",
      "from": {{ parseTimeTz "01-06-2021 00:00" }},
      "to": {{ parseTimeTz "31-10-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTimeTz "01-06-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-07-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-08-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "01-09-2021 00:00" }},
              "value": 25
            },
            {
              "timestamp": {{ parseTimeTz "01-10-2021 00:00" }},
              "value": 25
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given get request with invalid query params should return bad request
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "from": "From is missing.",
        "parameters": "Parameters is missing.",
        "sampling": "Sampling is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "filter": "not-exist",
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
        "filter": "Filter doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
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
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "total_user_activity"}
      ],
      "sampling": "day",
      "from": {{ now }},
      "to": {{ now }}
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
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "sampling": "not-exist",
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "from": {{ now }},
      "to": {{ now }}
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
    When I do POST /api/v4/cat/metrics/alarm
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics/alarm
    Then the response code should be 403

  @concurrent
  Scenario: given get request with all parameters should return all metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"},
        {"metric": "active_alarms"},
        {"metric": "non_displayed_alarms"},
        {"metric": "instruction_alarms"},
        {"metric": "pbehavior_alarms"},
        {"metric": "correlation_alarms"},
        {"metric": "ack_alarms"},
        {"metric": "cancel_ack_alarms"},
        {"metric": "ack_active_alarms"},
        {"metric": "ticket_active_alarms"},
        {"metric": "without_ticket_active_alarms"},
        {"metric": "ratio_correlation"},
        {"metric": "ratio_instructions"},
        {"metric": "ratio_tickets"},
        {"metric": "ratio_non_displayed"},
        {"metric": "average_ack"},
        {"metric": "average_resolve"}
      ],
      "filter": "test-kpi-filter-to-all-alarm-metrics-get",
      "sampling": "day",
      "from": {{ parseTimeTz "22-11-2021 00:00" }},
      "to": {{ parseTimeTz "24-11-2021 00:00" }}
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
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 3
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 3
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 3
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "pbehavior_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "correlation_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "cancel_ack_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "ack_active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 1
            }
          ]
        },
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 1
            }
          ]
        },
        {
          "title": "without_ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 2
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 2
            }
          ]
        },
        {
          "title": "ratio_correlation",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "ratio_instructions",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "ratio_non_displayed",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 33.33
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 33.33
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 500
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        },
        {
          "title": "average_resolve",
          "data": [
            {
              "timestamp": {{ parseTimeTz "22-11-2021 00:00" }},
              "value": 0
            },
            {
              "timestamp": {{ parseTimeTz "23-11-2021 00:00" }},
              "value": 1000
            },
            {
              "timestamp": {{ parseTimeTz "24-11-2021 00:00" }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given filter with old pattern should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "test-kpi-filter-to-alarm-metrics-get-by-old-pattern",
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
              "value": 1
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
              "value": 1
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
  Scenario: given get created_alarms hour request should return metrics by widget filter
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "widget_filters": ["test-widget-filter-to-alarm-metrics-get-1"],
      "parameters": [
        {"metric": "created_alarms"}
      ],
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
              "value": 1
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
              "value": 1
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
  Scenario: given get created_alarms hour request should return metrics by widget filters
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "widget_filters": [
        "test-widget-filter-to-alarm-metrics-get-2",
        "test-widget-filter-to-alarm-metrics-get-1"
      ],
      "parameters": [
        {"metric": "created_alarms"}
      ],
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
              "value": 1
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
              "value": 1
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
  Scenario: given get created_alarms hour request with both metrics and widget filters should return error
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "widget_filters": ["test-widget-filter-to-alarm-metrics-get"],
      "filter": "test-kpi-filter-to-alarm-metrics-get",
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "sampling": "hour",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }}
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
  Scenario: given get created_alarms hour request with not exist widget filter should return error
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "widget_filters": ["test-widget-filter-not-exist"],
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "sampling": "hour",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }}
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget_filters.0": "WidgetFilter doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given get created_alarms hour request with invalid widget filter should return error
    When I am admin
    When I do POST /api/v4/cat/metrics/alarm:
    """json
    {
      "widget_filters": ["test-widget-filter-to-alarm-metrics-get-3"],
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "sampling": "hour",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "23-11-2021 00:00" }}
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget_filters.0": "WidgetFilter is not applicable."
      }
    }
    """
