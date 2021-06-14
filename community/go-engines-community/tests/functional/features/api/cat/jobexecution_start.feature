Feature: run a job
  I need to be able to run a job
  Only admin should be able to run a job

  Scenario: given job should start job for operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-execution-start-1-1-name",
      "config": "test-job-config-to-link",
      "job_id": "test-job-execution-start-1-job-id",
      "payload": "{\"resource1\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity1\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    Then I save response job1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-execution-start-1-2-name",
      "config": "test-job-config-to-link",
      "job_id": "test-job-execution-start-1-job-id",
      "payload": "{\"resource2\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity2\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    Then I save response job2ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-job-execution-start-1-name",
      "alarm_patterns": [
        {
          "_id": "test-job-execution-start-1"
        }
      ],
      "description": "test-job-execution-start-1-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-job-execution-start-1-step-1",
          "operations": [
            {
              "name": "test-job-execution-start-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-1-step-1-operation-1-description",
              "jobs": []
            },
            {
              "name": "test-job-execution-start-1-step-1-operation-2",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-1-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-job-execution-start-1-step-1-endpoint"
        },
        {
          "name": "test-job-execution-start-1-step-2",
          "operations": [
            {
              "name": "test-job-execution-start-1-step-2-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-1-step-2-operation-1-description",
              "jobs": []
            },
            {
              "name": "test-job-execution-start-1-step-2-operation-2",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-1-step-2-operation-2-description",
              "jobs": ["{{ .job1ID }}", "{{ .job2ID }}"]
            },
            {
              "name": "test-job-execution-start-1-step-2-operation-3",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-1-step-2-operation-3-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-job-execution-start-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-job-execution-start-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    Then I save response operationID={{ (index (index .lastResponse.steps 1).operations 1).operation_id }}
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .job2ID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-job-execution-start-1-2-name",
      "status": 0,
      "fail_reason": "",
      "payload": "{\"resource2\": \"test-job-execution-start-resource-1\",\"entity2\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "steps": [
        {},
        {
          "operations": [
            {},
            {
              "name": "test-job-execution-start-1-step-2-operation-2",
              "description": "test-job-execution-start-1-step-2-operation-2-description",
              "completed_at": 0,
              "time_to_complete": {
                "seconds": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "_id": "",
                  "name": "test-job-execution-start-1-1-name",
                  "status": null,
                  "fail_reason": "",
                  "payload": "",
                  "started_at": 0,
                  "launched_at": 0,
                  "completed_at": 0
                },
                {
                  "name": "test-job-execution-start-1-2-name",
                  "status": 0,
                  "fail_reason": "",
                  "payload": "{\"resource2\": \"test-job-execution-start-resource-1\",\"entity2\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
                  "launched_at": 0,
                  "completed_at": 0
                }
              ]
            },
            {}
          ]
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/ping
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-job-execution-start-1-step-2-operation-2",
      "description": "test-job-execution-start-1-step-2-operation-2-description",
      "completed_at": 0,
      "time_to_complete": {
        "seconds": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "_id": "",
          "name": "test-job-execution-start-1-1-name",
          "status": null,
          "fail_reason": "",
          "payload": "",
          "started_at": 0,
          "launched_at": 0,
          "completed_at": 0
        },
        {
          "name": "test-job-execution-start-1-2-name",
          "status": 0,
          "fail_reason": "",
          "payload": "{\"resource2\": \"test-job-execution-start-resource-1\",\"entity2\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
          "launched_at": 0,
          "completed_at": 0
        }
      ]
    }
    """
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "{{ .executionID }}",
      "operation": "{{ .operationID }}",
      "job": "{{ .job1ID }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-job-execution-start-1-1-name",
      "status": 0,
      "fail_reason": "",
      "payload": "{\"resource1\": \"test-job-execution-start-resource-1\",\"entity1\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
      "launched_at": 0,
      "completed_at": 0
    }
    """
    Then the response key "started_at" should not be "0"
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "steps": [
        {},
        {
          "operations": [
            {},
            {
              "name": "test-job-execution-start-1-step-2-operation-2",
              "description": "test-job-execution-start-1-step-2-operation-2-description",
              "completed_at": 0,
              "time_to_complete": {
                "seconds": 1,
                "unit": "s"
              },
              "jobs": [
                {
                  "name": "test-job-execution-start-1-1-name",
                  "status": 0,
                  "fail_reason": "",
                  "payload": "{\"resource1\": \"test-job-execution-start-resource-1\",\"entity1\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
                  "launched_at": 0,
                  "completed_at": 0
                },
                {
                  "name": "test-job-execution-start-1-2-name",
                  "status": 0,
                  "fail_reason": "",
                  "payload": "{\"resource2\": \"test-job-execution-start-resource-1\",\"entity2\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
                  "launched_at": 0,
                  "completed_at": 0
                }
              ]
            },
            {}
          ]
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/ping
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-job-execution-start-1-step-2-operation-2",
      "description": "test-job-execution-start-1-step-2-operation-2-description",
      "completed_at": 0,
      "time_to_complete": {
        "seconds": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "name": "test-job-execution-start-1-1-name",
          "status": 0,
          "fail_reason": "",
          "payload": "{\"resource1\": \"test-job-execution-start-resource-1\",\"entity1\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
          "launched_at": 0,
          "completed_at": 0
        },
        {
          "name": "test-job-execution-start-1-2-name",
          "status": 0,
          "fail_reason": "",
          "payload": "{\"resource2\": \"test-job-execution-start-resource-1\",\"entity2\": \"test-job-execution-start-resource-1/test-job-execution-start-component\"}",
          "launched_at": 0,
          "completed_at": 0
        }
      ]
    }
    """

  Scenario: given job should not start job for operation of instruction multiple times
    When I am admin
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "test-job-execution-to-start-multiple-times-1",
      "operation": "test-job-execution-to-start-multiple-times-1-step-1-operation-1",
      "job": "test-job-execution-multiple-times-1"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "test-job-execution-to-start-multiple-times-1",
      "operation": "test-job-execution-to-start-multiple-times-1-step-1-operation-1",
      "job": "test-job-execution-multiple-times-1"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "job is already running for operation"
    }
    """

  Scenario: given job should start job for different instructions multiple times
    When I am admin
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "test-job-execution-to-start-multiple-times-2",
      "operation": "test-job-execution-to-start-multiple-times-2-step-1-operation-1",
      "job": "test-job-execution-multiple-times-2"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "test-job-execution-to-start-multiple-times-3",
      "operation": "test-job-execution-to-start-multiple-times-3-step-1-operation-1",
      "job": "test-job-execution-multiple-times-2"
    }
    """
    Then the response code should be 200

  Scenario: given job should not start job for not running operation of instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-execution-start-2-name",
      "config": "test-job-config-to-link",
      "job_id": "test-job-execution-start-2-job-id",
      "payload": "{}"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-job-execution-start-2-name",
      "alarm_patterns": [
        {
          "_id": "test-job-execution-start-2"
        }
      ],
      "description": "test-job-execution-start-2-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-job-execution-start-2-step-1",
          "operations": [
            {
              "name": "test-job-execution-start-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-2-step-1-operation-1-description",
              "jobs": ["{{ .lastResponse._id }}"]
            },
            {
              "name": "test-job-execution-start-2-step-1-operation-2",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-job-execution-start-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-job-execution-start-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "{{ .lastResponse._id }}",
      "operation": "{{ (index (index .lastResponse.steps 0).operations 0).operation_id }}",
      "job": "{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).job_id }}"
    }
    """
    Then the response code should be 404

  Scenario: given job should not start job for not running instruction
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-execution-start-3-name",
      "config": "test-job-config-to-link",
      "job_id": "test-job-execution-start-3-job-id",
      "payload": "{}"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-job-execution-start-3-name",
      "alarm_patterns": [
        {
          "_id": "test-job-execution-start-3"
        }
      ],
      "description": "test-job-execution-start-3-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-job-execution-start-3-step-1",
          "operations": [
            {
              "name": "test-job-execution-start-3-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-3-step-1-operation-1-description",
              "jobs": ["{{ .lastResponse._id }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-job-execution-start-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-job-execution-start-3",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "{{ .lastResponse._id }}",
      "operation": "{{ (index (index .lastResponse.steps 0).operations 0).operation_id }}",
      "job": "{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).job_id }}"
    }
    """
    Then the response code should be 404

  Scenario: given job should not start job for operation of instruction with invalid payload
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-execution-start-4-name",
      "config": "test-job-config-to-link",
      "job_id": "test-job-execution-start-4-job-id",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.ResourceBadValue }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-job-execution-start-4-name",
      "alarm_patterns": [
        {
          "_id": "test-job-execution-start-4"
        }
      ],
      "description": "test-job-execution-start-4-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-job-execution-start-4-step-1",
          "operations": [
            {
              "name": "test-job-execution-start-4-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-job-execution-start-4-step-1-operation-1-description",
              "jobs": ["{{ .lastResponse._id }}"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-job-execution-start-4-step-1-endpoint"
        }
      ]
    }
    """
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-job-execution-start-4",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    When I do POST /api/v4/cat/job-executions:
    """
    {
      "execution": "{{ .lastResponse._id }}",
      "operation": "{{ (index (index .lastResponse.steps 0).operations 0).operation_id }}",
      "job": "{{ (index (index (index .lastResponse.steps 0).operations 0).jobs 0).job_id }}"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "payload is not valid"
    }
    """

  Scenario: given unauth request should not allow access
    When I do POST /api/v4/cat/job-executions
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/job-executions
    Then the response code should be 403