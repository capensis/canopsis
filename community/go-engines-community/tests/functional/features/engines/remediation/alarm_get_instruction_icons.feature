Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  Scenario: given alarm and running manual instruction execution should return manual instruction running icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-1",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-1"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 1,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and running auto instruction execution should return auto instruction running icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-2",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-2",
      "resource": "test-resource-to-alarm-instruction-get-icons-2",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 2,
          "running_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-2-name"]
        }
      ]
    }
    """
    When I wait the end of event processing

  Scenario: given alarm and failed auto instruction execution should return failed auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-3",
      "resource": "test-resource-to-alarm-instruction-get-icons-3",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-3"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-3-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed auto and failed auto instruction execution should return failed auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-4",
      "resource": "test-resource-to-alarm-instruction-get-icons-4",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-4"
    }
    """
    When I wait the end of 5 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-4&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-4-1-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-4-2-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed manual and failed auto instruction execution should return failed auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-5",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-5",
      "resource": "test-resource-to-alarm-instruction-get-icons-5",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-5"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-5-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-5&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-5-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-5-2-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-6",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-6",
      "resource": "test-resource-to-alarm-instruction-get-icons-6",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-6"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-6
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-6"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-6&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-6-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed automatic and failed manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-7",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-7",
      "resource": "test-resource-to-alarm-instruction-get-icons-7",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-7"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-7
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-7-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-7&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-7-1-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-7-2-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed manual and failed manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-8",
      "resource": "test-resource-to-alarm-instruction-get-icons-8",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-8
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-8-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-8-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-8&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-8-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-8-2-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed automatic and failed auto and manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-9",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-9",
      "resource": "test-resource-to-alarm-instruction-get-icons-9",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-9"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-9
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-9-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-9-3"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-9&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-9-2-name"],
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-9-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-9-3-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed manual and running manual instruction execution should return failed manual and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-10",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-10",
      "resource": "test-resource-to-alarm-instruction-get-icons-10",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-10"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-10
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-10-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-10-2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-10&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 5,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-10-2-name"],
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-10-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed auto and running manual instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-11",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-11",
      "resource": "test-resource-to-alarm-instruction-get-icons-11",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-11"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-11
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-11-2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-11&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 6,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-11-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-11-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed auto and running auto instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-12",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-12",
      "resource": "test-resource-to-alarm-instruction-get-icons-12",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-12"
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-12&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 6,
          "running_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-12-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-12-1-name"]
        }
      ]
    }
    """
    When I wait the end of event processing

  Scenario: given alarm and failed manual and available manual instruction should return failed manual and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-13",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-13",
      "resource": "test-resource-to-alarm-instruction-get-icons-13",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-13"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-13
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-13-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-13&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 7,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-13-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed auto and available manual instruction should return failed auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-14",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-14",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-14",
      "resource": "test-resource-to-alarm-instruction-get-icons-14",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-14"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-14&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 8,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-14-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and available manual instruction should return available manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-15",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-15",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-15",
      "resource": "test-resource-to-alarm-instruction-get-icons-15",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-15"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-15&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 9
        }
      ]
    }
    """
