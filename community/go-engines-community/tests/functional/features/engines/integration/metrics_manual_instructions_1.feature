Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given manual instruction and new event should count assigned instruction metric only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-1-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-1-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-1-2",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-1-name-2",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-1-3",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-1-name-3",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-1-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-1"
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
      "connector": "test-connector-to-manual-instruction-metrics-1",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-1",
      "resource": "test-resource-to-manual-instruction-metrics-1",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-1-1 until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-1-2 until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-1-3 until response code is 200 and body contains:
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

  @concurrent
  Scenario: given manual instruction and new event should increase assigned metric on instruction creation only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-2-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-2"
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
        "connector": "test-connector-to-manual-instruction-metrics-2",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-2",
        "resource": "test-resource-to-manual-instruction-metrics-2-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-2",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-2",
        "resource": "test-resource-to-manual-instruction-metrics-2-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-2",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-2",
        "resource": "test-resource-to-manual-instruction-metrics-2-3",
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
      "_id": "test-instruction-to-manual-metrics-2-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-2-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-2-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-2-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-2-step-1-endpoint"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-2-1 until response code is 200 and body contains:
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
      "_id": "test-instruction-to-manual-metrics-2-2",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-2-name-2",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-2-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-2-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-2-step-1-endpoint"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-2-1 until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-2-2 until response code is 200 and body contains:
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
  Scenario: given manual instruction and new events should decrease assigned instruction alarms metric on resolve
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-3-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-3-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-3-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-3-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-3"
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
        "connector": "test-connector-to-manual-instruction-metrics-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-3",
        "resource": "test-resource-to-manual-instruction-metrics-3-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-3",
        "resource": "test-resource-to-manual-instruction-metrics-3-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-3",
        "resource": "test-resource-to-manual-instruction-metrics-3-3",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-3-1 until response code is 200 and body contains:
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
      "_id": "test-resolve-rule-to-manual-instruction-metrics-3",
      "name": "test-resolve-rule-to-manual-instruction-metrics-3-name",
      "description": "test-resolve-rule-to-manual-instruction-metrics-3-desc",
      "entity_pattern":[
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-3"
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
      "connector": "test-connector-to-manual-instruction-metrics-3",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-3",
      "resource": "test-resource-to-manual-instruction-metrics-3-3",
      "state": 0
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-to-manual-instruction-metrics-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-3",
        "component": "test-component-to-manual-instruction-metrics-3",
        "resource": "test-resource-to-manual-instruction-metrics-3-3",
        "source_type": "resource"
      },
      {
        "event_type": "resolve_close",
        "connector": "test-connector-to-manual-instruction-metrics-3",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-3",
        "component": "test-component-to-manual-instruction-metrics-3",
        "resource": "test-resource-to-manual-instruction-metrics-3-3",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-3-1 until response code is 200 and body contains:
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
  Scenario: given manual instruction and new event should count executed instruction metric only one time
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-4-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-4-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-4"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-4-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-4-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-4-2",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-4-name-2",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-4"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-4-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-4-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-4-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-4"
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
      "connector": "test-connector-to-manual-instruction-metrics-4",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-4",
      "resource": "test-resource-to-manual-instruction-metrics-4-1",
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
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-4-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-4-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-manual-instruction-metrics-4",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-4",
      "component": "test-component-to-manual-instruction-metrics-4",
      "resource": "test-resource-to-manual-instruction-metrics-4-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-manual-instruction-metrics-4",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-4",
      "component": "test-component-to-manual-instruction-metrics-4",
      "resource": "test-resource-to-manual-instruction-metrics-4-1",
      "source_type": "resource"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-4-1 until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-4-2 until response code is 200 and body contains:
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
      "instruction": "test-instruction-to-manual-metrics-4-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-manual-instruction-metrics-4",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-4",
      "component": "test-component-to-manual-instruction-metrics-4",
      "resource": "test-resource-to-manual-instruction-metrics-4-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-manual-instruction-metrics-4",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-4",
      "component": "test-component-to-manual-instruction-metrics-4",
      "resource": "test-resource-to-manual-instruction-metrics-4-1",
      "source_type": "resource"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-4-1 until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-4-2 until response code is 200 and body contains:
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
  Scenario: given manual instruction and new events should decrease executed instruction metric alarms if alarm is resolved
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-5-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-5-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-5"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-5-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-5-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-5-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-5"
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
        "connector": "test-connector-to-manual-instruction-metrics-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-5",
        "resource": "test-resource-to-manual-instruction-metrics-5-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-5",
        "resource": "test-resource-to-manual-instruction-metrics-5-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-5",
        "resource": "test-resource-to-manual-instruction-metrics-5-3",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-5-1 until response code is 200 and body contains:
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
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-5-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-5-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-manual-instruction-metrics-5",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
      "component": "test-component-to-manual-instruction-metrics-5",
      "resource": "test-resource-to-manual-instruction-metrics-5-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-manual-instruction-metrics-5",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
      "component": "test-component-to-manual-instruction-metrics-5",
      "resource": "test-resource-to-manual-instruction-metrics-5-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-manual-instruction-metrics-5-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-manual-metrics-5-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-manual-instruction-metrics-5",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
      "component": "test-component-to-manual-instruction-metrics-5",
      "resource": "test-resource-to-manual-instruction-metrics-5-3",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-manual-instruction-metrics-5",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
      "component": "test-component-to-manual-instruction-metrics-5",
      "resource": "test-resource-to-manual-instruction-metrics-5-3",
      "source_type": "resource"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-5-1 until response code is 200 and body contains:
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
      "_id": "test-resolve-rule-to-manual-instruction-metrics-5",
      "name": "test-resolve-rule-to-manual-instruction-metrics-5-name",
      "description": "test-resolve-rule-to-manual-instruction-metrics-5-desc",
      "entity_pattern":[
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-5"
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
      "connector": "test-connector-to-manual-instruction-metrics-5",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-5",
      "resource": "test-resource-to-manual-instruction-metrics-5-3",
      "state": 0
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-to-manual-instruction-metrics-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
        "component": "test-component-to-manual-instruction-metrics-5",
        "resource": "test-resource-to-manual-instruction-metrics-5-3",
        "source_type": "resource"
      },
      {
        "event_type": "resolve_close",
        "connector": "test-connector-to-manual-instruction-metrics-5",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-5",
        "component": "test-component-to-manual-instruction-metrics-5",
        "resource": "test-resource-to-manual-instruction-metrics-5-3",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-5-1 until response code is 200 and body contains:
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
  Scenario: given manual instruction and new event should increase assigned metric on instruction update
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-to-manual-instruction-metrics-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-6",
        "resource": "test-resource-to-manual-instruction-metrics-6-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-6",
        "resource": "test-resource-to-manual-instruction-metrics-6-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-6",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-6",
        "resource": "test-resource-to-manual-instruction-metrics-6-3",
        "state": 1
      }
    ]
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-6-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-6-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-6-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-6-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-6-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-6-2",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-6-name-2",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-6-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-6-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-6-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-6-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-6"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-6-1 until response code is 200 and body contains:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-6-1:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-manual-metrics-6-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-6-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-6-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-6-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-6-step-1-endpoint"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-6-1 until response code is 200 and body contains:
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
  Scenario: given manual instruction with create approval and new event should increase assigned metric on instruction creation only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-7-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-7"
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
        "connector": "test-connector-to-manual-instruction-metrics-7",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-7",
        "resource": "test-resource-to-manual-instruction-metrics-7-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-7",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-7",
        "resource": "test-resource-to-manual-instruction-metrics-7-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-7",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-7",
        "resource": "test-resource-to-manual-instruction-metrics-7-3",
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
      "_id": "test-instruction-to-manual-metrics-7-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-7-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-7-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-7-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-7-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-7-step-1-endpoint"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-7-1 until response code is 200 and body contains:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-7-1/approval:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-7-1 until response code is 200 and body contains:
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
  Scenario: given manual instruction with update approval and new event should increase assigned metric on instruction creation only one time for alarm metric and every time for instruction metric
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-8-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-8"
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
        "connector": "test-connector-to-manual-instruction-metrics-8",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-8",
        "resource": "test-resource-to-manual-instruction-metrics-8-1",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-8",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-8",
        "resource": "test-resource-to-manual-instruction-metrics-8-2",
        "state": 1
      },
      {
        "connector": "test-connector-to-manual-instruction-metrics-8",
        "connector_name": "test-connector-name-to-manual-instruction-metrics-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-manual-instruction-metrics-8",
        "resource": "test-resource-to-manual-instruction-metrics-8-3",
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
      "_id": "test-instruction-to-manual-metrics-8-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-8-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-8-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-8-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-8-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-8-step-1-endpoint"
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-8-1 until response code is 200 and body contains:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-8-1:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-manual-metrics-8-name-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-8-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-8-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-8-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-8-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-8-step-1-endpoint"
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-manual-metrics-8-1/approval:
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-8-1 until response code is 200 and body contains:
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
  Scenario: given manual instruction disabled in pbh and new event should not count assigned instruction metric if entity in pbh
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-manual-metrics-9-1",
      "type": 0,
      "name": "test-instruction-to-manual-metrics-9-name-1",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-9"
            }
          }
        ]
      ],
      "description": "test-instruction-to-manual-metrics-9-description",
      "enabled": true,
      "disabled_on_pbh": ["test-maintenance-type-to-engine"],
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-manual-metrics-9-step-1",
          "operations": [
            {
              "name": "test-instruction-to-manual-metrics-9-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-manual-metrics-9-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-manual-metrics-9-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-manual-instruction-metrics-9",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-9",
      "resource": "test-resource-to-manual-instruction-metrics-9",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-manual-metrics-9",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-manual-instruction-metrics-9"
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
      "connector": "test-connector-to-manual-instruction-metrics-9",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-9",
      "component": "test-component-to-manual-instruction-metrics-9",
      "resource": "test-resource-to-manual-instruction-metrics-9",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-to-manual-instruction-metrics-9-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-manual-instruction-metrics-9"
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
      "connector": "test-connector-to-manual-instruction-metrics-9",
      "connector_name": "test-connector-name-to-manual-instruction-metrics-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-manual-instruction-metrics-9",
      "resource": "test-resource-to-manual-instruction-metrics-9",
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
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDate }}&to={{ nowDate }}&instruction=test-instruction-to-manual-metrics-9-1 until response code is 200 and body contains:
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
