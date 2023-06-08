Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  @concurrent
  Scenario: given alarm and failed manual and running manual instruction execution should return failed manual and running instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-1",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-second-1-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-1",
      "component": "test-component-to-alarm-instruction-get-icons-second-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionfailed",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-1",
      "component": "test-component-to-alarm-instruction-get-icons-second-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-1",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-second-1-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-1",
      "component": "test-component-to-alarm-instruction-get-icons-second-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 5,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-1-2-name"],
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-1-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed auto and running manual instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-2",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-2",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-2",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-2",
        "component": "test-component-to-alarm-instruction-get-icons-second-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-2",
        "component": "test-component-to-alarm-instruction-get-icons-second-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-2",
        "component": "test-component-to-alarm-instruction-get-icons-second-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-second-2-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-2",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-2",
      "component": "test-component-to-alarm-instruction-get-icons-second-2",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 6,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-2-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-2-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed auto and running auto instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-3",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-3",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-3",
        "component": "test-component-to-alarm-instruction-get-icons-second-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-3",
        "component": "test-component-to-alarm-instruction-get-icons-second-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-3",
        "component": "test-component-to-alarm-instruction-get-icons-second-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-3",
        "component": "test-component-to-alarm-instruction-get-icons-second-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 6,
          "running_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-3-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-3-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed manual and available manual instruction should return failed manual and available instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-4",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-4",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-second-4-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-4",
      "component": "test-component-to-alarm-instruction-get-icons-second-4",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-4",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionfailed",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-4",
      "component": "test-component-to-alarm-instruction-get-icons-second-4",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-4&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 7,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-4-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed auto and available manual instruction should return failed auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-5",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-5",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-5",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-5"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-5",
        "component": "test-component-to-alarm-instruction-get-icons-second-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-5",
        "component": "test-component-to-alarm-instruction-get-icons-second-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-5",
        "component": "test-component-to-alarm-instruction-get-icons-second-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-5",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-5&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 8,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-5-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and available manual instruction should return available manual instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-6",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-6",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-6",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-6&with_instructions=true
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

  @concurrent
  Scenario: given alarm and successful auto instruction execution should return successful auto instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-7",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-7",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-7",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-7"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-7",
        "component": "test-component-to-alarm-instruction-get-icons-second-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-7",
        "component": "test-component-to-alarm-instruction-get-icons-second-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-7",
        "component": "test-component-to-alarm-instruction-get-icons-second-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-7",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-7&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-7-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful manual instruction execution should return successful manual instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-8",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-8",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-8
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-second-8"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-8",
      "component": "test-component-to-alarm-instruction-get-icons-second-8",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-8",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-8",
      "component": "test-component-to-alarm-instruction-get-icons-second-8",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-8",
      "source_type": "resource"
    }
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-8&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 11,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-8-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful auto and running manual instruction execution should return successful auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-second-9",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-second-9",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-9",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-second-9"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-9",
        "component": "test-component-to-alarm-instruction-get-icons-second-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-9",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-9",
        "component": "test-component-to-alarm-instruction-get-icons-second-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-9",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-second-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-9",
        "component": "test-component-to-alarm-instruction-get-icons-second-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-second-9",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-9
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-second-9-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-icons-second-9",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-second-9",
      "component": "test-component-to-alarm-instruction-get-icons-second-9",
      "resource": "test-resource-to-alarm-instruction-get-icons-second-9",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-second-9&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 12,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-9-2-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-second-9-1-name"]
        }
      ]
    }
    """
