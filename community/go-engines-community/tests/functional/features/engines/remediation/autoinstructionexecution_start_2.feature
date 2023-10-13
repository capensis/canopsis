Feature: run an auto instruction
  I need to be able to run an auto instruction

  @concurrent
  Scenario: given new alarm should run auto instructions, first instruction jobs should be stopped because of stopOnFail flag, but next instruction execution shouldn't be stopped
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-second-1",
      "connector_name": "test-connector-name-to-run-auto-instruction-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-second-1",
      "resource": "test-resource-to-run-auto-instruction-second-1",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-second-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-run-auto-instruction-second-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-1",
        "component": "test-component-to-run-auto-instruction-second-1",
        "resource": "test-resource-to-run-auto-instruction-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-1",
        "component": "test-component-to-run-auto-instruction-second-1",
        "resource": "test-resource-to-run-auto-instruction-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-1",
        "component": "test-component-to-run-auto-instruction-second-1",
        "resource": "test-resource-to-run-auto-instruction-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-1",
        "component": "test-component-to-run-auto-instruction-second-1",
        "resource": "test-resource-to-run-auto-instruction-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-1",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-1",
        "component": "test-component-to-run-auto-instruction-second-1",
        "resource": "test-resource-to-run-auto-instruction-second-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-second-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3
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
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-1-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-1-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-1-2-name."
      }
    ]
    """

  @concurrent
  Scenario: given new alarm should run auto instructions, auto instruction with failed jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-second-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-second-2",
      "resource": "test-resource-to-run-auto-instruction-second-2",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-second-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-run-auto-instruction-second-2",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-2",
        "component": "test-component-to-run-auto-instruction-second-2",
        "resource": "test-resource-to-run-auto-instruction-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-2",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-2",
        "component": "test-component-to-run-auto-instruction-second-2",
        "resource": "test-resource-to-run-auto-instruction-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-2",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-2",
        "component": "test-component-to-run-auto-instruction-second-2",
        "resource": "test-resource-to-run-auto-instruction-second-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-second-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3
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
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-2-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-2-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-2-name."
      }
    ]
    """

  @concurrent
  Scenario: given auto instructions should not add steps to resolved alarm and new alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-second-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-second-3",
      "resource": "test-resource-to-run-auto-instruction-second-3",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-second-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-run-auto-instruction-second-3",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-3",
        "component": "test-component-to-run-auto-instruction-second-3",
        "resource": "test-resource-to-run-auto-instruction-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-run-auto-instruction-second-3",
        "connector_name": "test-connector-name-to-run-auto-instruction-second-3",
        "component": "test-component-to-run-auto-instruction-second-3",
        "resource": "test-resource-to-run-auto-instruction-second-3",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-second-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-second-3",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-to-run-auto-instruction-second-3",
      "resource": "test-resource-to-run-auto-instruction-second-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-second-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-second-3",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-to-run-auto-instruction-second-3",
      "resource": "test-resource-to-run-auto-instruction-second-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-second-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-second-3",
      "resource": "test-resource-to-run-auto-instruction-second-3",
      "state": 2,
      "output": "test-output-to-run-auto-instruction-second-3"
    }
    """
    When I wait 5s
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-second-3&opened=false
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": false,
        "steps": {
          "page": 1
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
        "_t": "autoinstructionstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-3-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-3-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-3-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "initiator": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-second-3-1-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "cancel"
      },
      {
        "_t": "statusinc",
        "val": 4
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-second-3&opened=true
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": true,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """
