Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given auto instruction should count executed instruction metric
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "priority": 60,
      "name": "test-instruction-to-auto-metrics-1-1-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-auto-instruction-metrics-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-auto-metrics-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": [
        {
          "type": "create"
        },
        {
          "type": "stateinc"
        }
      ],
      "priority": 61,
      "name": "test-instruction-to-auto-metrics-1-2-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-auto-instruction-metrics-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-auto-metrics-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId2={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": [
        {
          "type": "stateinc"
        }
      ],
      "priority": 62,
      "name": "test-instruction-to-auto-metrics-1-3-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-auto-instruction-metrics-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-auto-metrics-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId3={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-to-auto-instruction-metrics-1",
      "connector_name": "test-connector-name-to-auto-instruction-metrics-1",
      "component": "test-component-to-auto-instruction-metrics-1",
      "resource": "test-resource-to-auto-instruction-metrics-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId1 }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId2 }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-to-auto-instruction-metrics-1",
      "connector_name": "test-connector-name-to-auto-instruction-metrics-1",
      "component": "test-component-to-auto-instruction-metrics-1",
      "resource": "test-resource-to-auto-instruction-metrics-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId3 }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId1 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """

  @concurrent
  Scenario: given auto instruction and fixed alarm should count assigned instruction metric
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "priority": 63,
      "name": "test-instruction-to-auto-metrics-2-1-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-auto-instruction-metrics-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-auto-metrics-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-5"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "priority": 64,
      "name": "test-instruction-to-auto-metrics-2-2-name",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-auto-instruction-metrics-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-auto-metrics-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId2={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-to-auto-instruction-metrics-2",
      "connector_name": "test-connector-name-to-auto-instruction-metrics-2",
      "component": "test-component-to-auto-instruction-metrics-2",
      "resource": "test-resource-to-auto-instruction-metrics-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "connector": "test-connector-to-auto-instruction-metrics-2",
      "connector_name": "test-connector-name-to-auto-instruction-metrics-2",
      "component": "test-component-to-auto-instruction-metrics-2",
      "resource": "test-resource-to-auto-instruction-metrics-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-auto-instruction-metrics-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId1 }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 1,
          "ratio": 100
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/remediation?sampling=day&from={{ nowDateTz }}&to={{ nowDateTz }}&instruction={{ .instructionId2 }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDateTz }},
          "assigned": 1,
          "executed": 0,
          "ratio": 0
        }
      ]
    }
    """
