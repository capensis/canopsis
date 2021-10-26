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
      "entity_patterns": [
        {
          "name": "test-resource-remediation-instruction-execution-abort-1"
        }
      ],
      "description": "test-remediation-instruction-execution-abort-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "seconds": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-1-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-1-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-1-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
              "time_to_complete": {"seconds": 6, "unit":"s"},
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

  Scenario: given running instruction and alarm in ok state should cancel execution after long time of inactivity
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-remediation-instruction-execution-abort-2"
        }
      ],
      "description": "test-remediation-instruction-execution-abort-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "seconds": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-2-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-2-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-2-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-2-step-1-endpoint"
        },
        {
          "name": "test-remediation-instruction-execution-abort-2-step-2",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-2-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-remediation-instruction-execution-abort-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-remediation-instruction-execution-abort-2",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-2",
      "resource": "test-resource-remediation-instruction-execution-abort-2",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-2
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
      "connector": "test-connector-remediation-instruction-execution-abort-2",
      "connector_name": "test-connector-name-remediation-instruction-execution-abort-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-remediation-instruction-execution-abort-2",
      "resource": "test-resource-remediation-instruction-execution-abort-2",
      "state": 0
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 410
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-2&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
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
                "m": "Instruction test-remediation-instruction-execution-abort-2-name."
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
                "m": "Instruction test-remediation-instruction-execution-abort-2-name."
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given paused instruction should cancel it on ok check event
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-3-name",
      "entity_patterns": [
        {
          "name": "test-resource-remediation-instruction-execution-abort-3"
        }
      ],
      "description": "test-remediation-instruction-execution-abort-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "seconds": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-3-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-3-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-3-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-3-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
              "time_to_complete": {"seconds": 6, "unit":"s"},
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
    When I wait 1s
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 410
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-3&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
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
      ]
    }
    """

  Scenario: given paused instruction should cancel it on resolve
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-4-name",
      "entity_patterns": [
        {
          "name": "test-resource-remediation-instruction-execution-abort-4"
        }
      ],
      "description": "test-remediation-instruction-execution-abort-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "seconds": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-4-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-4-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-4-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-4-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
              "time_to_complete": {"seconds": 6, "unit":"s"},
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
    When I wait 1s
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 410
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-4&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
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
      ]
    }
    """

  Scenario: given paused instruction should cancel it on ok change state action
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-remediation-instruction-execution-abort-5-name",
      "entity_patterns": [
        {
          "name": "test-resource-remediation-instruction-execution-abort-5"
        }
      ],
      "description": "test-remediation-instruction-execution-abort-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "seconds": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-remediation-instruction-execution-abort-5-step-1",
          "operations": [
            {
              "name": "test-remediation-instruction-execution-abort-5-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-remediation-instruction-execution-abort-5-step-1-operation-1-description"
            },
            {
              "name": "test-remediation-instruction-execution-abort-5-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
              "time_to_complete": {"seconds": 6, "unit":"s"},
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
      "enabled": true,
      "priority": 293,
      "triggers": ["ack"],
      "actions": [
        {
          "entity_patterns": [
            {
              "name": "test-resource-remediation-instruction-execution-abort-5"
            }
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
    When I wait 1s
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 410
    When I do GET /api/v4/alarms?search=test-resource-remediation-instruction-execution-abort-5&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
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
      ]
    }
    """
