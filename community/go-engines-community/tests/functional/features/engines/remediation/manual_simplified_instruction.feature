Feature: run an manual simplified instruction
  I need to be able to run an manual simplified instruction

  Scenario: given new alarm should run manual simplified instructions
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-1",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-1",
      "resource": "test-resource-to-run-manual-simplified-instruction-1",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-1"
    }
    """
    Then the response code should be 200
    When I wait the end of 3 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name."
      }
    ]
    """

  Scenario: given new alarm should run manual instruction, instruction jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-2",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-2",
      "resource": "test-resource-to-run-manual-simplified-instruction-2",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-2"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name."
      }
    ]
    """

  Scenario: given new alarm should run manual instruction, instruction jobs shouldn't stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-3",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-3",
      "resource": "test-resource-to-run-manual-simplified-instruction-3",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-3"
    }
    """
    Then the response code should be 200
    When I wait the end of 3 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name."
      }
    ]
    """

  Scenario: given new alarm should run manual instruction with failed jobs, instruction jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-4",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-4",
      "resource": "test-resource-to-run-manual-simplified-instruction-4",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-4"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "instructionfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name."
      }
    ]
    """
