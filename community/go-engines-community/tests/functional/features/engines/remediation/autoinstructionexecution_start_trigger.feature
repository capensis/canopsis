Feature: run an auto instruction
  I need to be able to run an auto instruction

  @concurrent
  Scenario: given trigger should run auto instructions
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["stateinc"],
      "name": "test-instruction-to-run-auto-instruction-trigger-1-1-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-1-1-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-1"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1",
          "stop_on_fail": false
        },
        {
          "job": "test-job-to-run-auto-instruction-2"
        }
      ],
      "priority": 80
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["stateinc", "create"],
      "name": "test-instruction-to-run-auto-instruction-trigger-1-2-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-1-2-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-1"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ],
      "priority": 81
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["stateinc"],
      "name": "test-instruction-to-run-auto-instruction-trigger-1-3-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-1-3-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-1"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ],
      "priority": 82
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-trigger-1",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-trigger-1",
      "resource": "test-resource-to-run-auto-instruction-trigger-1",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-trigger-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-trigger-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc",
        "val": 1
      },
      {
        "_t": "statusinc",
        "val": 1
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name."
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-trigger-1",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-trigger-1",
      "resource": "test-resource-to-run-auto-instruction-trigger-1",
      "state": 2,
      "output": "test-output-to-run-auto-instruction-trigger-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-1",
        "resource": "test-resource-to-run-auto-instruction-trigger-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-trigger-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc",
        "val": 1
      },
      {
        "_t": "statusinc",
        "val": 1
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name."
      },
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-3-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-3-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-3-name."
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
      {
        "_t": "stateinc",
        "val": 1
      },
      {
        "_t": "statusinc",
        "val": 1
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-2-name."
      },
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-1-3-name."
      }
    ]
    """

  @concurrent
  Scenario: given trigger should run auto instructions after manual instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-run-auto-instruction-trigger-2-1-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-2-1-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-2"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-run-auto-instruction-trigger-2-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-run-auto-instruction-trigger-2-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-run-auto-instruction-trigger-2-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-run-auto-instruction-trigger-2-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response manualInstructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["stateinc"],
      "name": "test-instruction-to-run-auto-instruction-trigger-2-2-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-2-2-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-2"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ],
      "priority": 83
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-trigger-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-trigger-2",
      "resource": "test-resource-to-run-auto-instruction-trigger-2",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-trigger-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-trigger-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .manualInstructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "instructionstarted",
      "connector": "test-connector-to-run-auto-instruction-trigger-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-2",
      "source_type": "resource",
      "component": "test-component-to-run-auto-instruction-trigger-2",
      "resource": "test-resource-to-run-auto-instruction-trigger-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-trigger-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-trigger-2",
      "resource": "test-resource-to-run-auto-instruction-trigger-2",
      "state": 2,
      "output": "test-output-to-run-auto-instruction-trigger-2"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "instructioncompleted",
      "connector": "test-connector-to-run-auto-instruction-trigger-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-2",
      "source_type": "resource",
      "component": "test-component-to-run-auto-instruction-trigger-2",
      "resource": "test-resource-to-run-auto-instruction-trigger-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-2",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-2",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-2",
        "resource": "test-resource-to-run-auto-instruction-trigger-2"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-2",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-2",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-2",
        "resource": "test-resource-to-run-auto-instruction-trigger-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-trigger-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
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
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-1-name."
      },
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-2-name."
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
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
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-1-name."
      },
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-2-2-name."
      }
    ]
    """

  @concurrent
  Scenario: given trigger should run auto instructions after simplified manual instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-run-auto-instruction-trigger-3-1-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-3-1-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-3"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
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
    When I save response manualInstructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["stateinc"],
      "name": "test-instruction-to-run-auto-instruction-trigger-3-2-name",
      "description": "test-instruction-to-run-auto-instruction-trigger-3-2-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-run-auto-instruction-trigger-3"
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-1"
        }
      ],
      "priority": 84
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-trigger-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-trigger-3",
      "resource": "test-resource-to-run-auto-instruction-trigger-3",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-trigger-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-trigger-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .manualInstructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "instructionstarted",
      "connector": "test-connector-to-run-auto-instruction-trigger-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-3",
      "source_type": "resource",
      "component": "test-component-to-run-auto-instruction-trigger-3",
      "resource": "test-resource-to-run-auto-instruction-trigger-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-trigger-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-trigger-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-trigger-3",
      "resource": "test-resource-to-run-auto-instruction-trigger-3",
      "state": 2,
      "output": "test-output-to-run-auto-instruction-trigger-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-3",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-3",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-3",
        "resource": "test-resource-to-run-auto-instruction-trigger-3"
      },
      {
        "connector": "test-connector-to-run-auto-instruction-trigger-3",
        "connector_name": "test-connector-name-to-run-auto-instruction-trigger-3",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-to-run-auto-instruction-trigger-3",
        "resource": "test-resource-to-run-auto-instruction-trigger-3"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-trigger-3&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
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
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-1-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-1-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-2-name."
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
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
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-1-name."
      },
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-trigger-3-2-name."
      }
    ]
    """
