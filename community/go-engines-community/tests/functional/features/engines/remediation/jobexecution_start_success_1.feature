Feature: run a job
  I need to be able to run a job
  Only admin should be able to run a job

  @concurrent
  Scenario: given job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-1-1-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response job1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-1-2-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource2\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity2\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response job2ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-1-step-1-operation-1-description",
              "jobs": ["{{ .job1ID }}", "{{ .job2ID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-1-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-1",
      "connector": "test-connector-to-job-execution-start-1",
      "connector_name": "test-connector-name-to-job-execution-start-1",
      "component": "test-component-to-job-execution-start-1",
      "resource": "test-resource-to-job-execution-start-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-1
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
      "resource": "test-resource-to-job-execution-start-1",
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
      "job": "{{ .job2ID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-job-execution-start-1-2-name",
      "job_id": "{{ .job2ID }}",
      "status": 0,
      "fail_reason": "",
      "started_at": null,
      "launched_at": null,
      "completed_at": null,
      "queue_number": 0
    }
    """
    When I do GET /api/v4/cat/job-executions/{{ .lastResponse._id }}/output
    Then the response code should be 404
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-1-step-1-operation-1",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "_id": "",
                  "job_id": "{{ .job1ID }}",
                  "name": "test-job-to-job-execution-start-1-1-name",
                  "status": null,
                  "fail_reason": "",
                  "started_at": null,
                  "launched_at": null,
                  "completed_at": null,
                  "queue_number": null
                },
                {
                  "job_id": "{{ .job2ID }}",
                  "name": "test-job-to-job-execution-start-1-2-name",
                  "status": 1,
                  "fail_reason": "",
                  "queue_number": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response job2StartedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 1).started_at }}
    When I save response job2LaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 1).launched_at }}
    When I save response job2CompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 1).completed_at }}
    When I save response expectedStartedAt=1
    Then "job2StartedAt" >= "expectedStartedAt"
    Then "job2LaunchedAt" >= "job2StartedAt"
    Then "job2CompletedAt" >= "job2LaunchedAt"
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 1)._id }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-long-succeeded-output
    """
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .job1ID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-job-execution-start-1-1-name",
      "job_id": "{{ .job1ID }}",
      "status": 0,
      "fail_reason": "",
      "started_at": null,
      "launched_at": null,
      "completed_at": null,
      "queue_number": 0
    }
    """
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-1-step-1-operation-1",
              "completed_at": null,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "job_id": "{{ .job1ID }}",
                  "name": "test-job-to-job-execution-start-1-1-name",
                  "status": 1,
                  "fail_reason": "",
                  "queue_number": null
                },
                {
                  "job_id": "{{ .job2ID }}",
                  "name": "test-job-to-job-execution-start-1-2-name",
                  "status": 1,
                  "fail_reason": "",
                  "queue_number": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response job1StartedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).started_at }}
    When I save response job1LaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).launched_at }}
    When I save response job1CompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).completed_at }}
    When I save response expectedStartedAt=1
    Then "job1StartedAt" >= "expectedStartedAt"
    Then "job1LaunchedAt" >= "job1StartedAt"
    Then "job1CompletedAt" >= "job1LaunchedAt"
    When I do GET /api/v4/cat/job-executions/{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0)._id }}/output
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
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-1-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-1-name."
      }
    ]
    """

  @concurrent
  Scenario: given error during job execution should return failed job status
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-2-1-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-http-error",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobId1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-2-2-name",
      "config": "test-job-config-to-run-manual-job-2",
      "job_id": "test-job",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobId2={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-2-3-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-running",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobId3={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-2-4-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-running",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobId4={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-2-5-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-failed",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobId5={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-2-step-1-operation-1-description",
              "jobs": [
                "{{ .jobId1 }}",
                "{{ .jobId2 }}",
                "{{ .jobId3 }}",
                "{{ .jobId4 }}",
                "{{ .jobId5 }}"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-2-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-2",
      "connector": "test-connector-to-job-execution-start-2",
      "connector_name": "test-connector-name-to-job-execution-start-2",
      "component": "test-component-to-job-execution-start-2",
      "resource": "test-resource-to-job-execution-start-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-2
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
      "resource": "test-resource-to-job-execution-start-2",
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
      "job": "{{ .jobId1 }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobId2 }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobId3 }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobId4 }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobId5 }}"
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
                  "fail_reason": "http-error",
                  "status": 2
                },
                {
                  "fail_reason": "url POST http://not-exist/api/35/job/test-job/run cannot be connected",
                  "status": 2
                },
                {
                  "fail_reason": "job is executing too long, cannot retrieve status after retries",
                  "status": 2
                },
                {
                  "fail_reason": "job is executing too long, cannot retrieve status after retries",
                  "status": 2
                },
                {
                  "fail_reason": "see {{ .dummyApiURL }}/rundeck/execution/show/test-job-execution-failed",
                  "status": 2
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response jobExecId1={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0)._id }}
    When I save response jobExecId2={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 1)._id }}
    When I save response jobExecId3={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 2)._id }}
    When I save response jobExecId4={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 3)._id }}
    When I save response jobExecId5={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 4)._id }}
    When I do GET /api/v4/cat/job-executions/{{ .jobExecId1 }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    """
    When I do GET /api/v4/cat/job-executions/{{ .jobExecId2 }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    """
    When I do GET /api/v4/cat/job-executions/{{ .jobExecId3 }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    """
    When I do GET /api/v4/cat/job-executions/{{ .jobExecId4 }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    """
    When I do GET /api/v4/cat/job-executions/{{ .jobExecId5 }}/output
    Then the response code should be 200
    Then the response raw body should be:
    """
    test-job-execution-failed-output
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "opened": true,
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
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-1-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-2-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-4-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-4-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-5-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-5-name."
      }
    ]
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given job should start job for operation of different instructions multiple times
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-3-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-3-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-3-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-3-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-3-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-3-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-3-1-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-3-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response firstInstructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-3-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-3-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-3-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-3-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-3-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-3-2-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-3-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response secondInstructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-to-job-execution-start-3",
      "connector": "test-connector-to-job-execution-start-3",
      "connector_name": "test-connector-name-to-job-execution-start-3",
      "component": "test-component-to-job-execution-start-3",
      "resource": "test-resource-to-job-execution-start-3-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-3-1
    Then the response code should be 200
    When I save response firstAlarmID={{ (index .lastResponse.data 0)._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-to-job-execution-start-3",
      "connector": "test-connector-to-job-execution-start-3",
      "connector_name": "test-connector-name-to-job-execution-start-3",
      "component": "test-component-to-job-execution-start-3",
      "resource": "test-resource-to-job-execution-start-3-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-3-2
    Then the response code should be 200
    When I save response secondAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .firstAlarmID }}",
      "instruction": "{{ .firstInstructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-3-1",
      "source_type": "resource"
    }
    """
    When I save response firstExecutionID={{ .lastResponse._id }}
    When I save response firstOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .secondAlarmID }}",
      "instruction": "{{ .secondInstructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-3-2",
      "source_type": "resource"
    }
    """
    When I save response secondExecutionID={{ .lastResponse._id }}
    When I save response secondOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .firstExecutionID }}",
      "operation": "{{ .firstOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .secondExecutionID }}",
      "operation": "{{ .secondOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .firstExecutionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .secondExecutionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given jenkins job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-4-name",
      "config": "test-job-config-to-run-manual-jenkins-job",
      "job_id": "test-job-succeeded",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-4"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-4-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-4-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-4",
      "connector": "test-connector-to-job-execution-start-4",
      "connector_name": "test-connector-name-to-job-execution-start-4",
      "component": "test-component-to-job-execution-start-4",
      "resource": "test-resource-to-job-execution-start-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-4
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
      "resource": "test-resource-to-job-execution-start-4",
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
                  "name": "test-job-to-job-execution-start-4-name",
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
  Scenario: given jenkins job with parameters should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-5-name",
      "config": "test-job-config-to-run-manual-jenkins-job",
      "job_id": "test-job-params-succeeded",
      "query": {
        "resource1": "{{ `{{ .Alarm.Value.Resource }}` }}",
        "entity1": "{{ `{{ .Entity.ID }}` }}"
      },
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-5"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-5-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-5-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-5",
      "connector": "test-connector-to-job-execution-start-5",
      "connector_name": "test-connector-name-to-job-execution-start-5",
      "component": "test-component-to-job-execution-start-5",
      "resource": "test-resource-to-job-execution-start-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-5
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
      "resource": "test-resource-to-job-execution-start-5",
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
                  "name": "test-job-to-job-execution-start-5-name",
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
    test-job-execution-params-succeeded-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given awx job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-6-name",
      "config": "test-job-config-to-run-manual-awx-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\", \"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-6"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-6-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-6-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-6",
      "connector": "test-connector-to-job-execution-start-6",
      "connector_name": "test-connector-name-to-job-execution-start-6",
      "component": "test-component-to-job-execution-start-6",
      "resource": "test-resource-to-job-execution-start-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-6
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
      "resource": "test-resource-to-job-execution-start-6",
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
                  "name": "test-job-to-job-execution-start-6-name",
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
  Scenario: given job should queue exclusive job for different executions
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-7-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-7-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-job-execution-start-7-1",
                "test-resource-to-job-execution-start-7-2",
                "test-resource-to-job-execution-start-7-3"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-7-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-7-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-job-execution-start-7",
        "connector": "test-connector-to-job-execution-start-7",
        "connector_name": "test-connector-name-to-job-execution-start-7",
        "component": "test-component-to-job-execution-start-7",
        "resource": "test-resource-to-job-execution-start-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-job-execution-start-7",
        "connector": "test-connector-to-job-execution-start-7",
        "connector_name": "test-connector-name-to-job-execution-start-7",
        "component": "test-component-to-job-execution-start-7",
        "resource": "test-resource-to-job-execution-start-7-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-job-execution-start-7",
        "connector": "test-connector-to-job-execution-start-7",
        "connector_name": "test-connector-name-to-job-execution-start-7",
        "component": "test-component-to-job-execution-start-7",
        "resource": "test-resource-to-job-execution-start-7-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-7-1
    Then the response code should be 200
    When I save response firstAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-7-2
    Then the response code should be 200
    When I save response secondAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-7-3
    Then the response code should be 200
    When I save response thirdAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .firstAlarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-7-1",
      "source_type": "resource"
    }
    """
    When I save response firstExecutionID={{ .lastResponse._id }}
    When I save response firstOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .secondAlarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-7-2",
      "source_type": "resource"
    }
    """
    When I save response secondExecutionID={{ .lastResponse._id }}
    When I save response secondOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .thirdAlarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-7-3",
      "source_type": "resource"
    }
    """
    When I save response thirdExecutionID={{ .lastResponse._id }}
    When I save response thirdOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .firstExecutionID }}",
      "operation": "{{ .firstOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .firstExecutionID }} until response code is 200 and response key "steps.0.operations.0.jobs.0.launched_at" is greater or equal than 1
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .secondExecutionID }}",
      "operation": "{{ .secondOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "queue_number": 1
    }
    """
    When I do GET /api/v4/cat/executions/{{ .secondExecutionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 0,
                  "queue_number": 1
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
      "execution": "{{ .thirdExecutionID }}",
      "operation": "{{ .thirdOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "queue_number": 2
    }
    """
    When I do GET /api/v4/cat/executions/{{ .thirdExecutionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 0,
                  "queue_number": 2
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/executions/{{ .firstExecutionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 1
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response firstJobCompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).completed_at }}
    When I do GET /api/v4/cat/executions/{{ .secondExecutionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 0,
                  "queue_number": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/executions/{{ .thirdExecutionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 0,
                  "queue_number": 1
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/executions/{{ .secondExecutionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 1
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response secondJobLaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).launched_at }}
    When I save response secondJobCompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).completed_at }}
    Then "secondJobLaunchedAt" >= "firstJobCompletedAt"
    When I do GET /api/v4/cat/executions/{{ .thirdExecutionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 1
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response thirdJobLaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).launched_at }}
    Then "thirdJobLaunchedAt" >= "secondJobCompletedAt"
    When I do PUT /api/v4/cat/executions/{{ .firstExecutionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .secondExecutionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .thirdExecutionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given job should run job in parallel for different executions
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-8-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": true
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-8-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-job-execution-start-8-1",
                "test-resource-to-job-execution-start-8-2"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-8-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-8-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-job-execution-start-8",
        "connector": "test-connector-to-job-execution-start-8",
        "connector_name": "test-connector-name-to-job-execution-start-8",
        "component": "test-component-to-job-execution-start-8",
        "resource": "test-resource-to-job-execution-start-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-job-execution-start-8",
        "connector": "test-connector-to-job-execution-start-8",
        "connector_name": "test-connector-name-to-job-execution-start-8",
        "component": "test-component-to-job-execution-start-8",
        "resource": "test-resource-to-job-execution-start-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-8-1
    Then the response code should be 200
    When I save response firstAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-8-2
    Then the response code should be 200
    When I save response secondAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .firstAlarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-8-1",
      "source_type": "resource"
    }
    """
    When I save response firstExecutionID={{ .lastResponse._id }}
    When I save response firstOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .secondAlarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "resource": "test-resource-to-job-execution-start-8-2",
      "source_type": "resource"
    }
    """
    When I save response secondExecutionID={{ .lastResponse._id }}
    When I save response secondOperationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .firstExecutionID }}",
      "operation": "{{ .firstOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .firstExecutionID }} until response code is 200 and response key "steps.0.operations.0.jobs.0.launched_at" is greater or equal than 1
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .secondExecutionID }}",
      "operation": "{{ .secondOperationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "queue_number": null
    }
    """
    When I do GET /api/v4/cat/executions/{{ .secondExecutionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 0,
                  "queue_number": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/executions/{{ .firstExecutionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 1
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response firstJobCompletedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).completed_at }}
    When I do GET /api/v4/cat/executions/{{ .secondExecutionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "status": 1
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I save response secondJobLaunchedAt={{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).launched_at }}
    Then "secondJobLaunchedAt" < "firstJobCompletedAt"
    When I do PUT /api/v4/cat/executions/{{ .firstExecutionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .secondExecutionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given job should not add step to resolved alarm and new alarm
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-9-name",
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
      "name": "test-instruction-to-job-execution-start-9-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-9"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-9-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-9-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-9-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-9-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-9-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-9",
      "connector": "test-connector-to-job-execution-start-9",
      "connector_name": "test-connector-name-to-job-execution-start-9",
      "component": "test-component-to-job-execution-start-9",
      "resource": "test-resource-to-job-execution-start-9",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-9
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
      "resource": "test-resource-to-job-execution-start-9",
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
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and response key "steps.0.operations.0.jobs.0.started_at" is greater or equal than 1
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-to-job-execution-start-9",
      "connector_name": "test-connector-name-to-job-execution-start-9",
      "component": "test-component-to-job-execution-start-9",
      "resource": "test-resource-to-job-execution-start-9",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "connector": "test-connector-to-job-execution-start-9",
      "connector_name": "test-connector-name-to-job-execution-start-9",
      "component": "test-component-to-job-execution-start-9",
      "resource": "test-resource-to-job-execution-start-9",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-to-job-execution-start-9",
      "connector": "test-connector-to-job-execution-start-9",
      "connector_name": "test-connector-name-to-job-execution-start-9",
      "component": "test-component-to-job-execution-start-9",
      "resource": "test-resource-to-job-execution-start-9",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-9-name",
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
    test-job-execution-long-succeeded-output
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I save request:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "opened": false,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains only:
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
        "m": "Instruction test-instruction-to-job-execution-start-9-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-job-execution-start-9-name. Job test-job-to-job-execution-start-9-name."
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
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-9&opened=true
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
