Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  @concurrent
  Scenario: given manual instruction execution should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-stats-update-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-1-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-1-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-stats-update-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-stats-update-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-stats-update-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-stats-update-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-stats-update-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I save response creationTime={{ .lastResponse.last_modified }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-1",
      "connector_name": "test-connector-name-to-stats-update-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-1",
      "state": 1,
      "output": "test-output-to-stats-update-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-1",
      "connector_name": "test-connector-name-to-stats-update-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-2",
      "state": 2,
      "output": "test-output-to-stats-update-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-1-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-1-2
    Then the response code should be 200
    When I save response alarm2ID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm1ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-1"
    }
    """
    When I wait 2s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save response execution1Time={{ .lastResponse.completed_at }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-1",
      "connector_name": "test-connector-name-to-stats-update-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-1",
      "state": 0,
      "output": "test-output-to-stats-update-1"
    }
    """
    When I wait 5s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm2ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-2"
    }
    """
    When I wait 1s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save response execution2Time={{ .lastResponse.completed_at }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "component": "test-component-to-stats-update-1",
      "resource": "test-resource-to-stats-update-1-2"
    }
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/summary until response code is 200 and body contains:
    """json
    {
      "_id": "{{ .instructionID }}",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 1,
          "to": 1
        },
        "minor": {
          "from": 1,
          "to": 0
        }
      },
      "ok_alarm_states": 1,
      "execution_count": 2,
      "last_executed_on": {{ .execution2Time }},
      "name": "test-instruction-to-stats-update-1-name",
      "type": 0
    }
    """
    Then the response key "avg_complete_time" should not be "0"
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/changes
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 0
            },
            "major": {
              "from": 1,
              "to": 1
            },
            "minor": {
              "from": 1,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "execution_count": 2,
          "modified_on": {{ .creationTime }}
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
    Then the response key "data.0.avg_complete_time" should not be "0"
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": {{ .execution2Time }},
          "alarm": {
            "_id": "{{ .alarm2ID }}",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "val": 2
                },
                {
                  "_t": "statusinc",
                  "val": 1
                },
                {
                  "_t": "instructionstart",
                  "m": "Instruction test-instruction-to-stats-update-1-name."
                },
                {
                  "_t": "instructioncomplete",
                  "m": "Instruction test-instruction-to-stats-update-1-name."
                }
              ]
            }
          }
        },
        {
          "executed_on": {{ .execution1Time }},
          "alarm": {
            "_id": "{{ .alarm1ID }}",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "val": 1
                },
                {
                  "_t": "statusinc",
                  "val": 1
                },
                {
                  "_t": "instructionstart",
                  "m": "Instruction test-instruction-to-stats-update-1-name."
                },
                {
                  "_t": "instructioncomplete",
                  "m": "Instruction test-instruction-to-stats-update-1-name."
                },
                {
                  "_t": "statedec",
                  "val": 0
                }
              ]
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-stats-update-1-name&with_month_executions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "month_executions": 2,
          "last_executed_on": {{ .execution2Time }}
        }
      ]
    }
    """

  @concurrent
  Scenario: given auto instruction execution should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["create"],
      "name": "test-instruction-to-stats-update-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-2-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-2-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-stats-update-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-5"
        }
      ],
      "priority": 20
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I save response creationTime={{ .lastResponse.last_modified }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-2",
      "connector_name": "test-connector-name-to-stats-update-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-2",
      "resource": "test-resource-to-stats-update-2-1",
      "state": 1,
      "output": "test-output-to-stats-update-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-2-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-2",
      "connector_name": "test-connector-name-to-stats-update-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-2",
      "resource": "test-resource-to-stats-update-2-1",
      "state": 0,
      "output": "test-output-to-stats-update-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-2",
        "resource": "test-resource-to-stats-update-2-1"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-2",
        "resource": "test-resource-to-stats-update-2-1"
      }
    ]
    """
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-2",
      "connector_name": "test-connector-name-to-stats-update-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-2",
      "resource": "test-resource-to-stats-update-2-2",
      "state": 2,
      "output": "test-output-to-stats-update-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-2",
        "resource": "test-resource-to-stats-update-2-2"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-2",
        "resource": "test-resource-to-stats-update-2-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-2-2
    Then the response code should be 200
    When I save response alarm2ID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/summary until response code is 200 and body contains:
    """json
    {
      "_id": "{{ .instructionID }}",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 1,
          "to": 1
        },
        "minor": {
          "from": 1,
          "to": 0
        }
      },
      "ok_alarm_states": 1,
      "execution_count": 2,
      "name": "test-instruction-to-stats-update-2-name",
      "type": 1
    }
    """
    Then the response key "avg_complete_time" should not be "0"
    Then the response key "last_executed_on" should not be "0"
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/changes
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 0
            },
            "major": {
              "from": 1,
              "to": 1
            },
            "minor": {
              "from": 1,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "execution_count": 2,
          "modified_on": {{ .creationTime }}
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
    Then the response key "data.0.avg_complete_time" should not be "0"
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "{{ .alarm2ID }}"
          }
        },
        {
          "alarm": {
            "_id": "{{ .alarm1ID }}"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    Then the response array key "data.0.alarm.v.steps" should contain:
    """json
    [
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "statusinc",
        "val": 1
      },
      {
        "_t": "autoinstructionstart",
        "m": "Instruction test-instruction-to-stats-update-2-name."
      },
      {
        "_t": "instructionjobstart",
        "m": "Instruction test-instruction-to-stats-update-2-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "instructionjobcomplete",
        "m": "Instruction test-instruction-to-stats-update-2-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "m": "Instruction test-instruction-to-stats-update-2-name."
      }
    ]
    """
    Then the response array key "data.1.alarm.v.steps" should contain:
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
        "m": "Instruction test-instruction-to-stats-update-2-name."
      },
      {
        "_t": "instructionjobstart",
        "m": "Instruction test-instruction-to-stats-update-2-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "statedec",
        "val": 0
      },
      {
        "_t": "statusdec",
        "val": 0
      },
      {
        "_t": "instructionjobcomplete",
        "m": "Instruction test-instruction-to-stats-update-2-name. Job test-job-to-run-auto-instruction-5-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "m": "Instruction test-instruction-to-stats-update-2-name."
      }
    ]
    """
    Then the response key "data.0.executed_on" should not be "0"
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-stats-update-2-name&with_month_executions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "month_executions": 2
        }
      ]
    }
    """
    Then the response key "data.0.last_executed_on" should not be "0"

  @concurrent
  Scenario: given failed execution of manual instruction should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-stats-update-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-3-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-3-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-stats-update-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-stats-update-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-stats-update-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-stats-update-3-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-stats-update-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I save response creationTime={{ .lastResponse.last_modified }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-3",
      "connector_name": "test-connector-name-to-stats-update-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-1",
      "state": 1,
      "output": "test-output-to-stats-update-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-3",
      "connector_name": "test-connector-name-to-stats-update-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-2",
      "state": 2,
      "output": "test-output-to-stats-update-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-3-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-3-2
    Then the response code should be 200
    When I save response alarm2ID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm1ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-1"
    }
    """
    When I wait 2s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save response execution1Time={{ .lastResponse.completed_at }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-3",
      "connector_name": "test-connector-name-to-stats-update-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-1",
      "state": 0,
      "output": "test-output-to-stats-update-3"
    }
    """
    When I wait 5s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm2ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-2"
    }
    """
    When I wait 1s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step:
    """json
    {
        "failed": true
    }
    """
    Then the response code should be 200
    When I save response execution2Time={{ .lastResponse.completed_at }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionfailed",
      "component": "test-component-to-stats-update-3",
      "resource": "test-resource-to-stats-update-3-2"
    }
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/summary until response code is 200 and body contains:
    """json
    {
      "_id": "{{ .instructionID }}",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 0,
          "to": 0
        },
        "minor": {
          "from": 1,
          "to": 0
        }
      },
      "ok_alarm_states": 1,
      "execution_count": 2,
      "successful": 1,
      "avg_complete_time": 2,
      "last_executed_on": {{ .execution1Time }},
      "name": "test-instruction-to-stats-update-3-name",
      "type": 0
    }
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/changes
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 0
            },
            "major": {
              "from": 0,
              "to": 0
            },
            "minor": {
              "from": 1,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "execution_count": 2,
          "successful": 1,
          "avg_complete_time": 2,
          "modified_on": {{ .creationTime }}
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
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": {{ .execution1Time }},
          "alarm": {
            "_id": "{{ .alarm1ID }}",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "val": 1
                },
                {
                  "_t": "statusinc",
                  "val": 1
                },
                {
                  "_t": "instructionstart",
                  "m": "Instruction test-instruction-to-stats-update-3-name."
                },
                {
                  "_t": "instructioncomplete",
                  "m": "Instruction test-instruction-to-stats-update-3-name."
                },
                {
                  "_t": "statedec",
                  "val": 0
                }
              ]
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
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-stats-update-3-name&with_month_executions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "month_executions": 2,
          "last_executed_on": {{ .execution1Time }}
        }
      ]
    }
    """

  @concurrent
  Scenario: given failed execution of auto instruction should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["create"],
      "name": "test-instruction-to-stats-update-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-4-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-4-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-stats-update-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-9"
        }
      ],
      "priority": 20
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I save response creationTime={{ .lastResponse.last_modified }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-4",
      "connector_name": "test-connector-name-to-stats-update-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-4",
      "resource": "test-resource-to-stats-update-4-1",
      "state": 1,
      "output": "test-output-to-stats-update-4"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-4",
        "resource": "test-resource-to-stats-update-4-1"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-4",
        "resource": "test-resource-to-stats-update-4-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-4-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-4",
      "connector_name": "test-connector-name-to-stats-update-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-4",
      "resource": "test-resource-to-stats-update-4-1",
      "state": 0,
      "output": "test-output-to-stats-update-4"
    }
    """
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-4",
      "connector_name": "test-connector-name-to-stats-update-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-4",
      "resource": "test-resource-to-stats-update-4-2",
      "state": 2,
      "output": "test-output-to-stats-update-4"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-4",
        "resource": "test-resource-to-stats-update-4-2"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-4",
        "resource": "test-resource-to-stats-update-4-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-4-2
    Then the response code should be 200
    When I save response alarm2ID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/summary until response code is 200 and body contains:
    """json
    {
      "_id": "{{ .instructionID }}",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 0,
          "to": 0
        },
        "minor": {
          "from": 0,
          "to": 0
        }
      },
      "ok_alarm_states": 0,
      "execution_count": 2,
      "successful": 0,
      "avg_complete_time": 0,
      "last_executed_on": null,
      "name": "test-instruction-to-stats-update-4-name",
      "type": 1
    }
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/changes
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 0
            },
            "major": {
              "from": 0,
              "to": 0
            },
            "minor": {
              "from": 0,
              "to": 0
            }
          },
          "ok_alarm_states": 0,
          "execution_count": 2,
          "successful": 0,
          "avg_complete_time": 0,
          "modified_on": {{ .creationTime }}
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
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/executions
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-stats-update-4-name&with_month_executions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "month_executions": 2
        }
      ]
    }
    """
    Then the response key "data.0.last_executed_on" should not exist

  @concurrent
  Scenario: given simplified manual instruction execution should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-stats-update-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-5-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-5-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-stats-update-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-4"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I save response creationTime={{ .lastResponse.last_modified }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-5",
      "connector_name": "test-connector-name-to-stats-update-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-5",
      "resource": "test-resource-to-stats-update-5-1",
      "state": 1,
      "output": "test-output-to-stats-update-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-5",
      "connector_name": "test-connector-name-to-stats-update-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-5",
      "resource": "test-resource-to-stats-update-5-2",
      "state": 2,
      "output": "test-output-to-stats-update-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-5-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-5-2
    Then the response code should be 200
    When I save response alarm2ID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm1ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "component": "test-component-to-stats-update-5",
        "resource": "test-resource-to-stats-update-5-1"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-5",
        "resource": "test-resource-to-stats-update-5-1"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-5",
      "connector_name": "test-connector-name-to-stats-update-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-5",
      "resource": "test-resource-to-stats-update-5-1",
      "state": 0,
      "output": "test-output-to-stats-update-5"
    }
    """
    When I wait 5s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm2ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "component": "test-component-to-stats-update-5",
        "resource": "test-resource-to-stats-update-5-2"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-5",
        "resource": "test-resource-to-stats-update-5-2"
      }
    ]
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/summary until response code is 200 and body contains:
    """json
    {
      "_id": "{{ .instructionID }}",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 1,
          "to": 1
        },
        "minor": {
          "from": 1,
          "to": 0
        }
      },
      "ok_alarm_states": 1,
      "execution_count": 2,
      "name": "test-instruction-to-stats-update-5-name",
      "type": 2
    }
    """
    Then the response key "avg_complete_time" should not be "0"
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/changes
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 0
            },
            "major": {
              "from": 1,
              "to": 1
            },
            "minor": {
              "from": 1,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "execution_count": 2,
          "modified_on": {{ .creationTime }}
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
    Then the response key "data.0.avg_complete_time" should not be "0"
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/executions
    Then the response code should be 200
    Then the response array key "data.0.alarm.v.steps" should contain only:
    """json
    [
      {
        "_t": "stateinc",
        "val": 2
      },
      {
        "_t": "statusinc",
        "val": 1
      },
      {
        "_t": "instructionstart",
        "m": "Instruction test-instruction-to-stats-update-5-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name. Job test-job-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name. Job test-job-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name."
      }
    ]
    """
    Then the response array key "data.1.alarm.v.steps" should contain only:
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
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name. Job test-job-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name. Job test-job-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-stats-update-5-name."
      },
      {
        "_t": "statedec",
        "val": 0
      }
    ]
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-stats-update-5-name&with_month_executions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "month_executions": 2
        }
      ]
    }
    """

  @concurrent
  Scenario: given failed execution of simplified manual instruction should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-stats-update-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-6-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-stats-update-6-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-stats-update-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-manual-simplified-instruction-5"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I save response creationTime={{ .lastResponse.last_modified }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-6",
      "connector_name": "test-connector-name-to-stats-update-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-6",
      "resource": "test-resource-to-stats-update-6-1",
      "state": 1,
      "output": "test-output-to-stats-update-6"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-stats-update-6",
      "connector_name": "test-connector-name-to-stats-update-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-stats-update-6",
      "resource": "test-resource-to-stats-update-6-2",
      "state": 2,
      "output": "test-output-to-stats-update-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-6-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-6-2
    Then the response code should be 200
    When I save response alarm2ID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm1ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "component": "test-component-to-stats-update-6",
        "resource": "test-resource-to-stats-update-6-1"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-6",
        "resource": "test-resource-to-stats-update-6-1"
      }
    ]
    """
    When I wait 5s
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarm2ID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "component": "test-component-to-stats-update-6",
        "resource": "test-resource-to-stats-update-6-2"
      },
      {
        "event_type": "trigger",
        "component": "test-component-to-stats-update-6",
        "resource": "test-resource-to-stats-update-6-2"
      }
    ]
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/summary until response code is 200 and body contains:
    """json
    {
      "_id": "{{ .instructionID }}",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 0,
          "to": 0
        },
        "minor": {
          "from": 0,
          "to": 0
        }
      },
      "ok_alarm_states": 0,
      "execution_count": 2,
      "successful": 0,
      "avg_complete_time": 0,
      "last_executed_on": null,
      "name": "test-instruction-to-stats-update-6-name",
      "type": 2
    }
    """
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/changes
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 0
            },
            "major": {
              "from": 0,
              "to": 0
            },
            "minor": {
              "from": 0,
              "to": 0
            }
          },
          "ok_alarm_states": 0,
          "execution_count": 2,
          "successful": 0,
          "avg_complete_time": 0,
          "modified_on": {{ .creationTime }}
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
    When I do GET /api/v4/cat/instruction-stats/{{ .instructionID }}/executions
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-stats-update-6-name&with_month_executions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "month_executions": 2
        }
      ]
    }
    """
    Then the response key "data.0.last_executed_on" should not exist
