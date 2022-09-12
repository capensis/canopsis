Feature: update an instruction statistics
  I need to be able to update an instruction statistics
  
  Scenario: given auto instruction execution should return flags in alarm API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-alarm-get-auto-instruction-flags-1-1-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-alarm-get-auto-instruction-flags-1"
        }
      ],
      "description": "test-instruction-to-alarm-get-auto-instruction-flags-1-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-5"
        }
      ],
      "priority": 30
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-alarm-get-auto-instruction-flags-1-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-alarm-get-auto-instruction-flags-1"
        }
      ],
      "description": "test-instruction-to-alarm-get-auto-instruction-flags-1-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-5"
        }
      ],
      "priority": 31
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-get-auto-instruction-flags-1",
      "connector_name": "test-connector-name-to-alarm-get-auto-instruction-flags-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-get-auto-instruction-flags-1",
      "resource": "test-resource-to-alarm-get-auto-instruction-flags-1",
      "state": 1,
      "output": "test-output-to-alarm-get-auto-instruction-flags-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": false,
          "is_all_auto_instructions_completed": true,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """

  Scenario: given manual instruction execution should return flags in alarm API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-alarm-get-auto-instruction-flags-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-alarm-get-auto-instruction-flags-2"
        }
      ],
      "description": "test-instruction-to-alarm-get-auto-instruction-flags-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 3,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-alarm-get-auto-instruction-flags-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-alarm-get-auto-instruction-flags-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-alarm-get-auto-instruction-flags-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-alarm-get-auto-instruction-flags-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-get-auto-instruction-flags-2",
      "connector_name": "test-connector-name-to-alarm-get-auto-instruction-flags-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-get-auto-instruction-flags-2",
      "resource": "test-resource-to-alarm-get-auto-instruction-flags-2",
      "state": 1,
      "output": "test-output-to-alarm-get-auto-instruction-flags-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-2
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_manual_instruction_running": true,
          "is_manual_instruction_waiting_result": false
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_manual_instruction_running": false,
          "is_manual_instruction_waiting_result": true
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_manual_instruction_running": false,
          "is_manual_instruction_waiting_result": false
        }
      ]
    }
    """

  Scenario: given auto failed instruction execution should return flags in alarm API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-alarm-get-auto-instruction-flags-3-1-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-alarm-get-auto-instruction-flags-3"
        }
      ],
      "description": "test-instruction-to-alarm-get-auto-instruction-flags-3-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-6"
        }
      ],
      "priority": 32
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-alarm-get-auto-instruction-flags-3-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-alarm-get-auto-instruction-flags-3"
        }
      ],
      "description": "test-instruction-to-alarm-get-auto-instruction-flags-3-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-5"
        }
      ],
      "priority": 33
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-get-auto-instruction-flags-3",
      "connector_name": "test-connector-name-to-alarm-get-auto-instruction-flags-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-get-auto-instruction-flags-3",
      "resource": "test-resource-to-alarm-get-auto-instruction-flags-3",
      "state": 1,
      "output": "test-output-to-alarm-get-auto-instruction-flags-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-3&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": true
        }
      ]
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-auto-instruction-flags-3&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": false,
          "is_all_auto_instructions_completed": true,
          "is_auto_instruction_failed": true
        }
      ]
    }
    """
