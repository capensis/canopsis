Feature: run a job
  I need to be able to run a job
  Only admin should be able to run a job

  @concurrent
  Scenario: given unauth start job request should not allow access
    When I do POST /api/v4/cat/job-executions
    Then the response code should be 401

  @concurrent
  Scenario: given start job request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/job-executions
    Then the response code should be 403

  @concurrent
  Scenario: given unauth get output request should not allow access
    When I do GET /api/v4/cat/job-executions/test-job-not-exist/output
    Then the response code should be 401

  @concurrent
  Scenario: given get output request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/job-executions/test-job-not-exist/output
    Then the response code should be 403

  @concurrent
  Scenario: given not exist id in get output request should return error
    When I am admin
    When I do GET /api/v4/cat/job-executions/test-job-not-exist/output
    Then the response code should be 404

  @concurrent
  Scenario: given job should not start job for operation of instruction multiple times
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-fail-1-name",
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
      "name": "test-instruction-to-job-execution-start-fail-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-fail-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-fail-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-fail-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-fail-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-fail-1-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-fail-1-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-fail-1",
      "connector": "test-connector-to-job-execution-start-fail-1",
      "connector_name": "test-connector-name-to-job-execution-start-fail-1",
      "component": "test-component-to-job-execution-start-fail-1",
      "resource": "test-resource-to-job-execution-start-fail-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-fail-1
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
      "resource": "test-resource-to-job-execution-start-fail-1",
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
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "trigger",
      "resource": "test-resource-to-job-execution-start-fail-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given job should not start job for not running operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-fail-2-name",
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
      "name": "test-instruction-to-job-execution-start-fail-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-fail-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-fail-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-fail-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-fail-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-fail-2-step-1-operation-1-description",
              "jobs": []
            },
            {
              "name": "test-instruction-to-job-execution-start-fail-2-step-1-operation-2",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-fail-2-step-1-operation-2-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-fail-2-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-fail-2",
      "connector": "test-connector-to-job-execution-start-fail-2",
      "connector_name": "test-connector-name-to-job-execution-start-fail-2",
      "component": "test-component-to-job-execution-start-fail-2",
      "resource": "test-resource-to-job-execution-start-fail-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-fail-2
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
      "resource": "test-resource-to-job-execution-start-fail-2",
      "source_type": "resource"
    }
    """
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
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200

  @concurrent
  Scenario: given job should not start job for not running instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-fail-3-name",
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
      "name": "test-instruction-to-job-execution-start-fail-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-fail-3"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-fail-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-fail-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-fail-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-fail-3-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-fail-3-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-fail-3",
      "connector": "test-connector-to-job-execution-start-fail-3",
      "connector_name": "test-connector-name-to-job-execution-start-fail-3",
      "component": "test-component-to-job-execution-start-fail-3",
      "resource": "test-resource-to-job-execution-start-fail-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-fail-3
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
      "resource": "test-resource-to-job-execution-start-fail-3",
      "source_type": "resource"
    }
    """
    When I save response executionID={{ .lastResponse._id }}
    When I save response operationID={{ (index (index .lastResponse.steps 0).operations 0).operation_id }}
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "resource": "test-resource-to-job-execution-start-fail-3",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/cat/job-executions:
    """json
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .jobID }}"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given job with invalid payload should return error
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-job-execution-start-fail-4-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.ResourceBadValue }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response jobID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-job-execution-start-fail-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-job-execution-start-fail-4"
            }
          }
        ]
      ],
      "description": "test-instruction-to-job-execution-start-fail-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-job-execution-start-fail-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-job-execution-start-fail-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-job-execution-start-fail-4-step-1-operation-1-description",
              "jobs": ["{{ .jobID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-job-execution-start-fail-4-step-1-endpoint"
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
      "output": "test-output-to-job-execution-start-fail-4",
      "connector": "test-connector-to-job-execution-start-fail-4",
      "connector_name": "test-connector-name-to-job-execution-start-fail-4",
      "component": "test-component-to-job-execution-start-fail-4",
      "resource": "test-resource-to-job-execution-start-fail-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-job-execution-start-fail-4
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
      "resource": "test-resource-to-job-execution-start-fail-4",
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
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "payload is not valid"
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
