Feature: send activation event on unsnooze
  I need to be able to trigger rule on alarm activation

  @standalone
  Scenario: given event for new alarm and autoremediation should activate alarm after autoremediation is finished
    Given I am admin
    When I set config parameter alarm.activatealarmafterautoremediation=true
    When I wait the next periodical process

    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-1",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-1",
      "resource": "test-resource-to-test-auto-instruction-activate-event-1",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-1"
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 6,9
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
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-1-name. Job test-job-to-test-auto-instruction-activate-event-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-1-name. Job test-job-to-test-auto-instruction-activate-event-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-1-name."
      }
    ]
    """
    When I save response instructionCompleteTime={{ (index (index .lastResponse 0).data.steps.data 5).t }}
    Then the difference between alarmActivationDate instructionCompleteTime is in range 3,5

    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-2",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-2",
      "resource": "test-resource-to-test-auto-instruction-activate-event-2",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-2"
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 2,5
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
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-2-name. Job test-job-to-test-auto-instruction-activate-event-2-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-2-name. Job test-job-to-test-auto-instruction-activate-event-2-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-2-name."
      }
    ]
    """
    When I save response instructionCompleteTime={{ (index (index .lastResponse 0).data.steps.data 5).t }}
    Then the difference between alarmActivationDate instructionCompleteTime is in range -1,1

    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-3",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-3",
      "resource": "test-resource-to-test-auto-instruction-activate-event-3",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-3"
    }
    """
    When I wait the end of 6 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-3&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 3,6
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
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-1-name. Job test-job-to-test-auto-instruction-activate-event-3-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-1-name. Job test-job-to-test-auto-instruction-activate-event-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-2-name. Job test-job-to-test-auto-instruction-activate-event-3-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-2-name. Job test-job-to-test-auto-instruction-activate-event-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-3-2-name."
      }
    ]
    """
    When I save response instructionCompleteTime={{ (index (index .lastResponse 0).data.steps.data 9).t }}
    Then the difference between alarmActivationDate instructionCompleteTime is in range 1,3

    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-4",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-4",
      "resource": "test-resource-to-test-auto-instruction-activate-event-4",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-4"
    }
    """
    When I wait the end of 6 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-4&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 5,8
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
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-1-name. Job test-job-to-test-auto-instruction-activate-event-4-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-1-name. Job test-job-to-test-auto-instruction-activate-event-4-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-2-name. Job test-job-to-test-auto-instruction-activate-event-4-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-2-name. Job test-job-to-test-auto-instruction-activate-event-4-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-4-2-name."
      }
    ]
    """
    When I save response instructionCompleteTime={{ (index (index .lastResponse 0).data.steps.data 9).t }}
    Then the difference between alarmActivationDate instructionCompleteTime is in range -1,1

    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-test-auto-instruction-activate-event-5-name",
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-to-test-auto-instruction-activate-event-5"
                }
              }
            ]
          ],
          "type":"snooze",
          "parameters":{
            "duration": {
              "value": 2,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-5",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-5",
      "resource": "test-resource-to-test-auto-instruction-activate-event-5",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-5"
    }
    """
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-5&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 6,9
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
        "_t": "snooze"
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-5-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-5-name. Job test-job-to-test-auto-instruction-activate-event-5-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-5-name. Job test-job-to-test-auto-instruction-activate-event-5-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-5-name."
      }
    ]
    """

    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-test-auto-instruction-activate-event-6",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "8s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-test-auto-instruction-activate-event-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-6",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-6",
      "resource": "test-resource-to-test-auto-instruction-activate-event-6",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-6"
    }
    """
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-6&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 6,9
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
        "_t": "pbhenter"
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-6-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-6-name. Job test-job-to-test-auto-instruction-activate-event-6-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-6-name. Job test-job-to-test-auto-instruction-activate-event-6-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-6-name."
      },
      {
        "_t": "pbhleave"
      }
    ]
    """

    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-7",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-7",
      "resource": "test-resource-to-test-auto-instruction-activate-event-7",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-7"
    }
    """
    When I wait the end of 8 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-7&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 10,13
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
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-1-name. Job test-job-to-test-auto-instruction-activate-event-7-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-1-name. Job test-job-to-test-auto-instruction-activate-event-7-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-1-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-2-name. Job test-job-to-test-auto-instruction-activate-event-7-2-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-2-name. Job test-job-to-test-auto-instruction-activate-event-7-2-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-2-name."
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-3-name. Job test-job-to-test-auto-instruction-activate-event-7-3-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-3-name. Job test-job-to-test-auto-instruction-activate-event-7-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-7-3-name."
      }
    ]
    """
    When I save response instructionCompleteTime={{ (index (index .lastResponse 0).data.steps.data 13).t }}
    Then the difference between alarmActivationDate instructionCompleteTime is in range 1,3

    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-test-auto-instruction-activate-event-8-name",
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-to-test-auto-instruction-activate-event-8"
                }
              }
            ]
          ],
          "type":"snooze",
          "parameters":{
            "duration": {
              "value": 6,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-8",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-8",
      "resource": "test-resource-to-test-auto-instruction-activate-event-8",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-8"
    }
    """
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-8&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 3
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 5,8
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
        "_t": "snooze"
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-8-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-8-name. Job test-job-to-test-auto-instruction-activate-event-8-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-8-name. Job test-job-to-test-auto-instruction-activate-event-8-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-8-name."
      }
    ]
    """

    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-test-auto-instruction-activate-event-9",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "4s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-test-auto-instruction-activate-event-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-to-test-auto-instruction-activate-event-9",
      "connector_name": "test-connector-name-to-test-auto-instruction-activate-event-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-test-auto-instruction-activate-event-9",
      "resource": "test-resource-to-test-auto-instruction-activate-event-9",
      "state": 1,
      "output": "test-output-to-test-auto-instruction-activate-event-9"
    }
    """
    When I wait the end of 5 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-test-auto-instruction-activate-event-9&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 10
        }
      ]
    }
    """
    When I save response createTimestamp={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range 6,9
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
        "_t": "pbhenter"
      },
      {
        "_t": "autoinstructionstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-9-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-9-name. Job test-job-to-test-auto-instruction-activate-event-9-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-9-name. Job test-job-to-test-auto-instruction-activate-event-9-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-to-test-auto-instruction-activate-event-9-name."
      },
      {
        "_t": "pbhleave"
      }
    ]
    """    

    When I set config parameter alarm.activatealarmafterautoremediation=false
    When I wait the next periodical process
