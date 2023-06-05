Feature: scenarios should be triggered by remediation triggers
  I need to be able to trigger scenarios by remediation triggers

  @concurrent
  Scenario: given scenario and failed instruction should be triggered by autoinstructionresultfail trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-second-1-name",
      "enabled": true,
      "triggers": ["autoinstructionresultfail"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-remediation-triggers-second-1/test-component-action-remediation-triggers-second-1"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-second-1-ack"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-action-remediation-triggers-second-1",
      "connector_name": "test-connector-name-action-remediation-triggers-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-action-remediation-triggers-second-1",
      "resource": "test-resource-action-remediation-triggers-second-1",
      "state": 1,
      "output": "test output"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-action-remediation-triggers-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-second-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "system",
              "m": "test-resource-action-remediation-triggers-second-1-ack"
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
        "m": "Instruction test-instruction-action-remediation-triggers-second-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-1-name. Job test-job-action-remediation-triggers-2-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-1-name. Job test-job-action-remediation-triggers-2-name."
      },
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-1-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-second-1-ack"
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
      {
        "_t": "autoinstructionfail",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-1-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-second-1-ack"
      }
    ]
    """

  @concurrent
  Scenario: given scenario and completed instruction and not ok alarm should be triggered by autoinstructionresultfail trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-second-2-name",
      "enabled": true,
      "triggers": ["autoinstructionresultfail"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-remediation-triggers-second-2/test-component-action-remediation-triggers-second-2"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-resource-action-remediation-triggers-second-2-ack"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-action-remediation-triggers-second-2",
      "connector_name": "test-connector-name-action-remediation-triggers-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-action-remediation-triggers-second-2",
      "resource": "test-resource-action-remediation-triggers-second-2",
      "state": 1,
      "output": "test output"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-action-remediation-triggers-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-second-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "system",
              "m": "test-resource-action-remediation-triggers-second-2-ack"
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
        "m": "Instruction test-instruction-action-remediation-triggers-second-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-2-name. Job test-job-action-remediation-triggers-3-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-2-name. Job test-job-action-remediation-triggers-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-2-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-second-2-ack"
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-2-name."
      },
      {
        "_t": "ack",
        "a": "system",
        "m": "test-resource-action-remediation-triggers-second-2-ack"
      }
    ]
    """
