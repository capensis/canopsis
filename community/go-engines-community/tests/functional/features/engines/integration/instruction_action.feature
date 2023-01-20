Feature: run an auto instruction
  I need to be able to run an auto instruction

  @concurrent
  Scenario: given changestate trigger should run auto instructions on changestate action
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["changestate"],
      "name": "test-instruction-instruction-action-1-name",
      "description": "test-instruction-instruction-action-1-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-instruction-action-1"
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
      "priority": 90
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-instruction-action-1-name",
      "priority": 10070,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-instruction-action-1"
                }
              }
            ]
          ],
          "parameters": {
            "state": 3
          },
          "type": "changestate",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-instruction-action-1",
      "connector_name": "test-connector-name-instruction-action-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-instruction-action-1",
      "resource": "test-resource-instruction-action-1",
      "state": 1,
      "output": "test-output-instruction-action-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "connector": "test-connector-instruction-action-1",
        "connector_name": "test-connector-name-instruction-action-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-instruction-action-1",
        "resource": "test-resource-instruction-action-1"
      },
      {
        "connector": "test-connector-instruction-action-1",
        "connector_name": "test-connector-name-instruction-action-1",
        "source_type": "resource",
        "event_type": "trigger",
        "component": "test-component-instruction-action-1",
        "resource": "test-resource-instruction-action-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-instruction-action-1&with_instructions=true until response code is 200 and body contains:
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
        "_t": "changestate",
        "a": "system",
        "val": 3
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-instruction-action-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-instruction-action-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-instruction-action-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-instruction-action-1-name."
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
        "_t": "changestate",
        "a": "system",
        "val": 3
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-instruction-action-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-instruction-action-1-name."
      }
    ]
    """
