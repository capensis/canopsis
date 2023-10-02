Feature: run a job
  I need to be able to run a job
  Only admin should be able to run a job

  @concurrent
  Scenario: given multiple job executions for one instruction execution should get only last one
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-second-1-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-second-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-second-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-second-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-second-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-second-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-second-1-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-second-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-to-job-execution-start-second-1",
      "connector": "test-connector-to-job-execution-start-second-1",
      "connector_name": "test-connector-name-to-job-execution-start-second-1",
      "component": "test-component-to-job-execution-start-second-1",
      "resource": "test-resource-to-job-execution-start-second-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-second-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-second-1",
      "source_type": "resource"
    }
    """
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-second-1-step-1-operation-1",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-1-name",
                  "status": 1,
                  "fail_reason": ""
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response job1CompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).completed_at }}
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-long-succeeded-output
    """
    When I wait 1s
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I save request:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "opened": true,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-1-name. Job test-job-to-job-execution-start-second-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-1-name. Job test-job-to-job-execution-start-second-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-1-name. Job test-job-to-job-execution-start-second-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-1-name. Job test-job-to-job-execution-start-second-1-name."
      }
    ]
    """
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-second-1-step-1-operation-1",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-1-name",
                  "status": 1,
                  "fail_reason": ""
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response job2StartedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).started_at }}
    When I save response job2LaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).launched_at }}
    When I save response job2CompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).completed_at }}
    Then "job2StartedAt" > "job1CompletedAt"
    Then "job2LaunchedAt" >= "job2StartedAt"
    Then "job2CompletedAt" >= "job2LaunchedAt"
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-long-succeeded-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given start previous operation should not return job
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-second-2-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-second-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-second-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-second-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-second-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-1-description",
              "jobs": []
            },
            {
              "name": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-2",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-2-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-second-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-to-job-execution-start-second-2",
      "connector": "test-connector-to-job-execution-start-second-2",
      "connector_name": "test-connector-name-to-job-execution-start-second-2",
      "component": "test-component-to-job-execution-start-second-2",
      "resource": "test-resource-to-job-execution-start-second-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-second-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-second-2",
      "source_type": "resource"
    }
    """
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 1).operation_id }}
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {},
            {
              "name": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-2",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-2-name",
                  "status": 1,
                  "fail_reason": ""
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response job1CompletedAt={{ (index (index (index .lastResponse.steps 0).operations 1).jobs 0).completed_at }}
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 1).jobs 0)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-long-succeeded-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/previous
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {},
            {
              "name": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-2",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-2-name",
                  "status": null,
                  "fail_reason": "",
                  "started_at": null,
                  "launched_at": null,
                  "completed_at": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I wait 1s
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {},
            {
              "name": "test-instruction-to-job-execution-start-second-2-step-1-operation-1-2",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-2-name",
                  "status": 1,
                  "fail_reason": ""
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response job2StartedAt={{ (index (index (index .lastResponse.steps 0).operations 1).jobs 0).started_at }}
    When I save response job2LaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 1).jobs 0).launched_at }}
    When I save response job2CompletedAt={{ (index (index (index .lastResponse.steps 0).operations 1).jobs 0).completed_at }}
    Then "job2StartedAt" > "job1CompletedAt"
    Then "job2LaunchedAt" >= "job2StartedAt"
    Then "job2CompletedAt" >= "job2LaunchedAt"
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 1).jobs 0)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-long-succeeded-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save request:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "opened": true,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-2-name. Job test-job-to-job-execution-start-second-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-2-name. Job test-job-to-job-execution-start-second-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-2-name. Job test-job-to-job-execution-start-second-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-2-name. Job test-job-to-job-execution-start-second-2-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root John Doe admin@canopsis.net",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-second-2-name."
      }
    ]
    """

  @concurrent
  Scenario: given vtom job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-second-3-name",
      "config": "test-job-config-to-run-manual-vtom-job-1",
      "job_id": "test-job-succeeded",
      "payload": "{\"parameters\": [\"{{ `{{ .Alarm.Value.Resource }}` }}\",\"{{ `{{ .Entity.ID }}` }}\"]}",
      "multiple_executions": false,
      "job_wait_interval": {
        "value": 8,
        "unit": "s"
      }
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-second-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-second-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-second-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-second-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-second-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-second-3-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-second-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-to-job-execution-start-second-3",
      "connector": "test-connector-to-job-execution-start-second-3",
      "connector_name": "test-connector-name-to-job-execution-start-second-3",
      "component": "test-component-to-job-execution-start-second-3",
      "resource": "test-resource-to-job-execution-start-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-second-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-second-3",
      "source_type": "resource"
    }
    """
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-3-name",
                  "status": 1,
                  "fail_reason": ""
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-succeeded-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given failed vtom job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-second-4-name",
      "config": "test-job-config-to-run-manual-vtom-job-2",
      "job_id": "test-job-failed",
      "payload": "{\"parameters\": [\"{{ `{{ .Alarm.Value.Resource }}` }}\",\"{{ `{{ .Entity.ID }}` }}\"]}",
      "multiple_executions": false,
      "job_wait_interval": {
        "value": 8,
        "unit": "s"
      }
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-second-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-second-4"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-second-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-second-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-second-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-second-4-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-second-4-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-to-job-execution-start-second-4",
      "connector": "test-connector-to-job-execution-start-second-4",
      "connector_name": "test-connector-name-to-job-execution-start-second-4",
      "component": "test-component-to-job-execution-start-second-4",
      "resource": "test-resource-to-job-execution-start-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-second-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-second-4",
      "source_type": "resource"
    }
    """
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-second-4-name",
                  "status": 2,
                  "fail_reason": "job failed, check {{ .dummyApiURL }} for more info"
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-failed-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
