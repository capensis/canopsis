Feature: Get alarm metrics
  I need to be able to get alarm metrics
  Only admin should be able to get alarm metrics

  @concurrent
  Scenario: given get request with username group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms"
        },
        {
          "metric": "ack_alarms"
        },
        {
          "metric": "average_ack"
        },
        {
          "metric": "max_ack"
        },
        {
          "metric": "min_ack"
        }
      ],
      "criteria": 3,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-user-to-group-metrics-get-1-username",
          "data": [
            [
              {
                "title": "",
                "value": 3
              }
            ],
            [
              {
                "title": "",
                "value": 3
              }
            ],
            [
              {
                "title": "",
                "value": 200
              }
            ],
            [
              {
                "title": "",
                "value": 300
              }
            ],
            [
              {
                "title": "",
                "value": 100
              }
            ]
          ]
        },
        {
          "title": "test-user-to-group-metrics-get-2-username",
          "data": [
            [
              {
                "title": "",
                "value": 1
              }
            ],
            [
              {
                "title": "",
                "value": 1
              }
            ],
            [
              {
                "title": "",
                "value": 400
              }
            ],
            [
              {
                "title": "",
                "value": 400
              }
            ],
            [
              {
                "title": "",
                "value": 400
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with username group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        },
        {
          "metric": "ack_alarms",
          "criteria": 11
        },
        {
          "metric": "average_ack",
          "criteria": 11
        },
        {
          "metric": "max_ack",
          "criteria": 11
        },
        {
          "metric": "min_ack",
          "criteria": 11
        },
        {
          "metric": "ticket_active_alarms",
          "criteria": 11
        }
      ],
      "criteria": 3,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-user-to-group-metrics-get-1-username",
          "data": [
            [
              {
                "title": "value1",
                "value": 2
              },
              {
                "title": "value2",
                "value": 1
              }
            ],
            [
              {
                "title": "value1",
                "value": 2
              },
              {
                "title": "value2",
                "value": 1
              }
            ],
            [
              {
                "title": "value1",
                "value": 150
              },
              {
                "title": "value2",
                "value": 300
              }
            ],
            [
              {
                "title": "value1",
                "value": 200
              },
              {
                "title": "value2",
                "value": 300
              }
            ],
            [
              {
                "title": "value1",
                "value": 100
              },
              {
                "title": "value2",
                "value": 300
              }
            ],
            []
          ]
        },
        {
          "title": "test-user-to-group-metrics-get-2-username",
          "data": [
            [
              {
                "title": "value2",
                "value": 1
              }
            ],
            [
              {
                "title": "value2",
                "value": 1
              }
            ],
            [
              {
                "title": "value2",
                "value": 400
              }
            ],
            [
              {
                "title": "value2",
                "value": 400
              }
            ],
            [
              {
                "title": "value2",
                "value": 400
              }
            ],
            []
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with username group and different subgroups should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        },
        {
          "metric": "average_ack",
          "criteria": 15
        },
        {
          "metric": "max_ack",
          "criteria": 16
        },
        {
          "metric": "min_ack"
        },
        {
          "metric": "cancel_ack_alarms"
        }
      ],
      "criteria": 3,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-user-to-group-metrics-get-1-username",
          "data": [
            [
              {
                "title": "value1",
                "value": 2
              },
              {
                "title": "value2",
                "value": 1
              }
            ],
            [
              {
                "title": "value1",
                "value": 150
              }
            ],
            [],
            [
              {
                "title": "",
                "value": 100
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ]
          ]
        },
        {
          "title": "test-user-to-group-metrics-get-2-username",
          "data": [
            [
              {
                "title": "value2",
                "value": 1
              }
            ],
            [],
            [],
            [
              {
                "title": "",
                "value": 400
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with user role group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms"
        }
      ],
      "criteria": 12,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-role-to-group-metrics-get-1",
          "data": [
            [
              {
                "title": "",
                "value": 3
              }
            ]
          ]
        },
        {
          "title": "test-role-to-group-metrics-get-2",
          "data": [
            [
              {
                "title": "",
                "value": 1
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with user role group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        }
      ],
      "criteria": 12,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-role-to-group-metrics-get-1",
          "data": [
            [
              {
                "title": "value1",
                "value": 2
              },
              {
                "title": "value2",
                "value": 1
              }
            ]
          ]
        },
        {
          "title": "test-role-to-group-metrics-get-2",
          "data": [
            [
              {
                "title": "value2",
                "value": 1
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with category group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms"
        },
        {
          "metric": "active_alarms"
        },
        {
          "metric": "ack_active_alarms"
        },
        {
          "metric": "not_acked_alarms"
        }
      ],
      "criteria": 13,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-category-to-group-metrics-get-1",
          "data": [
            [
              {
                "title": "",
                "value": 0.75
              }
            ],
            [
              {
                "title": "",
                "value": 1.25
              }
            ],
            [
              {
                "title": "",
                "value": 0.75
              }
            ],
            [
              {
                "title": "",
                "value": 0.5
              }
            ]
          ]
        },
        {
          "title": "test-category-to-group-metrics-get-2",
          "data": [
            [
              {
                "title": "",
                "value": 0.25
              }
            ],
            [
              {
                "title": "",
                "value": 0.5
              }
            ],
            [
              {
                "title": "",
                "value": 0.5
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with category group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        },
        {
          "metric": "active_alarms",
          "criteria": 11
        },
        {
          "metric": "ack_active_alarms",
          "criteria": 11
        },
        {
          "metric": "not_acked_alarms",
          "criteria": 11
        }
      ],
      "criteria": 13,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test-category-to-group-metrics-get-1",
          "data": [
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 1
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0
              }
            ]
          ]
        },
        {
          "title": "test-category-to-group-metrics-get-2",
          "data": [
            [
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value2",
                "value": 0.5
              }
            ],
            [
              {
                "title": "value2",
                "value": 0.5
              }
            ],
            [
              {
                "title": "value2",
                "value": 0
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with impact level group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms"
        }
      ],
      "criteria": 14,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "1",
          "data": [
            [
              {
                "title": "",
                "value": 0.75
              }
            ]
          ]
        },
        {
          "title": "2",
          "data": [
            [
              {
                "title": "",
                "value": 0.25
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with impact level group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        }
      ],
      "criteria": 14,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "1",
          "data": [
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ]
          ]
        },
        {
          "title": "2",
          "data": [
            [
              {
                "title": "value2",
                "value": 0.25
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with entity patterns group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms"
        },
        {
          "metric": "ack_alarms"
        },
        {
          "metric": "active_alarms"
        },
        {
          "metric": "ack_active_alarms"
        },
        {
          "metric": "not_acked_alarms"
        },
        {
          "metric": "average_ack"
        },
        {
          "metric": "max_ack"
        },
        {
          "metric": "min_ack"
        }
      ],
      "entity_patterns": [
        {
          "title": "test2",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "is_one_of",
                  "value": [
                    "test-entity-to-group-metrics-get-1",
                    "test-entity-to-group-metrics-get-2",
                    "test-entity-to-group-metrics-get-3",
                    "test-entity-to-group-metrics-get-5"
                  ]
                }
              }
            ]
          ]
        },
        {
          "title": "test1",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-entity-to-group-metrics-get-4"
                }
              }
            ]
          ]
        }
      ],
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test2",
          "data": [
            [
              {
                "title": "",
                "value": 0.75
              }
            ],
            [
              {
                "title": "",
                "value": 0.75
              }
            ],
            [
              {
                "title": "",
                "value": 1.25
              }
            ],
            [
              {
                "title": "",
                "value": 0.75
              }
            ],
            [
              {
                "title": "",
                "value": 0.5
              }
            ],
            [
              {
                "title": "",
                "value": 200
              }
            ],
            [
              {
                "title": "",
                "value": 300
              }
            ],
            [
              {
                "title": "",
                "value": 100
              }
            ]
          ]
        },
        {
          "title": "test1",
          "data": [
            [
              {
                "title": "",
                "value": 0.25
              }
            ],
            [
              {
                "title": "",
                "value": 0.25
              }
            ],
            [
              {
                "title": "",
                "value": 0.5
              }
            ],
            [
              {
                "title": "",
                "value": 0.5
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ],
            [
              {
                "title": "",
                "value": 400
              }
            ],
            [
              {
                "title": "",
                "value": 400
              }
            ],
            [
              {
                "title": "",
                "value": 400
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with entity patterns group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        },
        {
          "metric": "active_alarms",
          "criteria": 11
        },
        {
          "metric": "ack_active_alarms",
          "criteria": 11
        },
        {
          "metric": "not_acked_alarms",
          "criteria": 11
        },
        {
          "metric": "average_ack",
          "criteria": 11
        },
        {
          "metric": "max_ack",
          "criteria": 11
        },
        {
          "metric": "min_ack",
          "criteria": 11
        }
      ],
      "entity_patterns": [
        {
          "title": "test2",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "is_one_of",
                  "value": [
                    "test-entity-to-group-metrics-get-1",
                    "test-entity-to-group-metrics-get-2",
                    "test-entity-to-group-metrics-get-3",
                    "test-entity-to-group-metrics-get-5"
                  ]
                }
              }
            ]
          ]
        },
        {
          "title": "test1",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-entity-to-group-metrics-get-4"
                }
              }
            ]
          ]
        }
      ],
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test2",
          "data": [
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 1
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0
              }
            ],
            [
              {
                "title": "value1",
                "value": 150
              },
              {
                "title": "value2",
                "value": 300
              }
            ],
            [
              {
                "title": "value1",
                "value": 200
              },
              {
                "title": "value2",
                "value": 300
              }
            ],
            [
              {
                "title": "value1",
                "value": 100
              },
              {
                "title": "value2",
                "value": 300
              }
            ]
          ]
        },
        {
          "title": "test1",
          "data": [
            [
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value2",
                "value": 0.5
              }
            ],
            [
              {
                "title": "value2",
                "value": 0.5
              }
            ],
            [
              {
                "title": "value2",
                "value": 0
              }
            ],
            [
              {
                "title": "value2",
                "value": 400
              }
            ],
            [
              {
                "title": "value2",
                "value": 400
              }
            ],
            [
              {
                "title": "value2",
                "value": 400
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request with entity patterns group and different subgroups should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 11
        },
        {
          "metric": "active_alarms",
          "criteria": 15
        },
        {
          "metric": "ack_active_alarms",
          "criteria": 16
        },
        {
          "metric": "not_acked_alarms"
        },
        {
          "metric": "cancel_ack_alarms"
        }
      ],
      "entity_patterns": [
        {
          "title": "test2",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "is_one_of",
                  "value": [
                    "test-entity-to-group-metrics-get-1",
                    "test-entity-to-group-metrics-get-2",
                    "test-entity-to-group-metrics-get-3",
                    "test-entity-to-group-metrics-get-5"
                  ]
                }
              }
            ]
          ]
        },
        {
          "title": "test1",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-entity-to-group-metrics-get-4"
                }
              }
            ]
          ]
        },
        {
          "title": "test3",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-entity-to-group-metrics-get-1000"
                }
              }
            ]
          ]
        }
      ],
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "title": "test2",
          "data": [
            [
              {
                "title": "value1",
                "value": 0.5
              },
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [
              {
                "title": "value1",
                "value": 1
              }
            ],
            [],
            [
              {
                "title": "",
                "value": 0.5
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ]
          ]
        },
        {
          "title": "test1",
          "data": [
            [
              {
                "title": "value2",
                "value": 0.25
              }
            ],
            [],
            [],
            [
              {
                "title": "",
                "value": 0
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ]
          ]
        },
        {
          "title": "test3",
          "data": [
            [],
            [],
            [],
            [
              {
                "title": "",
                "value": 0
              }
            ],
            [
              {
                "title": "",
                "value": 0
              }
            ]
          ]
        }
      ],
      "meta": {
        "min_date": {{ parseTimeTz "01-07-2021 00:00" }}
      }
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics/group
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics/group
    Then the response code should be 403

  @concurrent
  Scenario: given invalid get request should return error
    When I am admin
    When I do POST /api/v4/cat/metrics/group:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "from": "From is missing.",
        "to": "To is missing.",
        "criteria": "Criteria is missing.",
        "entity_patterns": "EntityPatterns is missing.",
        "parameters": "Parameters is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 11,
      "entity_patterns": [
        {
          "title": "test",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test"
                }
              }
            ]
          ]
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_patterns": "Can't be present both EntityPatterns and Criteria.",
        "parameters": "Parameters is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 2000
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "criteria": "Criteria doesn't exist.",
        "parameters": "Parameters is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "entity_patterns": [
        {
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_patterns.0.title": "Title is missing.",
        "entity_patterns.0.pattern": "Pattern is missing.",
        "parameters": "Parameters is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "entity_patterns": [
        {
          "title": "test",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "is_one_of",
                  "value": "test"
                }
              }
            ]
          ]
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_patterns.0.pattern": "Pattern is invalid entity pattern.",
        "parameters": "Parameters is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 11,
      "parameters": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters": "Parameters should not be blank."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 11,
      "parameters": [
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.metric": "Metric is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 11,
      "parameters": [
        {
          "metric": "test"
        }
      ]
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
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 11,
      "parameters": [
        {
          "metric": "ack_alarms",
          "criteria": 2000
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.0.criteria": "Criteria doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 3,
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
        "criteria": "Criteria doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 11,
      "parameters": [
        {
          "criteria": 3,
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
        "parameters.0.criteria": "Criteria doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/metrics/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "entity_patterns": [
        {
          "title": "test",
          "pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test"
                }
              }
            ]
          ]
        }
      ],
      "parameters": [
        {
          "metric": "total_user_activity"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_patterns": "EntityPatterns is not empty."
      }
    }
    """
