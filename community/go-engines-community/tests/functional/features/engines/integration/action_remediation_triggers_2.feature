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

  @concurrent
  Scenario: given resolved ok alarm and scenario should be triggered by autoinstructionresultok trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-second-3-name",
      "enabled": true,
      "triggers": ["autoinstructionresultok"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-remediation-triggers-second-3/test-component-action-remediation-triggers-second-3"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"{{ `{{ .Alarm.Value.Output }}` }}\",\"enabled\":true,\"triggers\":[\"comment\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\": \"eq\", \"value\": \"{{ `{{ .Entity.ID }}` }}\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "ticket_id": "_id"
            }
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
      "event_type": "check",
      "state": 1,
      "output": "test-output-action-remediation-triggers-second-3-1",
      "connector": "test-connector-action-remediation-triggers-second-3",
      "connector_name": "test-connector-name-action-remediation-triggers-second-3",
      "component": "test-component-action-remediation-triggers-second-3",
      "resource": "test-resource-action-remediation-triggers-second-3",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-action-remediation-triggers-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-3",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-action-remediation-triggers-second-3-1",
      "connector": "test-connector-action-remediation-triggers-second-3",
      "connector_name": "test-connector-name-action-remediation-triggers-second-3",
      "component": "test-component-action-remediation-triggers-second-3",
      "resource": "test-resource-action-remediation-triggers-second-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_close",
      "connector": "test-connector-action-remediation-triggers-second-3",
      "connector_name": "test-connector-name-action-remediation-triggers-second-3",
      "component": "test-component-action-remediation-triggers-second-3",
      "resource": "test-resource-action-remediation-triggers-second-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-remediation-triggers-second-3-2",
      "connector": "test-connector-action-remediation-triggers-second-3",
      "connector_name": "test-connector-name-action-remediation-triggers-second-3",
      "component": "test-component-action-remediation-triggers-second-3",
      "resource": "test-resource-action-remediation-triggers-second-3",
      "source_type": "resource"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "trigger",
      "resource": "test-resource-action-remediation-triggers-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-second-3&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-action-remediation-triggers-second-3"
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
        "m": "Instruction test-instruction-action-remediation-triggers-second-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-3-name. Job test-job-action-remediation-triggers-3-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-3-name. Job test-job-action-remediation-triggers-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-3-name."
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
    When I do GET /api/v4/scenarios?search=test-output-action-remediation-triggers-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-output-action-remediation-triggers-second-3-1"
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
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-second-3&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-action-remediation-triggers-second-3"
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
    Then the response body should contain:
    """json
    [
      {
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
            ]
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given resolved canceled alarm and scenario should be triggered by autoinstructionresultok trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-remediation-triggers-second-4-name",
      "enabled": true,
      "triggers": ["autoinstructionresultok"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-remediation-triggers-second-4/test-component-action-remediation-triggers-second-4"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"{{ `{{ .Alarm.Value.Output }}` }}\",\"enabled\":true,\"triggers\":[\"comment\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\": \"eq\", \"value\": \"{{ `{{ .Entity.ID }}` }}\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "ticket_id": "_id"
            }
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
      "event_type": "check",
      "state": 1,
      "output": "test-output-action-remediation-triggers-second-4-1",
      "connector": "test-connector-action-remediation-triggers-second-4",
      "connector_name": "test-connector-name-action-remediation-triggers-second-4",
      "component": "test-component-action-remediation-triggers-second-4",
      "resource": "test-resource-action-remediation-triggers-second-4",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-action-remediation-triggers-second-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "resource": "test-resource-action-remediation-triggers-second-4",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-action-remediation-triggers-second-4",
      "connector_name": "test-connector-name-action-remediation-triggers-second-4",
      "component": "test-component-action-remediation-triggers-second-4",
      "resource": "test-resource-action-remediation-triggers-second-4",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "connector": "test-connector-action-remediation-triggers-second-4",
      "connector_name": "test-connector-name-action-remediation-triggers-second-4",
      "component": "test-component-action-remediation-triggers-second-4",
      "resource": "test-resource-action-remediation-triggers-second-4",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-remediation-triggers-second-4-2",
      "connector": "test-connector-action-remediation-triggers-second-4",
      "connector_name": "test-connector-name-action-remediation-triggers-second-4",
      "component": "test-component-action-remediation-triggers-second-4",
      "resource": "test-resource-action-remediation-triggers-second-4",
      "source_type": "resource"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "trigger",
      "resource": "test-resource-action-remediation-triggers-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-second-4&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-action-remediation-triggers-second-4"
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
        "m": "Instruction test-instruction-action-remediation-triggers-second-4-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-4-name. Job test-job-action-remediation-triggers-3-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-4-name. Job test-job-action-remediation-triggers-3-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "a": "system",
        "m": "Instruction test-instruction-action-remediation-triggers-second-4-name."
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
    When I do GET /api/v4/scenarios?search=test-output-action-remediation-triggers-second-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-output-action-remediation-triggers-second-4-1"
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
    When I do GET /api/v4/alarms?search=test-resource-action-remediation-triggers-second-4&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-action-remediation-triggers-second-4"
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
    Then the response body should contain:
    """json
    [
      {
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
            ]
          }
        }
      }
    ]
    """
