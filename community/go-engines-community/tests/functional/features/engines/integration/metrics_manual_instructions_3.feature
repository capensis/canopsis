Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given simplified manual instruction disabled in pbh and new event should not count assigned instruction metric if entity in pbh
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-third-1-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-third-1-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-third-1-description",
      "enabled": true,
      "disabled_on_pbh": ["test-maintenance-type-to-engine"],
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-1",
          "stop_on_fail": false
        },
        {
          "job": "test-job-to-run-manual-simplified-instruction-2",
          "stop_on_fail": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-third-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-third-1",
      "resource": "test-resource-to-manual-instruction-metrics-third-1",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-manual-metrics-third-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-1"
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
      "event_type": "pbhenter",
      "connector": "test-connector-to-manual-instruction-metrics-third-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-1",
      "component": "test-component-to-manual-instruction-metrics-third-1",
      "resource": "test-resource-to-manual-instruction-metrics-third-1",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-third-1-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-1"
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
      "connector": "test-connector-to-manual-instruction-metrics-third-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-third-1",
      "resource": "test-resource-to-manual-instruction-metrics-third-1",
      "state": 1
    }
    """
    When I wait 3s
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-third-1-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual disabled in pbh instruction and new event should not count assigned metric if entity in pbh
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-third-2-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-2"
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
      "connector": "test-connector-to-manual-instruction-metrics-third-2",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-third-2",
      "resource": "test-resource-to-manual-instruction-metrics-third-2-1",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-manual-metrics-third-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-2"
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
      "event_type": "pbhenter",
      "connector": "test-connector-to-manual-instruction-metrics-third-2",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-2",
      "component": "test-component-to-manual-instruction-metrics-third-2",
      "resource": "test-resource-to-manual-instruction-metrics-third-2-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-third-2",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-third-2",
      "resource": "test-resource-to-manual-instruction-metrics-third-2-1",
      "state": 1
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-third-2-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-third-2-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-third-2-1"
            }
          }
        ]
      ],
      "disabled_on_pbh": ["test-maintenance-type-to-engine"],
      "description": "test-instruction-to-manual-metrics-third-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-1",
          "stop_on_fail": false
        },
        {
          "job": "test-job-to-run-manual-simplified-instruction-2",
          "stop_on_fail": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait 3s
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-third-2-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 0,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction and new event should count assigned instruction metric only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-third-3-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-third-3-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-third-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-1",
          "stop_on_fail": false
        },
        {
          "job": "test-job-to-run-manual-simplified-instruction-2",
          "stop_on_fail": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-third-3-2",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-third-3-name-2",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-third-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-1",
          "stop_on_fail": false
        },
        {
          "job": "test-job-to-run-manual-simplified-instruction-2",
          "stop_on_fail": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-third-3-3",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-third-3-name-3",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-third-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-1",
          "stop_on_fail": false
        },
        {
          "job": "test-job-to-run-manual-simplified-instruction-2",
          "stop_on_fail": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-third-3-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-third-3"
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
      "connector": "test-connector-to-manual-instruction-metrics-third-3",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-third-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-third-3",
      "resource": "test-resource-to-manual-instruction-metrics-third-3",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-third-3-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 1,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-third-3-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 1,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-third-3-3 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 1,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
