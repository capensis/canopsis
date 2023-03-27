Feature: get an execution status 
  I need to be able to get an execution status

  @concurrent
  Scenario: given manual instruction should get jobs statuses from websocket room
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-manual-instruction-websocket-1-1-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response job1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-manual-instruction-websocket-1-2-name",
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
      "name": "test-instruction-manual-instruction-websocket-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-manual-instruction-websocket-1"
            }
          }
        ]
      ],
      "description": "test-instruction-manual-instruction-websocket-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-manual-instruction-websocket-1-step-1",
          "operations": [
            {
              "name": "test-instruction-manual-instruction-websocket-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-manual-instruction-websocket-1-step-1-operation-1-description",
              "jobs": ["{{ .job1ID }}", "{{ .job2ID }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-manual-instruction-websocket-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-manual-instruction-websocket-1",
      "connector_name": "test-connector-name-manual-instruction-websocket-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-manual-instruction-websocket-1",
      "resource": "test-resource-manual-instruction-websocket-1",
      "state": 1,
      "output": "test-output-manual-instruction-websocket-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-manual-instruction-websocket-1
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
    When I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "execution/{{ .executionID }}"
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
    Then I wait message from websocket room "execution/{{ .executionID }}" which contains:
    """json
    {
      "name": "test-instruction-manual-instruction-websocket-1-step-1-operation-1",
      "operation_id": "{{ .operationID }}",
      "description": "test-instruction-manual-instruction-websocket-1-step-1-operation-1-description",
      "completed_at": null,
      "time_to_complete": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "_id": "",
          "job_id": "{{ .job1ID }}",
          "name": "test-job-manual-instruction-websocket-1-1-name",
          "status": null,
          "fail_reason": "",
          "started_at": null,
          "launched_at": null,
          "completed_at": null,
          "queue_number": null
        },
        {
          "job_id": "{{ .job2ID }}",
          "name": "test-job-manual-instruction-websocket-1-2-name",
          "status": 1,
          "fail_reason": "",
          "queue_number": null
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
    Then I wait message from websocket room "execution/{{ .executionID }}" which contains:
    """json
    {
      "name": "test-instruction-manual-instruction-websocket-1-step-1-operation-1",
      "operation_id": "{{ .operationID }}",
      "description": "test-instruction-manual-instruction-websocket-1-step-1-operation-1-description",
      "completed_at": null,
      "time_to_complete": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job_id": "{{ .job1ID }}",
          "name": "test-job-manual-instruction-websocket-1-1-name",
          "status": 1,
          "fail_reason": "",
          "queue_number": null
        },
        {
          "job_id": "{{ .job2ID }}",
          "name": "test-job-manual-instruction-websocket-1-2-name",
          "status": 1,
          "fail_reason": "",
          "queue_number": null
        }
      ]
    }
    """

  @concurrent
  Scenario: given simplified manual instruction should get jobs statuses from websocket room
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-manual-instruction-websocket-2-1-name",
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
      "name": "test-job-manual-instruction-websocket-2-2-name",
      "config": "test-job-config-to-run-manual-job-1",
      "job_id": "test-job-succeeded",
      "payload": "{\"resource2\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity2\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I save response job2ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-manual-instruction-websocket-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-manual-instruction-websocket-2"
            }
          }
        ]
      ],
      "description": "test-instruction-manual-instruction-websocket-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "{{ .job1ID }}",
          "stop_on_fail": false
        },
        {
          "job": "{{ .job2ID }}"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-manual-instruction-websocket-2",
      "connector_name": "test-connector-name-manual-instruction-websocket-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-manual-instruction-websocket-2",
      "resource": "test-resource-manual-instruction-websocket-2",
      "state": 1,
      "output": "test-output-manual-instruction-websocket-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-manual-instruction-websocket-2
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
    When I save response executionID={{ .lastResponse._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "execution/{{ .executionID }}"
    Then I wait message from websocket room "execution/{{ .executionID }}" which contains:
    """json
    [
      {
        "job_id": "{{ .job1ID }}",
        "name": "test-job-manual-instruction-websocket-2-1-name",
        "status": 0,
        "fail_reason": "",
        "completed_at": null,
        "queue_number": null
      },
      {
        "job_id": "{{ .job2ID }}",
        "name": "test-job-manual-instruction-websocket-2-2-name",
        "status": 0,
        "fail_reason": "",
        "started_at": null,
        "launched_at": null,
        "completed_at": null,
        "queue_number": null
      }
    ]
    """
    Then I wait message from websocket room "execution/{{ .executionID }}" which contains:
    """json
    [
      {
        "job_id": "{{ .job1ID }}",
        "name": "test-job-manual-instruction-websocket-2-1-name",
        "status": 1,
        "fail_reason": "",
        "queue_number": null
      },
      {
        "job_id": "{{ .job2ID }}",
        "name": "test-job-manual-instruction-websocket-2-2-name",
        "status": 1,
        "fail_reason": "",
        "queue_number": null
      }
    ]
    """
