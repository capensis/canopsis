Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  @concurrent
  Scenario: given alarm and successful manual and running manual instruction execution should return successful manual and running instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-1",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-1-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-third-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-1",
      "component": "test-component-to-alarm-instruction-get-icons-third-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-alarm-instruction-get-icons-third-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-1",
      "component": "test-component-to-alarm-instruction-get-icons-third-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-1",
      "source_type": "resource"
    }
    """
    When I wait 3s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-1-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-third-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-1",
      "component": "test-component-to-alarm-instruction-get-icons-third-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 13,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-1-2-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-1-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful auto and available manual instruction execution should return successful auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-2",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-2",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-2",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-2",
        "component": "test-component-to-alarm-instruction-get-icons-third-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-2",
        "component": "test-component-to-alarm-instruction-get-icons-third-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-2",
        "component": "test-component-to-alarm-instruction-get-icons-third-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-2",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 14,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-2-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful manual and available manual instruction execution should return successful manual and available instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-3",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-3",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-3-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-third-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-3",
      "component": "test-component-to-alarm-instruction-get-icons-third-3",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-3",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-alarm-instruction-get-icons-third-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-3",
      "component": "test-component-to-alarm-instruction-get-icons-third-3",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-3",
      "source_type": "resource"
    }
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 15,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-3-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and running simplified manual instruction execution should return manual instruction running icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-4",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-4",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-4"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-4&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 1,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-4-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and completed simplified manual and failed auto instruction execution should return failed auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-5",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-5",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-5",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-5"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-5",
        "component": "test-component-to-alarm-instruction-get-icons-third-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-5",
        "component": "test-component-to-alarm-instruction-get-icons-third-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-5",
        "component": "test-component-to-alarm-instruction-get-icons-third-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-5",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-5-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-5",
        "component": "test-component-to-alarm-instruction-get-icons-third-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-5",
        "component": "test-component-to-alarm-instruction-get-icons-third-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-5",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-5&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-5-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-5-2-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed simplified manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-6",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-6",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-6",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-6
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-6"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-6",
        "component": "test-component-to-alarm-instruction-get-icons-third-6",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-6",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-6",
        "component": "test-component-to-alarm-instruction-get-icons-third-6",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-6",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-6&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-6-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and completed automatic and failed simplified manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-7",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-7"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
        "component": "test-component-to-alarm-instruction-get-icons-third-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
        "component": "test-component-to-alarm-instruction-get-icons-third-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
        "component": "test-component-to-alarm-instruction-get-icons-third-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
        "component": "test-component-to-alarm-instruction-get-icons-third-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-7
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-7-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
        "component": "test-component-to-alarm-instruction-get-icons-third-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-7",
        "component": "test-component-to-alarm-instruction-get-icons-third-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-7",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-7&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-7-2-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-7-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and completed simplified manual and failed simplified manual instruction execution should return failed manual instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-8",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-8",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-8
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-8-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-8",
        "component": "test-component-to-alarm-instruction-get-icons-third-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-8",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-8",
        "component": "test-component-to-alarm-instruction-get-icons-third-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-8",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-8-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-8",
        "component": "test-component-to-alarm-instruction-get-icons-third-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-8",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-8",
        "component": "test-component-to-alarm-instruction-get-icons-third-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-8",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-8&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-8-1-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-8-2-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed manual and running manual instruction execution should return failed manual and running instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-third-9",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-third-9",
      "resource": "test-resource-to-alarm-instruction-get-icons-third-9",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-third-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-9
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-9-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-9",
        "component": "test-component-to-alarm-instruction-get-icons-third-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-9",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-third-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-third-9",
        "component": "test-component-to-alarm-instruction-get-icons-third-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-third-9",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-third-9-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-third-9&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 5,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-9-2-name"],
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-third-9-1-name"]
        }
      ]
    }
    """
