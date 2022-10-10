Feature: update an instruction statistics
  I need to be able to update an instruction statistics

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
    When I send an event:
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
    When I send an event:
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
    When I send an event:
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
    When I wait the end of 2 events processing
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
    When I send an event:
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
    When I send an event:
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
      "data": [
        {
          "executed_on": {{ .creationTime }},
          "duration": 0,
          "alarm": null
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
