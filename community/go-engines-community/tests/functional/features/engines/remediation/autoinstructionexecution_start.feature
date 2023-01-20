Feature: run an auto instruction
  I need to be able to run an auto instruction

  Scenario: given new alarm should run auto instructions
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-1",
      "connector_name": "test-connector-name-to-run-auto-instruction-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-1",
      "resource": "test-resource-to-run-auto-instruction-1",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-1"
    }
    """
    When I wait the end of 8 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
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
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-2-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-3-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-3-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-3-name."
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
      {
        "_t": "stateinc",
        "val": 1
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-2-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-1-3-name."
      }
    ]
    """

  Scenario: given new alarm should not run next auto instructions if alarm is ok
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-2",
      "resource": "test-resource-to-run-auto-instruction-2",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-2"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-2",
      "connector_name": "test-connector-name-to-run-auto-instruction-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-2",
      "resource": "test-resource-to-run-auto-instruction-2",
      "state": 0,
      "output": "test-output-to-run-auto-instruction-2"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
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
        "m": "Instruction test-instruction-to-run-auto-instruction-2-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-2-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-2-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-2-1-name."
      },
      {
        "_t": "statedec",
        "val": 0
      },
      {
        "_t": "statusdec",
        "val": 0
      }
    ]
    """

  Scenario: given new alarm should run next auto instructions if instruction failed
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-3",
      "connector_name": "test-connector-name-to-run-auto-instruction-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-3",
      "resource": "test-resource-to-run-auto-instruction-3",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-3"
    }
    """
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-3&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-3-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-1-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-1-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-3-2-name."
      }
    ]
    """

  Scenario: given multiple alarms for the same instruction should run the instruction for all
    When I am admin
    When I send an event:
   """json
   [
     {
       "connector": "test-connector-to-run-auto-instruction-4",
       "connector_name": "test-connector-name-to-run-auto-instruction-4",
       "source_type": "resource",
       "event_type": "check",
       "component": "test-component-to-run-auto-instruction-4",
       "resource": "test-resource-to-run-auto-instruction-4-1",
       "state": 1,
       "output": "test-output-to-run-auto-instruction-4"
     },
     {
       "connector": "test-connector-to-run-auto-instruction-4",
       "connector_name": "test-connector-name-to-run-auto-instruction-4",
       "source_type": "resource",
       "event_type": "check",
       "component": "test-component-to-run-auto-instruction-4",
       "resource": "test-resource-to-run-auto-instruction-4-2",
       "state": 1,
       "output": "test-output-to-run-auto-instruction-4"
     }
   ]
   """
    When I wait the end of 6 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-4-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
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
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name."
      },
      {
       "_t": "instructionjobstart",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
       "_t": "instructionjobcomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
       "_t": "autoinstructioncomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name."
      },
      {
       "_t": "autoinstructionstart",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-2-name."
      },
      {
       "_t": "instructionjobstart",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-2-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
       "_t": "instructionjobcomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-2-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
       "_t": "autoinstructioncomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-2-name."
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-4-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
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
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name."
      },
      {
       "_t": "instructionjobcomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
       "_t": "autoinstructioncomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-1-name."
      },
      {
       "_t": "autoinstructionstart",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-3-name."
      },
      {
       "_t": "instructionjobstart",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-3-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
       "_t": "instructionjobcomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-3-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
       "_t": "autoinstructioncomplete",
       "a": "system",
       "m": "Instruction test-instruction-to-run-auto-instruction-4-3-name."
      }
    ]
    """

  Scenario: given http error during job execution should return failed job status
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-5",
      "connector_name": "test-connector-name-to-run-auto-instruction-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-5",
      "resource": "test-resource-to-run-auto-instruction-5",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-5"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-5&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-5-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-5-name. Job test-job-to-run-auto-instruction-4-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-5-name. Job test-job-to-run-auto-instruction-4-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-5-name."
      }
    ]
    """

  Scenario: given multiple alarms for the same instruction should run next auto instructions for both alarms if instruction failed
    When I am admin
    When I send an event:
   """json
   [
     {
       "connector": "test-connector-to-run-auto-instruction-6",
       "connector_name": "test-connector-name-to-run-auto-instruction-6",
       "source_type": "resource",
       "event_type": "check",
       "component": "test-component-to-run-auto-instruction-6",
       "resource": "test-resource-to-run-auto-instruction-6-1",
       "state": 1,
       "output": "test-output-to-run-auto-instruction-6"
     },
     {
       "connector": "test-connector-to-run-auto-instruction-6",
       "connector_name": "test-connector-name-to-run-auto-instruction-6",
       "source_type": "resource",
       "event_type": "check",
       "component": "test-component-to-run-auto-instruction-6",
       "resource": "test-resource-to-run-auto-instruction-6-2",
       "state": 1,
       "output": "test-output-to-run-auto-instruction-6"
     }
   ]
   """
    When I wait the end of 6 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-6-1&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name. Job test-job-to-run-auto-instruction-6-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name. Job test-job-to-run-auto-instruction-6-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-2-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-2-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-2-name."
      }
     ]
     """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-6-2&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name. Job test-job-to-run-auto-instruction-6-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name. Job test-job-to-run-auto-instruction-6-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-3-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-3-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-6-3-name."
      }
    ]
    """

  Scenario: given new alarm should run auto instructions with old patterns
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-7",
      "connector_name": "test-connector-name-to-run-auto-instruction-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-7",
      "resource": "test-resource-to-run-auto-instruction-7-1",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-7"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-7-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
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
        "m": "Instruction test-instruction-to-run-auto-instruction-7-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-7-1-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-7-1-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-7-1-name."
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-7",
      "connector_name": "test-connector-name-to-run-auto-instruction-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-7",
      "resource": "test-resource-to-run-auto-instruction-7-2",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-7"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-7-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
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
        "m": "Instruction test-instruction-to-run-auto-instruction-7-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-7-2-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-7-2-name. Job test-job-to-run-auto-instruction-2-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-7-2-name."
      }
    ]
    """

  Scenario: given new alarm should not run auto instructions with empty patterns
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-8",
      "connector_name": "test-connector-name-to-run-auto-instruction-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-8",
      "resource": "test-resource-to-run-auto-instruction-8",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-8"
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-8
    Then the response code should be 200
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
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {"_t": "stateinc"},
              {"_t": "statusinc"}
            ]
          }
        }
      }
    ]
    """

  Scenario: given new alarm should run auto instructions, auto instruction jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-9",
      "connector_name": "test-connector-name-to-run-auto-instruction-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-9",
      "resource": "test-resource-to-run-auto-instruction-9",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-9"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-9&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-9-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-9-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-9-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-9-name."
      }
    ]
    """

  Scenario: given new alarm should run auto instructions, first instruction jobs should be stopped because of stopOnFail flag, but next instruction execution shouldn't be stopped
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-10",
      "connector_name": "test-connector-name-to-run-auto-instruction-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-10",
      "resource": "test-resource-to-run-auto-instruction-10",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-10"
    }
    """
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-10&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-10-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-1-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-1-name. Job test-job-to-run-auto-instruction-3-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-2-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-10-2-name."
      }
    ]
    """

  Scenario: given new alarm should run auto instructions, auto instruction with failed jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-11",
      "connector_name": "test-connector-name-to-run-auto-instruction-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-11",
      "resource": "test-resource-to-run-auto-instruction-11",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-11"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-11&with_instructions=true until response code is 200 and body contains:
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
        "m": "Instruction test-instruction-to-run-auto-instruction-11-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-11-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-11-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-11-name."
      }
    ]
    """

  Scenario: given auto instructions should not add steps to resolved alarm and new alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-12",
      "connector_name": "test-connector-name-to-run-auto-instruction-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-12",
      "resource": "test-resource-to-run-auto-instruction-12",
      "state": 1,
      "output": "test-output-to-run-auto-instruction-12"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-12",
      "connector_name": "test-connector-name-to-run-auto-instruction-12",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-to-run-auto-instruction-12",
      "resource": "test-resource-to-run-auto-instruction-12"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-12",
      "connector_name": "test-connector-name-to-run-auto-instruction-12",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-to-run-auto-instruction-12",
      "resource": "test-resource-to-run-auto-instruction-12"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-auto-instruction-12",
      "connector_name": "test-connector-name-to-run-auto-instruction-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-auto-instruction-12",
      "resource": "test-resource-to-run-auto-instruction-12",
      "state": 2,
      "output": "test-output-to-run-auto-instruction-12"
    }
    """
    When I wait the end of event processing
    When I wait 5s
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-12&opened=false
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
        "m": "Instruction test-instruction-to-run-auto-instruction-12-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-12-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-12-1-name. Job test-job-to-run-auto-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-run-auto-instruction-12-1-name. Job test-job-to-run-auto-instruction-5-name."
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
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-12&opened=true
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
