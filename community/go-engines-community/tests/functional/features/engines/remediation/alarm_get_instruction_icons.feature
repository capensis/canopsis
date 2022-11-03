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

  Scenario: given alarm and successful auto instruction execution should return successful auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-16",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-16",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-16",
      "resource": "test-resource-to-alarm-instruction-get-icons-16",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-16"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-16&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-16-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful manual instruction execution should return successful manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-17",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-17",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-17",
      "resource": "test-resource-to-alarm-instruction-get-icons-17",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-17"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-17
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-17"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-17&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 11,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-17-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful auto and running manual instruction execution should return successful auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-18",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-18",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-18",
      "resource": "test-resource-to-alarm-instruction-get-icons-18",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-18"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-18
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-18-2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-18&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 12,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-18-2-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-18-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful manual and running manual instruction execution should return successful manual and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-19",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-19",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-19",
      "resource": "test-resource-to-alarm-instruction-get-icons-19",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-19"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-19
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-19-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-19-2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-19&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 13,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-19-2-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-19-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful auto and available manual instruction execution should return successful auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-20",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-20",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-20",
      "resource": "test-resource-to-alarm-instruction-get-icons-20",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-20"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-20&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 14,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-20-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful manual and available manual instruction execution should return successful manual and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-21",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-21",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-21",
      "resource": "test-resource-to-alarm-instruction-get-icons-21",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-21"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-21
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-21-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-21&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 15,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-21-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and running simplified manual instruction execution should return manual instruction running icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-22",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-22",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-22",
      "resource": "test-resource-to-alarm-instruction-get-icons-22",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-22"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-22
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-22"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-22&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 1,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-22-name"]
        }
      ]
    }
    """
    When I wait the end of 2 events processing

  Scenario: given alarm and completed simplified manual and failed auto instruction execution should return failed auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-23",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-23",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-23",
      "resource": "test-resource-to-alarm-instruction-get-icons-23",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-23"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-23
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-23-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-23&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-23-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-23-2-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed simplified manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-24",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-24",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-24",
      "resource": "test-resource-to-alarm-instruction-get-icons-24",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-24"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-24
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-24"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-24&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-24-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed automatic and failed simplified manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-25",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-25",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-25",
      "resource": "test-resource-to-alarm-instruction-get-icons-25",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-25"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-25
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-25-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-25&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-25-2-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-25-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and completed simplified manual and failed simplified manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-26",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-26",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-26",
      "resource": "test-resource-to-alarm-instruction-get-icons-26",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-26"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-26
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-26-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-26-2"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-26&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-26-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-26-2-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed manual and running manual instruction execution should return failed manual and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-27",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-27",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-27",
      "resource": "test-resource-to-alarm-instruction-get-icons-27",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-27"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-27
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-27-1"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-27-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-27&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 5,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-27-2-name"],
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-27-1-name"]
        }
      ]
    }
    """
    When I wait the end of 2 events processing

  Scenario: given alarm and failed auto and running simplified manual instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-28",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-28",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-28",
      "resource": "test-resource-to-alarm-instruction-get-icons-28",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-28"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-28
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-28-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-28&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 6,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-28-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-28-1-name"]
        }
      ]
    }
    """
    When I wait the end of 2 events processing

  Scenario: given alarm and failed simplified manual and available simplified manual instruction should return failed manual and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-29",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-29",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-29",
      "resource": "test-resource-to-alarm-instruction-get-icons-29",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-29"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-29
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-29-1"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-29&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 7,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-29-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed auto and available simplified manual instruction should return failed auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-30",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-30",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-30",
      "resource": "test-resource-to-alarm-instruction-get-icons-30",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-30"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-30&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 8,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-30-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and available simplified manual instruction should return available manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-31",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-31",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-31",
      "resource": "test-resource-to-alarm-instruction-get-icons-31",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-31"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-31&with_instructions=true
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

  Scenario: given alarm and successful simplified manual instruction execution should return successful manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-32",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-32",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-32",
      "resource": "test-resource-to-alarm-instruction-get-icons-32",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-32"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-32
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-32"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-32&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 11,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-32-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful auto and running simplified manual instruction execution should return successful auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-33",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-33",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-33",
      "resource": "test-resource-to-alarm-instruction-get-icons-33",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-33"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-33
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-33-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-33&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 12,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-33-2-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-33-1-name"]
        }
      ]
    }
    """
    When I wait the end of 2 events processing

  Scenario: given alarm and successful simplified manual and running simplified manual instruction execution should return successful manual and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-34",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-34",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-34",
      "resource": "test-resource-to-alarm-instruction-get-icons-34",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-34"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-34
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-34-1"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I wait 3s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-34-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-34&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 13,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-34-2-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-34-1-name"]
        }
      ]
    }
    """
    When I wait the end of 2 events processing

  Scenario: given alarm and successful auto and available simplified manual instruction execution should return successful auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-35",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-35",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-35",
      "resource": "test-resource-to-alarm-instruction-get-icons-35",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-35"
    }
    """
    When I wait the end of 3 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-35&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 14,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-35-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and successful simplified manual and available simplified manual instruction execution should return successful manual and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-36",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-36",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-36",
      "resource": "test-resource-to-alarm-instruction-get-icons-36",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-36"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-36
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-36-1"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-36&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 15,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-36-1-name"]
        }
      ]
    }
    """

  Scenario: given alarm and failed auto and running manual instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-37",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-37",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-37",
      "resource": "test-resource-to-alarm-instruction-get-icons-37",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-37"
    }
    """
    When I wait the end of 3 events processing
    When I wait 2s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-37
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-37-2"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-37&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-37-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-37-1-name"]
        }
      ]
    }
    """
