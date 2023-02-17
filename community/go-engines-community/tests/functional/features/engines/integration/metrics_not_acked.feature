Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given alarm should add not_acked_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-not-acked-metrics-axe",
      "connector_name" : "test-connector-name-not-acked-metrics-axe",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-not-acked-metrics-axe",
      "resource" : "test-resource-not-acked-metrics-axe-1",
      "state" : 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_alarms&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm should add not_acked_in_hour_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm should add not_acked_in_four_hours_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm should add not_acked_in_day_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm should decrease not_acked_in_hour_alarms and increase not_acked_in_four_hours_alarms
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm should decrease not_acked_in_four_hours_alarms and increase not_acked_in_day_alarms
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given ack event should remove not_acked_in_hour_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-7-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-not-acked-metrics-axe",
      "connector_name": "test-connector-name-not-acked-metrics-axe",
      "source_type": "resource",
      "component":  "test-component-not-acked-metrics-axe",
      "resource": "test-resource-not-acked-metrics-axe-7",
      "output": "test-output",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm resolve should remove not_acked_in_hour_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-8-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-8"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-not-acked-metrics-axe-8",
      "name": "test-resolve-rule-not-acked-metrics-axe-8-name",
      "description": "test-resolve-rule-not-acked-metrics-axe-8-desc",
      "entity_pattern":[
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-8"
            }
          }
        ]
      ],
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "priority": 1
    }
    """
    Then the response code should be 201
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-not-acked-metrics-axe",
      "connector_name" : "test-connector-name-not-acked-metrics-axe",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-not-acked-metrics-axe",
      "resource" : "test-resource-not-acked-metrics-axe-8",
      "state" : 0
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "resolve_close",
      "connector" : "test-connector-not-acked-metrics-axe",
      "connector_name" : "test-connector-name-not-acked-metrics-axe",
      "source_type" : "resource",
      "component" : "test-component-not-acked-metrics-axe",
      "resource" : "test-resource-not-acked-metrics-axe-8"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm pbhenter should remove not_acked_in_hour_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-9-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name":" test-pbehavior-not-acked-metrics-axe-9",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "6s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhenter",
      "connector" : "test-connector-not-acked-metrics-axe",
      "connector_name" : "test-connector-name-not-acked-metrics-axe",
      "source_type" : "resource",
      "component" : "test-component-not-acked-metrics-axe",
      "resource" : "test-resource-not-acked-metrics-axe-9"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm pbhenter with active type should not remove not_acked_in_hour_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-not-acked-metrics-axe-10-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-10"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name":" test-pbehavior-not-acked-metrics-axe-10",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "6s" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-not-acked-metrics-axe-10"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhenter",
      "connector" : "test-connector-not-acked-metrics-axe",
      "connector_name" : "test-connector-name-not-acked-metrics-axe",
      "source_type" : "resource",
      "component" : "test-component-not-acked-metrics-axe",
      "resource" : "test-resource-not-acked-metrics-axe-10"
    }
    """
    When I wait the next periodical process
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=not_acked_in_hour_alarms&parameters[]=not_acked_in_four_hours_alarms&parameters[]=not_acked_in_day_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "not_acked_in_hour_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "not_acked_in_four_hours_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
        {
          "title": "not_acked_in_day_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
