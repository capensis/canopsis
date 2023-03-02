Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given manual disabled in pbh instruction and new event should not count assigned metric if entity in pbh
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-second-1-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-1"
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
      "connector": "test-connector-to-manual-instruction-metrics-second-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-second-1",
      "resource": "test-resource-to-manual-instruction-metrics-second-1-1",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-manual-metrics-second-1",
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
              "value": "test-component-to-manual-instruction-metrics-second-1"
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
      "connector": "test-connector-to-manual-instruction-metrics-second-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-1",
      "component": "test-component-to-manual-instruction-metrics-second-1",
      "resource": "test-resource-to-manual-instruction-metrics-second-1-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-second-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-second-1",
      "resource": "test-resource-to-manual-instruction-metrics-second-1-1",
      "state": 1
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-1-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-second-1-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-1-1"
            }
          }
        ]
      ],
      "disabled_on_pbh": ["test-maintenance-type-to-engine"],
      "description": "test-instruction-to-manual-metrics-second-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-second-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-second-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-second-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-second-1-step-1-endpoint"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-1-1 until response code is 200 and body contains:
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
  Scenario: given manual instruction and new events should not count executed instruction metric if an alarm wasn't assigned before
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test manual instruction metrics 11",
      "type": "enrichment",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-2"
            }
          }
        ]
      ],
      "priority": 0,
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "test-infos-to-manual-metrics-second-2",
            "value": "{{ `{{ .Event.Output }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-2-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-second-2-name-1",
      "entity_pattern": [
        [
          {
            "field": "infos.test-infos-to-manual-metrics-second-2",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test value"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-second-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-second-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-second-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-second-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-second-2-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-2"
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
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-2",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-2",
        "resource": "test-resource-to-manual-instruction-metrics-second-2-1",
        "output": "test value",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-2",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-2",
        "resource": "test-resource-to-manual-instruction-metrics-second-2-2",
        "output": "test value 2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-2",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-2",
        "resource": "test-resource-to-manual-instruction-metrics-second-2-3",
        "output": "test value",
        "state": 1
      }
    ]
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
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-2-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-second-2",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-second-2",
      "resource": "test-resource-to-manual-instruction-metrics-second-2-2",
      "output": "test value",
      "state": 1
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-second-2-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-second-2-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-manual-instruction-metrics-second-2",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-2",
      "component": "test-component-to-manual-instruction-metrics-second-2",
      "resource": "test-resource-to-manual-instruction-metrics-second-2-2",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-manual-instruction-metrics-second-2",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-2",
      "component": "test-component-to-manual-instruction-metrics-second-2",
      "resource": "test-resource-to-manual-instruction-metrics-second-2-2",
      "source_type": "resource"
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
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-2-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction and new event should increase assigned metric on instruction creation only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-second-3-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-3"
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
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-3",
        "resource": "test-resource-to-manual-instruction-metrics-second-3-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-3",
        "resource": "test-resource-to-manual-instruction-metrics-second-3-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-3",
        "resource": "test-resource-to-manual-instruction-metrics-second-3-3",
        "state": 1
      }
    ]
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
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-3-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-3-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-3-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-3-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-3-description",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-3-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-3-2",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-3-name-2",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-3-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-3-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-3-description",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 3
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-3-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-3-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction and new events should decrease assigned instruction metric alarms on resolve
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-4-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-4-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-4"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-4-description",
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
      "name": "test-filter-to-manual-instruction-metrics-second-4-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-4"
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
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-4",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-4",
        "resource": "test-resource-to-manual-instruction-metrics-second-4-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-4",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-4",
        "resource": "test-resource-to-manual-instruction-metrics-second-4-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-4",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-4",
        "resource": "test-resource-to-manual-instruction-metrics-second-4-3",
        "state": 1
      }
    ]
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
              "value": 3
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-4-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-to-manual-instruction-metrics-second-4",
      "name": "test-resolve-rule-to-manual-instruction-metrics-second-4-name",
      "description": "test-resolve-rule-to-manual-instruction-metrics-second-4-desc",
      "entity_pattern":[
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-4"
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
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-second-4",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-second-4",
      "resource": "test-resource-to-manual-instruction-metrics-second-4-3",
      "state": 0
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-to-manual-instruction-metrics-second-4",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-4",
        "component": "test-component-to-manual-instruction-metrics-second-4",
        "resource": "test-resource-to-manual-instruction-metrics-second-4-3",
        "source_type": "resource"
      },
      {
        "event_type": "resolve_close",
        "connector": "test-connector-to-manual-instruction-metrics-second-4",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-4",
        "component": "test-component-to-manual-instruction-metrics-second-4",
        "resource": "test-resource-to-manual-instruction-metrics-second-4-3",
        "source_type": "resource"
      }
    ]
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
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-4-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction and new event should count executed instruction metric only one time
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-5-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-5-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-5"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-5-description",
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
      "_id": "test-instruction-to-manual-metrics-second-5-2",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-5-name-2",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-5"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-5-description",
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
      "name": "test-filter-to-manual-instruction-metrics-second-5-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-5"
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
      "connector": "test-connector-to-manual-instruction-metrics-second-5",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-second-5",
      "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ratio_remediated_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ratio_remediated_alarms",
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
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-second-5-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-second-5-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-manual-instruction-metrics-second-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
        "component": "test-component-to-manual-instruction-metrics-second-5",
        "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
        "component": "test-component-to-manual-instruction-metrics-second-5",
        "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
        "component": "test-component-to-manual-instruction-metrics-second-5",
        "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ratio_remediated_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ratio_remediated_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 100
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-5-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-5-2 until response code is 200 and body contains:
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
    When I wait 5s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-second-5-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-manual-instruction-metrics-second-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
        "component": "test-component-to-manual-instruction-metrics-second-5",
        "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
        "component": "test-component-to-manual-instruction-metrics-second-5",
        "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-5",
        "component": "test-component-to-manual-instruction-metrics-second-5",
        "resource": "test-resource-to-manual-instruction-metrics-second-5-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ratio_remediated_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ratio_remediated_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 100
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-5-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-5-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction and new events should decrease executed instruction metrics if alarm is resolved
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-6-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-6-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-6"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-6-description",
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
      "name": "test-filter-to-manual-instruction-metrics-second-6-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-6"
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
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
        "state": 1
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ratio_remediated_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ratio_remediated_alarms",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-6-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-second-6-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-second-6-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-second-6-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-second-6-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ratio_remediated_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ratio_remediated_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 66.66
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-6-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 2,
          "ratio": 66.66
        }
      ]
    }
    """
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-to-manual-instruction-metrics-second-6",
      "name": "test-resolve-rule-to-manual-instruction-metrics-second-6-name",
      "description": "test-resolve-rule-to-manual-instruction-metrics-second-6-desc",
      "entity_pattern":[
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-6"
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
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-second-6",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-second-6",
      "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
      "state": 0
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
        "source_type": "resource"
      },
      {
        "event_type": "resolve_close",
        "connector": "test-connector-to-manual-instruction-metrics-second-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-6",
        "component": "test-component-to-manual-instruction-metrics-second-6",
        "resource": "test-resource-to-manual-instruction-metrics-second-6-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_executed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_executed_alarms",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ratio_remediated_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ratio_remediated_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 50
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-6-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 2,
          "ratio": 66.66
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction and new event should increase assigned metric on instruction update
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-7",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-7",
        "resource": "test-resource-to-manual-instruction-metrics-second-7-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-7",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-7",
        "resource": "test-resource-to-manual-instruction-metrics-second-7-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-7",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-7",
        "resource": "test-resource-to-manual-instruction-metrics-second-7-3",
        "state": 1
      }
    ]
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-7-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-7-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-7-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-7-description",
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
      "_id": "test-instruction-to-manual-metrics-second-7-2",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-7-name-2",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-7-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-7-description",
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
      "name": "test-filter-to-manual-instruction-metrics-second-7-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-7-1 until response code is 200 and body contains:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-second-7-1:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-7-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-7-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-7-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-7-description",
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
    Then the response code should be 200
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 3
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-7-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction with create approval and new event should increase assigned metric on instruction creation only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-second-8-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-8"
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
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-8",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-8",
        "resource": "test-resource-to-manual-instruction-metrics-second-8-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-8",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-8",
        "resource": "test-resource-to-manual-instruction-metrics-second-8-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-8",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-8",
        "resource": "test-resource-to-manual-instruction-metrics-second-8-3",
        "state": 1
      }
    ]
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
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-8-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-8-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-8-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-8-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-8-description",
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
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-8-1 until response code is 200 and body contains:
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
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-second-8-1/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-8-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction with update approval and new event should increase assigned metric on instruction creation only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-second-9-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-second-9"
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
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-9",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-9",
        "resource": "test-resource-to-manual-instruction-metrics-second-9-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-9",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-9",
        "resource": "test-resource-to-manual-instruction-metrics-second-9-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-second-9",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-second-9",
        "resource": "test-resource-to-manual-instruction-metrics-second-9-3",
        "state": 1
      }
    ]
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
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-second-9-1",
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-9-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-9-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-9-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-9-description",
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-9-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 2,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-second-9-1:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-manual-metrics-second-9-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-9-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-9-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-second-9-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-second-9-description",
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
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-second-9-1/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I am admin
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=manual_instruction_assigned_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "manual_instruction_assigned_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 3
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-second-9-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "assigned": 3,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
