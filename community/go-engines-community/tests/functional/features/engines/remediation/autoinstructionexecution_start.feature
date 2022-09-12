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
    When I wait the end of 6 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
          "is_all_auto_instructions_completed": true
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-7-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_all_auto_instructions_completed": true
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-7-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_all_auto_instructions_completed": true
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
