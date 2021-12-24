Feature: run a job
  I need to be able to run a job
  Only admin should be able to run a job

  Scenario: given job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-1-1-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response job1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-1-2-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource2\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity2\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response job2ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-1-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-1"
        }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-1",
      "connector_name": "test-connector-name-to-job-execution-start-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-1",
      "resource": "test-resource-to-job-execution-start-1",
      "state": 1,
      "output": "test-output-to-job-execution-start-1"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of event processing
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
      "status": 0,
      "fail_reason": "",
      "output": "",
      "payload": "{\"resource2\": \"test-resource-to-job-execution-start-1\",\"entity2\": \"test-resource-to-job-execution-start-1/test-component-to-job-execution-start-1\"}",
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-1-step-1-operation-1",
              "completed_at": 0,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "_id": "",
                  "name": "test-job-to-job-execution-start-1-1-name",
                  "status": null,
                  "fail_reason": "",
                  "payload": "",
                  "started_at": 0,
                  "launched_at": 0,
                  "completed_at": 0
                },
                {
                  "name": "test-job-to-job-execution-start-1-2-name",
                  "status": 1,
                  "fail_reason": "",
                  "output": "test-job-execution-succeeded-output",
                  "payload": "{\"resource2\": \"test-resource-to-job-execution-start-1\",\"entity2\": \"test-resource-to-job-execution-start-1/test-component-to-job-execution-start-1\"}"
                }
              ]
            }
          ]
        }
      ]
    }
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
      "status": 0,
      "fail_reason": "",
      "output": "",
      "payload": "{\"resource1\": \"test-resource-to-job-execution-start-1\",\"entity1\": \"test-resource-to-job-execution-start-1/test-component-to-job-execution-start-1\"}",
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-1-step-1-operation-1",
              "completed_at": 0,
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-1-1-name",
                  "status": 1,
                  "fail_reason": "",
                  "output": "test-job-execution-succeeded-output",
                  "payload": "{\"resource1\": \"test-resource-to-job-execution-start-1\",\"entity1\": \"test-resource-to-job-execution-start-1/test-component-to-job-execution-start-1\"}"
                },
                {
                  "name": "test-job-to-job-execution-start-1-2-name",
                  "status": 1,
                  "fail_reason": "",
                  "output": "test-job-execution-succeeded-output",
                  "payload": "{\"resource2\": \"test-resource-to-job-execution-start-1\",\"entity2\": \"test-resource-to-job-execution-start-1/test-component-to-job-execution-start-1\"}"
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-1&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-1-name."
              },
              {
                "_t": "instructionjobstart",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-2-name."
              },
              {
                "_t": "instructionjobcomplete",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-2-name."
              },
              {
                "_t": "instructionjobstart",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-1-name."
              },
              {
                "_t": "instructionjobcomplete",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-1-name. Job test-job-to-job-execution-start-1-1-name."
              },
              {
                "_t": "instructioncomplete",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-1-name."
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given http error during job execution should return failed job status
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-2-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-http-error",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-2"
        }
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
              "jobs": ["{{ .jobID }}"]
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-2",
      "connector_name": "test-connector-name-to-job-execution-start-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-2",
      "resource": "test-resource-to-job-execution-start-2",
      "state": 1,
      "output": "test-output-to-job-execution-start-2"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of event processing
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
    When I wait 2s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-2&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-2-name."
              },
              {
                "_t": "instructionjobstart",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-name."
              },
              {
                "_t": "instructionjobfail",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-2-name. Job test-job-to-job-execution-start-2-name."
              },
              {
                "_t": "instructioncomplete",
                "a": "root",
                "m": "Instruction test-instruction-to-job-execution-start-2-name."
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given job should not start job for operation of instruction multiple times
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-3-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-long-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-3-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-3"
        }
      ],
      "description": "test-instruction-to-job-execution-start-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-3-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-3",
      "connector_name": "test-connector-name-to-job-execution-start-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-3",
      "resource": "test-resource-to-job-execution-start-3",
      "state": 1,
      "output": "test-output-to-job-execution-start-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-3
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
    When I wait the end of event processing
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
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "job is already running for operation"
    }
    """
    When I wait 2s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing

  Scenario: given job should start job for operation of different instructions multiple times
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-4-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-4-1-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-4-1"
        }
      ],
      "description": "test-instruction-to-job-execution-start-4-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-4-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-4-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-4-1-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-4-1-step-1-endpoint"
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
      "name": "test-instruction-to-job-execution-start-4-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-4-2"
        }
      ],
      "description": "test-instruction-to-job-execution-start-4-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-4-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-4-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-4-2-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-4-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response secondInstructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-4",
      "connector_name": "test-connector-name-to-job-execution-start-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-4",
      "resource": "test-resource-to-job-execution-start-4-1",
      "state": 1,
      "output": "test-output-to-job-execution-start-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-4-1
    Then the response code should be 200
    When I save response firstAlarmID={{ (index .lastResponse.data 0)._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-4",
      "connector_name": "test-connector-name-to-job-execution-start-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-4",
      "resource": "test-resource-to-job-execution-start-4-2",
      "state": 1,
      "output": "test-output-to-job-execution-start-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-4-2
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
    When I wait the end of event processing
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
    When I wait the end of event processing
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
    When I wait 2s
    When I do PUT /api/v4/cat/executions/{{ .firstExecutionID }}/next-step
    When I do PUT /api/v4/cat/executions/{{ .secondExecutionID }}/next-step
    When I wait the end of 2 events processing

  Scenario: given job should not start job for not running operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-5-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-5-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-5"
        }
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
              "jobs": []
            },
            {
              "name": "test-instruction-to-job-execution-start-5-step-1-operation-2",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-5-step-1-operation-2-description",
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-5",
      "connector_name": "test-connector-name-to-job-execution-start-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-5",
      "resource": "test-resource-to-job-execution-start-5",
      "state": 1,
      "output": "test-output-to-job-execution-start-5"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of event processing
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 1).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 404
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing

  Scenario: given job should not start job for not running instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-6-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-6-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-6"
        }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-6",
      "connector_name": "test-connector-name-to-job-execution-start-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-6",
      "resource": "test-resource-to-job-execution-start-6",
      "state": 1,
      "output": "test-output-to-job-execution-start-6"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of event processing
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 404

  Scenario: given job with invalid payload should return error
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-7-name",
      "config": "test-job-config-to-run-manual-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.ResourceBadValue }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-7-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-7"
        }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-7",
      "connector_name": "test-connector-name-to-job-execution-start-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-7",
      "resource": "test-resource-to-job-execution-start-7",
      "state": 1,
      "output": "test-output-to-job-execution-start-7"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-7
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
    When I wait the end of event processing
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
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "payload is not valid"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing

  Scenario: given jenkins job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-8-name",
      "config": "test-job-config-to-run-manual-jenkins-job",
      "job_id": "test-job-succeeded"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-8-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-8"
        }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-8",
      "connector_name": "test-connector-name-to-job-execution-start-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-8",
      "resource": "test-resource-to-job-execution-start-8",
      "state": 1,
      "output": "test-output-to-job-execution-start-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-8
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
    When I wait the end of event processing
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
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-job-execution-start-8-name",
      "status": 0,
      "fail_reason": "",
      "output": "",
      "payload": "",
      "query": null,
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-8-name",
                  "status": 1,
                  "fail_reason": "",
                  "output": "test-job-execution-succeeded-output",
                  "payload": "",
                  "query": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing

  Scenario: given jenkins job with parameters should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-9-name",
      "config": "test-job-config-to-run-manual-jenkins-job",
      "job_id": "test-job-params-succeeded",
      "query": {
        "resource1": "{{ `{{ .Alarm.Value.Resource }}` }}",
        "entity1": "{{ `{{ .Entity.ID }}` }}"
      }
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-9-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-9"
        }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-9",
      "connector_name": "test-connector-name-to-job-execution-start-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-9",
      "resource": "test-resource-to-job-execution-start-9",
      "state": 1,
      "output": "test-output-to-job-execution-start-9"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of event processing
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
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-job-execution-start-9-name",
      "status": 0,
      "fail_reason": "",
      "output": "",
      "payload": "",
      "query": {
        "resource1": "test-resource-to-job-execution-start-9",
        "entity1": "test-resource-to-job-execution-start-9/test-component-to-job-execution-start-9"
      },
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
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
                  "fail_reason": "",
                  "output": "test-job-execution-params-succeeded-output",
                  "payload": "",
                  "query": {
                    "resource1": "test-resource-to-job-execution-start-9",
                    "entity1": "test-resource-to-job-execution-start-9/test-component-to-job-execution-start-9"
                  }
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing

  Scenario: given awx job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-10-name",
      "config": "test-job-config-to-run-manual-awx-job",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\", \"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-10-name",
      "entity_patterns": [
        {
          "name": "test-resource-to-job-execution-start-10"
        }
      ],
      "description": "test-instruction-to-job-execution-start-10-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-10-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-10-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-10-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-10-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-to-job-execution-start-10",
      "connector_name": "test-connector-name-to-job-execution-start-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-job-execution-start-10",
      "resource": "test-resource-to-job-execution-start-10",
      "state": 1,
      "output": "test-output-to-job-execution-start-10"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-10
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
    When I wait the end of event processing
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
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-job-execution-start-10-name",
      "status": 0,
      "fail_reason": "",
      "output": "",
      "payload": "{\"resource1\": \"test-resource-to-job-execution-start-10\", \"entity1\": \"test-resource-to-job-execution-start-10/test-component-to-job-execution-start-10\"}",
      "query": null,
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
    When I do GET /api/v4/cat/executions/{{ .executionID }} until response code is 200 and body contains:
    """json
    {
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                {
                  "name": "test-job-to-job-execution-start-10-name",
                  "status": 1,
                  "fail_reason": "",
                  "output": "test-job-execution-succeeded-output",
                  "payload": "{\"resource1\": \"test-resource-to-job-execution-start-10\", \"entity1\": \"test-resource-to-job-execution-start-10/test-component-to-job-execution-start-10\"}",
                  "query": null
                }
              ]
            }
          ]
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    When I wait the end of event processing

  Scenario: given unauth request should not allow access
    When I do POST /api/v4/cat/job-executions
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/job-executions
    Then the response code should be 403
