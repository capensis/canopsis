Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  @concurrent
  Scenario: given alarm and failed auto and running manual instruction execution should return failed auto and running instruction icon
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
      "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-icons-fifth-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
        "source_type": "resource"
      }
    ]
    """
    When I wait 2s
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fifth-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-icons-fifth-1-2"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-icons-fifth-1",
        "connector_name": "test-connector-name-to-alarm-instruction-get-icons-fifth-1",
        "component": "test-component-to-alarm-instruction-get-icons-fifth-1",
        "resource": "test-resource-to-alarm-instruction-get-icons-fifth-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-icons-fifth-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 4,
          "failed_manual_instructions": ["test-instruction-to-alarm-instruction-get-icons-fifth-1-2-name"],
          "failed_auto_instructions": ["test-instruction-to-alarm-instruction-get-icons-fifth-1-1-name"]
        }
      ]
    }
    """
