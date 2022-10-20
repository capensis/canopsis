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
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-1&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-2&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
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
    When I wait 100ms
    When the response array key "data.0.v.steps" should contain only one:
    """json
    {
      "_t": "autoinstructionstart"
    }
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
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-3&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
   When I wait the end of 4 events processing
   When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-4-1&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
   When I wait the end of 2 events processing
   When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-4-2&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-5&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-6-1&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-6-2&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
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
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-9&with_steps=true
    Then the response code should be 200
    Then the response array key "data.0.v.steps" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
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
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-10&with_steps=true
    Then the response code should be 200
    Then the response array key "data.0.v.steps" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
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
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-11&with_steps=true
    Then the response code should be 200
    Then the response array key "data.0.v.steps" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
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
    When I wait the end of event processing
    When I wait 500ms
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
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-12&opened=false&with_steps=true
    Then the response code should be 200
    Then the response array key "data.0.v.steps" should contain only:
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
    When I do GET /api/v4/alarms?search=test-resource-to-run-auto-instruction-12&opened=true&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              }
            ]
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
