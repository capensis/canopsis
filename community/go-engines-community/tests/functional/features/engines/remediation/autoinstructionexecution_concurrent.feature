Feature: run an auto instruction
  I need to be able to run an auto instruction

  Scenario: given new alarm should run auto instructions
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-13",
      "connector_name": "test-connector-name-to-run-auto-instruction-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-13",
      "resource": "test-resource-to-run-auto-instruction-13-1",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-13-1"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-13",
      "connector_name": "test-connector-name-to-run-auto-instruction-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-13",
      "resource": "test-resource-to-run-auto-instruction-13-2",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-13-2"
    }
    """
    When I wait the end of 6 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-13&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        },
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
