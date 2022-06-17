Feature: scenarios should be triggered by remediation triggers
  I need to be able to trigger scenarios by remediation triggers

  Scenario: given scenario should be triggered by instructionfail trigger
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-action-remediation-triggers-1",
      "connector_name" : "test-connector-name-action-remediation-triggers-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-remediation-triggers-1",
      "resource" : "test-resource-action-remediation-triggers-1",
      "state" : 1,
      "output" : "test output"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-1-name",
      "enabled": true,
      "triggers": ["instructionfail"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-resource-action-remediation-triggers-1/test-component-action-remediation-triggers-1"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-1-ack"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-action-remediation-triggers-1"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step:
    """
    {
      "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-1&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-action-remediation-triggers-1-name."
              },
              {
                "_t": "instructionfail",
                "a": "root",
                "m": "Instruction test-instruction-action-remediation-triggers-1-name."
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-resource-action-remediation-triggers-1-ack"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given scenario should be triggered by autoinstructionfail trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-2-name",
      "enabled": true,
      "triggers": ["autoinstructionfail"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-resource-action-remediation-triggers-2/test-component-action-remediation-triggers-2"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-2-ack"
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
    """
    {
      "connector" : "test-connector-action-remediation-triggers-2",
      "connector_name" : "test-connector-name-action-remediation-triggers-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-remediation-triggers-2",
      "resource" : "test-resource-action-remediation-triggers-2",
      "state" : 1,
      "output" : "test output"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-2&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-2-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-2-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-2-ack"
      }
    ]
    """

  Scenario: given scenario should be triggered by instructionjobfail trigger with auto instruction
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-3-name",
      "enabled": true,
      "triggers": ["instructionjobfail"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-resource-action-remediation-triggers-3/test-component-action-remediation-triggers-3"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-3-ack"
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
    """
    {
      "connector" : "test-connector-action-remediation-triggers-3",
      "connector_name" : "test-connector-name-action-remediation-triggers-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-remediation-triggers-3",
      "resource" : "test-resource-action-remediation-triggers-3",
      "state" : 1,
      "output" : "test output"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-3&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-3-name. Job test-job-action-remediation-triggers-1-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-3-ack"
      }
    ]
    """

  Scenario: given scenario should be triggered by instructionjobfail trigger with manual instruction
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-action-remediation-triggers-4",
      "connector_name" : "test-connector-name-action-remediation-triggers-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-remediation-triggers-4",
      "resource" : "test-resource-action-remediation-triggers-4",
      "state" : 1,
      "output" : "test output"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-4-name",
      "enabled": true,
      "triggers": ["instructionjobfail"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-resource-action-remediation-triggers-4/test-component-action-remediation-triggers-4"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-4-ack"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-action-remediation-triggers-4"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "test-job-action-remediation-triggers-2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-4&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "m": "Instruction test-instruction-action-remediation-triggers-4-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "m": "Instruction test-instruction-action-remediation-triggers-4-name. Job test-job-action-remediation-triggers-2-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-4-ack"
      }
    ]
    """

  Scenario: given scenario should be triggered by instructioncomplete trigger
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-action-remediation-triggers-5",
      "connector_name" : "test-connector-name-action-remediation-triggers-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-remediation-triggers-5",
      "resource" : "test-resource-action-remediation-triggers-5",
      "state" : 1,
      "output" : "test output"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-5-name",
      "enabled": true,
      "triggers": ["instructioncomplete"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-resource-action-remediation-triggers-5/test-component-action-remediation-triggers-5"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-5-ack"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-action-remediation-triggers-5"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-5&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-action-remediation-triggers-5-name."
              },
              {
                "_t": "instructioncomplete",
                "a": "root",
                "m": "Instruction test-instruction-action-remediation-triggers-5-name."
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-resource-action-remediation-triggers-5-ack"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given scenario should be triggered by autoinstructioncomplete trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-6-name",
      "enabled": true,
      "triggers": ["autoinstructioncomplete"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-resource-action-remediation-triggers-6/test-component-action-remediation-triggers-6"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-6-ack"
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
    """
    {
      "connector" : "test-connector-action-remediation-triggers-6",
      "connector_name" : "test-connector-name-action-remediation-triggers-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-action-remediation-triggers-6",
      "resource" : "test-resource-action-remediation-triggers-6",
      "state" : 1,
      "output" : "test output"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-6&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-6-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-6-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-6-ack"
      }
    ]
    """
