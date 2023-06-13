Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  @concurrent
  Scenario: given alarm and failed auto and running simplified manual instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-1",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-1-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 6,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-1-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-1-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed simplified manual and available simplified manual instruction should return failed manual and available instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-2",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-2",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-2",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-2-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-2",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-2",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-2",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 7,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-2-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and failed auto and available simplified manual instruction should return failed auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-3",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-3",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-3",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-3",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-3",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-3",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-3",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 8,
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-3-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and available simplified manual instruction should return available manual instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-4",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-4",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-4&with_instructions=true
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
  Scenario: given alarm and successful simplified manual instruction execution should return successful manual instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-5",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-5",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-5",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-5"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-5",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-5",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-5",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-5",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-5",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-5&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 11,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-5-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful auto and running simplified manual instruction execution should return successful auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-6",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-6",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-6",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-6"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-6",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-6",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-6",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-6",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-6",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-6",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-6",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-6",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-6",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-6
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-6-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-6&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 12,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-6-2-name"],
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-6-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful simplified manual and running simplified manual instruction execution should return successful manual and running instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-7",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-7",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-7",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-7
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-7-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-7",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-7",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-7",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-7",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-7",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-7-2"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-7&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 13,
          "running_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-7-2-name"],
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-7-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful auto and available simplified manual instruction execution should return successful auto and available instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-8",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-8",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-8"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-8",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-8",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-8",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-8",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-8",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-8",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-8",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-8&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 14,
          "successful_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-8-1-name"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm and successful simplified manual and available simplified manual instruction execution should return successful manual and available instruction icon
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fourth-9",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fourth-9",
      "resource": "test-resource-to-alarm-instruction-get-icons-fourth-9",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fourth-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-9
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fourth-9-1"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-9",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-9",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fourth-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fourth-9",
        "component": "test-component-to-alarm-instruction-get-icons-fourth-9",
        "resource": "test-resource-to-alarm-instruction-get-icons-fourth-9",
        "source_type": "resource"
      }
    ]
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fourth-9&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 15,
          "successful_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fourth-9-1-name"]
        }
      ]
    }
    """
