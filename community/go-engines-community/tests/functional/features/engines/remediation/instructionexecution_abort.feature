Feature: abort a instruction execution
  I need to be able to abort instruction operation
  Only admin should be able to abort a instruction

  Scenario: given running instruction and alarm in ok state should not allow pause execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-remediation-instruction-execution-abort-1"
            }
          }
        ]
      ],
      "description": "test-remediation-instruction-execution-abort-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-1-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-1-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-1-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-1-step-1-endpoint"
        },
        {
          "name": "test-remediation-instruction-execution-abort-1-step-2",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-1-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-1",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-1",
      "resource": "test-resource-remediation-instruction-execution-abort-1",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-1",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-1",
      "resource": "test-resource-remediation-instruction-execution-abort-1",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/pause
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "execution cannot be paused for closed alarm"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/cancel
    Then the response code should be 204
    When I wait the end of event processing

  Scenario: given paused instruction should cancel it on ok check event
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-remediation-instruction-execution-abort-3"
            }
          }
        ]
      ],
      "description": "test-remediation-instruction-execution-abort-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-3-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-3-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-3-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-3-step-1-endpoint"
        },
        {
          "name": "test-remediation-instruction-execution-abort-3-step-2",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-3-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-3",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-3",
      "resource": "test-resource-remediation-instruction-execution-abort-3",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/pause
    Then the response code should be 204
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-3",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-3",
      "resource": "test-resource-remediation-instruction-execution-abort-3",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "status": 3
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-3
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": true,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "instructionstart",
                "m": "Instruction test-remediation-instruction-execution-abort-3-name."
              },
              {
                "_t": "instructionpause",
                "m": "Instruction test-remediation-instruction-execution-abort-3-name."
              },
              {
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
                "val": 0
              },
              {
                "_t": "instructionabort",
                "m": "Instruction test-remediation-instruction-execution-abort-3-name."
              }
            ]
          }
        }
      }
    ]
    """

  Scenario: given paused instruction should cancel it on resolve
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-remediation-instruction-execution-abort-4"
            }
          }
        ]
      ],
      "description": "test-remediation-instruction-execution-abort-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-4-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-4-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-4-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-4-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-4-step-1-endpoint"
        },
        {
          "name": "test-remediation-instruction-execution-abort-4-step-2",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-4-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-4-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-4-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-4",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-4",
      "resource": "test-resource-remediation-instruction-execution-abort-4",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/pause
    Then the response code should be 204
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-4",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-4",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-remediation-instruction-execution-abort-4",
      "resource": "test-resource-remediation-instruction-execution-abort-4"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-4",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-4",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-remediation-instruction-execution-abort-4",
      "resource": "test-resource-remediation-instruction-execution-abort-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "status": 3
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-4
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "instructionstart",
                "m": "Instruction test-remediation-instruction-execution-abort-4-name."
              },
              {
                "_t": "instructionpause",
                "m": "Instruction test-remediation-instruction-execution-abort-4-name."
              },
              {
                "_t": "cancel"
              },
              {
                "_t": "statusinc",
                "val": 4
              }
            ]
          }
        }
      }
    ]
    """

  Scenario: given paused instruction should cancel it on ok change state action
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-remediation-instruction-execution-abort-5"
            }
          }
        ]
      ],
      "description": "test-remediation-instruction-execution-abort-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-5-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-5-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-5-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-5-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-5-step-1-endpoint"
        },
        {
          "name": "test-remediation-instruction-execution-abort-5-step-2",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-5-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-5-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-5-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-5",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-5",
      "resource": "test-resource-remediation-instruction-execution-abort-5",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/pause
    Then the response code should be 204
    When I wait the end of event processing
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-remediation-instruction-execution-abort-5-name",
      "priority": 10067,
      "enabled": true,
      "triggers": ["ack"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-remediation-instruction-execution-abort-5"
                }
              }
            ]
          ],
          "type": "changestate",
          "parameters": {
            "output": "test-output-remediation-instruction-execution-abort-5",
            "state": 0
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-5",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-5",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-remediation-instruction-execution-abort-5",
      "resource": "test-resource-remediation-instruction-execution-abort-5"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "status": 3
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-5
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": true,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "instructionstart",
                "m": "Instruction test-remediation-instruction-execution-abort-5-name."
              },
              {
                "_t": "instructionpause",
                "m": "Instruction test-remediation-instruction-execution-abort-5-name."
              },
              {
                "_t": "ack"
              },
              {
                "_t": "changestate",
                "val": 0
              },
              {
                "_t": "statusdec",
                "val": 0
              },
              {
                "_t": "instructionabort",
                "m": "Instruction test-remediation-instruction-execution-abort-5-name."
              }
            ]
          }
        }
      }
    ]
    """
