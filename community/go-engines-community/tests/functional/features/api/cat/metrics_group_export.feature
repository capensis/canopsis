Feature: Get alarm metrics
  I need to be able to export alarm metrics
  Only admin should be able to export alarm metrics

  @concurrent
  Scenario: given export request with username group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
    """json
    {
      "parameters": [
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
    metric,username,value
    ack_alarms,test-user-to-group-metrics-get-1-username,0.75
    ack_alarms,test-user-to-group-metrics-get-2-username,0.25
    average_ack,test-user-to-group-metrics-get-1-username,200
    average_ack,test-user-to-group-metrics-get-2-username,400
    max_ack,test-user-to-group-metrics-get-1-username,300
    max_ack,test-user-to-group-metrics-get-2-username,400
    min_ack,test-user-to-group-metrics-get-1-username,100
    min_ack,test-user-to-group-metrics-get-2-username,400

    """

  @concurrent
  Scenario: given export request with username group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
    """json
    {
      "parameters": [
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
        }
      ],
      "criteria": 3,
      "filter": "test-kpi-filter-to-group-metrics-get",
      "from": {{ parseTimeTz "23-11-2021 00:00" }},
      "to": {{ parseTimeTz "26-11-2021 00:00" }}
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
    infos.test-info-to-group-metrics-get,metric,username,value
    value1,ack_alarms,test-user-to-group-metrics-get-1-username,0.5
    value2,ack_alarms,test-user-to-group-metrics-get-1-username,0.25
    value2,ack_alarms,test-user-to-group-metrics-get-2-username,0.25
    value1,average_ack,test-user-to-group-metrics-get-1-username,150
    value2,average_ack,test-user-to-group-metrics-get-1-username,300
    value2,average_ack,test-user-to-group-metrics-get-2-username,400
    value1,max_ack,test-user-to-group-metrics-get-1-username,200
    value2,max_ack,test-user-to-group-metrics-get-1-username,300
    value2,max_ack,test-user-to-group-metrics-get-2-username,400
    value1,min_ack,test-user-to-group-metrics-get-1-username,100
    value2,min_ack,test-user-to-group-metrics-get-1-username,300
    value2,min_ack,test-user-to-group-metrics-get-2-username,400

    """

  @concurrent
  Scenario: given export request with username group and different subgroups should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
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
    infos.test-info-to-group-metrics-get,infos.test-info-to-group-metrics-get-2,infos.test-info-to-group-metrics-get-3,metric,username,value
    value1,,,ack_alarms,test-user-to-group-metrics-get-1-username,0.5
    value2,,,ack_alarms,test-user-to-group-metrics-get-1-username,0.25
    value2,,,ack_alarms,test-user-to-group-metrics-get-2-username,0.25
    ,value1,,average_ack,test-user-to-group-metrics-get-1-username,150
    ,,,min_ack,test-user-to-group-metrics-get-1-username,100
    ,,,min_ack,test-user-to-group-metrics-get-2-username,400

    """

  @concurrent
  Scenario: given export request with entity patterns group should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
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
    metric,pattern,value
    ack_alarms,test2,0.75
    ack_alarms,test1,0.25
    active_alarms,test2,1.25
    active_alarms,test1,0.5
    ack_active_alarms,test2,0.75
    ack_active_alarms,test1,0.5
    not_acked_alarms,test2,0.5
    not_acked_alarms,test1,0
    average_ack,test2,200
    average_ack,test1,400
    max_ack,test2,300
    max_ack,test1,400
    min_ack,test2,100
    min_ack,test1,400

    """

  @concurrent
  Scenario: given export request with entity patterns group and subgroup should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
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
    infos.test-info-to-group-metrics-get,metric,pattern,value
    value1,ack_alarms,test2,0.5
    value2,ack_alarms,test2,0.25
    value2,ack_alarms,test1,0.25
    value1,active_alarms,test2,1
    value2,active_alarms,test2,0.25
    value2,active_alarms,test1,0.5
    value1,ack_active_alarms,test2,0.5
    value2,ack_active_alarms,test2,0.25
    value2,ack_active_alarms,test1,0.5
    value1,not_acked_alarms,test2,0.5
    value2,not_acked_alarms,test2,0
    value2,not_acked_alarms,test1,0
    value1,average_ack,test2,150
    value2,average_ack,test2,300
    value2,average_ack,test1,400
    value1,max_ack,test2,200
    value2,max_ack,test2,300
    value2,max_ack,test1,400
    value1,min_ack,test2,100
    value2,min_ack,test2,300
    value2,min_ack,test1,400

    """

  @concurrent
  Scenario: given export request with entity patterns group and different subgroups should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
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
    infos.test-info-to-group-metrics-get,infos.test-info-to-group-metrics-get-2,infos.test-info-to-group-metrics-get-3,metric,pattern,value
    value1,,,ack_alarms,test2,0.5
    value2,,,ack_alarms,test2,0.25
    value2,,,ack_alarms,test1,0.25
    ,value1,,active_alarms,test2,1
    ,,,not_acked_alarms,test2,0.5
    ,,,not_acked_alarms,test1,0
    ,,,not_acked_alarms,test3,0
    ,,,cancel_ack_alarms,test2,0
    ,,,cancel_ack_alarms,test1,0
    ,,,cancel_ack_alarms,test3,0

    """

  @concurrent
  Scenario: given export request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics-export/group
    Then the response code should be 401

  @concurrent
  Scenario: given export request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics-export/group
    Then the response code should be 403

  @concurrent
  Scenario: given invalid export request should return error
    When I am admin
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
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
    When I do POST /api/v4/cat/metrics-export/group:
    """json
    {
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }},
      "criteria": 3,
      "filter": "test-kpi-filter-to-group-metrics-get",
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
        "filter": "KpiFilter is not empty."
      }
    }
    """
