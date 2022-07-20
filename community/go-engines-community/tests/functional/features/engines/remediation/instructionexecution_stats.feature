Feature: update an instruction statistics
  I need to be able to update an instruction statistics

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
    When I send an event:
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
    When I send an event:
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
    When I wait the end of 2 events processing
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
    When I wait the end of event processing
    When I wait 2s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save response execution1Time={{ .lastResponse.completed_at }}
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
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
    When I wait the end of event processing
    When I wait 1s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save response execution2Time={{ .lastResponse.completed_at }}
    When I wait the end of event processing
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
        },
        {
          "executed_on": {{ .creationTime }},
          "alarm": null
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
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

  Scenario: given auto instruction execution should update statistics
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
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
    When I send an event:
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-stats-update-2-1
    Then the response code should be 200
    When I save response alarm1ID={{ (index .lastResponse.data 0)._id }}
    When I send an event:
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
    When I wait the end of event processing
    When I wait 5s
    When I send an event:
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
    When I wait the end of 3 events processing
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
        },
        {
          "executed_on": {{ .creationTime }},
          "alarm": null
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
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
