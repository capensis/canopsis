Feature: Get metrics rating
  I need to be able to get metrics rating
  Only admin should be able to get metrics rating

  Scenario: given get created_alarms by name request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=created_alarms&criteria=1&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 3
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 1
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  Scenario: given get created_alarms by name with filter by entity infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=created_alarms&criteria=1&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get-by-entity-infos
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get created_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=created_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 3
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get active_alarms by name request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=active_alarms&criteria=1&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "25-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 4
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get non_displayed_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=non_displayed_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get instruction_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=instruction_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get pbehavior_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=pbehavior_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get correlation_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=correlation_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get ack_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ack_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 2
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get cancel_ack_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=cancel_ack_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 2
        }
      ]
    }
    """

  Scenario: given get ack_active_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ack_active_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 0
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get ticket_active_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ticket_active_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 1
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get without_ticket_active_alarms by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=without_ticket_active_alarms&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 3
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 0
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 1
        }
      ]
    }
    """

  Scenario: given get ratio_correlation by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ratio_correlation&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 25
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 100
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 0
        }
      ]
    }
    """

  Scenario: given get ratio_instructions by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ratio_instructions&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 25
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 100
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 0
        }
      ]
    }
    """

  Scenario: given get ratio_tickets by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ratio_tickets&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 25
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 100
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 0
        }
      ]
    }
    """

  Scenario: given get ratio_non_displayed by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=ratio_non_displayed&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 25
        },
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 100
        },
        {
          "label": "test-entity-to-metrics-rating-get-3",
          "value": 0
        }
      ]
    }
    """

  Scenario: given get average_ack by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=average_ack&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 300
        },
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 150
        }
      ]
    }
    """

  Scenario: given get average_resolve by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=average_resolve&criteria=2&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}&filter=test-kpi-filter-to-metrics-rating-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-entity-to-metrics-rating-get-2",
          "value": 300
        },
        {
          "label": "test-entity-to-metrics-rating-get-1",
          "value": 200
        }
      ]
    }
    """

  Scenario: given get total_user_activity by infos request should return metrics
    When I am admin
    When I do GET /api/v4/cat/metrics/rating?metric=total_user_activity&criteria=3&from={{ parseTimeTz "20-11-2021 00:00" }}&to={{ parseTimeTz "23-11-2021 00:00" }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "label": "test-user-to-metrics-rating-get-2",
          "value": 300
        },
        {
          "label": "test-user-to-metrics-rating-get-1",
          "value": 100
        }
      ]
    }
    """

  Scenario: given get request with invalid query params should return bad request
    When I am admin
    When I do GET /api/v4/cat/metrics/rating
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "criteria": "Criteria is missing.",
        "from": "From is missing.",
        "metric": "Metric is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/rating?metric=not-exist&from={{ nowDateTz }}&to={{ nowDateTz }}&criteria=1
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "metric": "Metric doesn't exist."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/rating?criteria=1000000&metric=created_alarms&from={{ nowDateTz }}&to={{ nowDateTz }}
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "criteria": "Criteria doesn't exist."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/rating?filter=not-exist&from={{ nowDateTz }}&to={{ nowDateTz }}&metric=created_alarms&criteria=1
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter doesn't exist."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/rating?metric=created_alarms&criteria=3&from={{ nowDateTz }}&to={{ nowDateTz }}
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "criteria": "Criteria doesn't exist."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/rating?metric=total_user_activity&criteria=1&from={{ nowDateTz }}&to={{ nowDateTz }}
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "criteria": "Criteria doesn't exist."
      }
    }
    """
    When I do GET /api/v4/cat/metrics/rating?metric=total_user_activity&filter=test-kpi-filter-to-metrics-rating-get&criteria=3&from={{ nowDateTz }}&to={{ nowDateTz }}
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter is not empty."
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metrics/rating
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metrics/rating
    Then the response code should be 403
